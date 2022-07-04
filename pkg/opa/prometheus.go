package opa

import (
	"context"
	"fmt"

	"github.com/csullivanupgrade/opa-exporter/internal/config"
	"github.com/csullivanupgrade/opa-exporter/internal/log"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	ConstraintInformation *prometheus.Desc
	ConstraintViolation   *prometheus.Desc
	Metrics               []prometheus.Metric
	Namespace             string
	Up                    *prometheus.Desc
}

func NewExporter(cfg config.Config) *Exporter {
	return &Exporter{
		Namespace: cfg.Namespace,
		Metrics:   []prometheus.Metric{},
		ConstraintInformation: prometheus.NewDesc(
			prometheus.BuildFQName(cfg.Namespace, "", "constraint_information"),
			"Some general information of all constraints",
			[]string{"kind", "name", "enforcementAction", "totalViolations"},
			nil,
		),
		ConstraintViolation: prometheus.NewDesc(
			prometheus.BuildFQName(cfg.Namespace, "", "constraint_violations"),
			"OPA violations for all constraints",
			[]string{
				"kind",
				"name",
				"violating_kind",
				"violating_name",
				"violating_namespace",
				"violation_msg",
				"violation_enforcement",
			},
			nil,
		),
		Up: prometheus.NewDesc(
			prometheus.BuildFQName(cfg.Namespace, "", "up"),
			"Was the last OPA exporter query successful.",
			nil,
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.ConstraintInformation
	ch <- e.ConstraintViolation
	ch <- e.Up
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		e.Up, prometheus.GaugeValue, 1,
	)
	for _, m := range e.Metrics {
		ch <- m
	}
}

func ExportViolations(ctx context.Context, cv *prometheus.Desc, constraints []Constraint) []prometheus.Metric {
	logger := log.FromContext(ctx)

	unique := make(map[string]bool)
	m := make([]prometheus.Metric, 0)
	for _, c := range constraints {
		for _, v := range c.Status.Violations {
			key := fmt.Sprintf("%v-%v-%v-%v-%v", c.Meta.Kind, c.Meta.Name, v.Name, v.Namespace, v.Message)
			if _, ok := unique[key]; ok {
				logger.Warn("found duplicate metric", zap.String("metric", key))
				continue
			}
			unique[key] = true
			metric := prometheus.MustNewConstMetric(
				cv,
				prometheus.GaugeValue,
				1,
				c.Meta.Kind,
				c.Meta.Name,
				v.Kind,
				v.Name,
				v.Namespace,
				v.Message,
				v.EnforcementAction,
			)
			m = append(m, metric)
		}
	}
	return m
}

func ExportConstraintInformation(ci *prometheus.Desc, constraints []Constraint) []prometheus.Metric {
	m := make([]prometheus.Metric, 0)
	for _, c := range constraints {
		metric := prometheus.MustNewConstMetric(
			ci,
			prometheus.GaugeValue,
			c.Status.TotalViolations,
			c.Meta.Kind,
			c.Meta.Name,
			c.Spec.EnforcementAction,
			fmt.Sprintf("%f", c.Status.TotalViolations),
		)
		m = append(m, metric)
	}
	return m
}
