// Package opa provides prometheus metrics for OPA and methods
package opa

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/csullivanupgrade/opa-exporter/internal/log"

	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	controllerClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ConstraintMeta represents meta information of a constraint
type ConstraintMeta struct {
	Kind string
	Name string
}

// Violation represents each constraintViolation under status
type Violation struct {
	Kind              string `json:"kind"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace,omitempty"`
	Message           string `json:"message"`
	EnforcementAction string `json:"enforcementAction"`
}

type ConstraintStatus struct {
	TotalViolations float64 `json:"totalViolations,omitempty"`
	Violations      []*Violation
}

type Constraint struct {
	Meta   ConstraintMeta
	Spec   ConstraintSpec
	Status ConstraintStatus
}

// ConstraintSpec collect general information about the overall constraints applied to the cluster
type ConstraintSpec struct {
	EnforcementAction string `json:"enforcementAction"`
}

const (
	constraintsGV = "constraints.gatekeeper.sh/v1beta1"
)

func createConfig(inCluster *bool) (*rest.Config, error) {
	if *inCluster {
		return rest.InClusterConfig()
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeconfig := filepath.Join(home, ".kube", "config")

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func createKubeClient(inCluster *bool) (*kubernetes.Clientset, error) {
	config, err := createConfig(inCluster)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func createKubeClientGroupVersion(inCluster *bool) (controllerClient.Client, error) {
	config, err := createConfig(inCluster)
	if err != nil {
		return nil, err
	}

	client, err := controllerClient.New(config, controllerClient.Options{})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func checkAndAddConstraint(obj map[string]interface{}, kind string, name string) (*Constraint, error) {
	var constraint Constraint
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &constraint)
	if err != nil {
		return nil, err
	}

	return &Constraint{
		Meta: ConstraintMeta{Kind: kind, Name: name},
		Status: ConstraintStatus{
			TotalViolations: constraint.Status.TotalViolations,
			Violations:      constraint.Status.Violations,
		},
		Spec: ConstraintSpec{EnforcementAction: constraint.Spec.EnforcementAction},
	}, nil
}

func getConstraints(group string, kind string, inCluster *bool) ([]unstructured.Unstructured, error) {
	cClient, err := createKubeClientGroupVersion(inCluster)
	if err != nil {
		return nil, err
	}

	actual := &unstructured.UnstructuredList{}
	actual.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: constraintsGV,
	})

	if err = cClient.List(context.Background(), actual); err != nil {
		return nil, err
	}

	return actual.Items, nil
}

func getAPIResources(inCluster *bool) (*v1.APIResourceList, error) {
	client, err := createKubeClient(inCluster)
	if err != nil {
		return nil, err
	}

	apiResources, err := client.ServerResourcesForGroupVersion(constraintsGV)
	if err != nil {
		return nil, err
	}

	return apiResources, nil
}

// GetConstraints returns a list of all OPA constraints
// nolint:gocognit // would be nice to reduce complexity - I don't see a straightforward path ATM.
func GetConstraints(ctx context.Context, inCluster *bool) ([]Constraint, error) {
	logger := log.FromContext(ctx)

	apiResources, err := getAPIResources(inCluster)
	if err != nil {
		return nil, err
	}

	var constraints []Constraint
	for _, r := range apiResources.APIResources {
		if strings.HasSuffix(r.Name, "/status") {
			continue
		}

		items, err := getConstraints(r.Group, r.Kind, inCluster)
		if err != nil {
			logger.Error("error listing", zap.Error(err))
			continue
		}

		if len(items) > 0 {
			for _, item := range items {
				kind := item.GetKind()
				name := item.GetName()
				namespace := item.GetNamespace()
				constraint, err := checkAndAddConstraint(item.Object, kind, name)
				if err != nil {
					logger.Error(
						"error when checking constraint",
						zap.Error(err), zap.String("kind", kind),
						zap.String("name", name),
						zap.String("namespace", namespace),
					)
					continue
				}
				logger.Info(
					"added constraint",
					zap.Error(err), zap.String("kind", kind),
					zap.String("name", name),
					zap.String("namespace", namespace),
				)
				constraints = append(constraints, *constraint)
			}
		} else {
			logger.Info("nothing returned for this kind", zap.String("kind", r.Kind))
		}
	}

	return constraints, nil
}
