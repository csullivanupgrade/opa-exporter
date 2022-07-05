// Package server handles HTTP requests for the exporter
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/csullivanupgrade/opa-exporter/internal/config"
	"github.com/csullivanupgrade/opa-exporter/internal/log"
	"github.com/csullivanupgrade/opa-exporter/pkg/opa"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startScheduled(ctx context.Context, e *opa.Exporter, cfg config.Config) {
	logger := log.FromContext(ctx)

	done := make(chan bool)
	ticker := time.NewTicker(cfg.Interval)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				logger.Info("ticker", zap.Time("tick", t))
				constraints, err := opa.GetConstraints(ctx, &cfg.InCluster)
				if err != nil {
					logger.Error("could not get contraints", zap.Error(err))
				}
				logger.Info("found constraints", zap.Int("no_constraints", len(constraints)))
				allMetrics := make([]prometheus.Metric, 0)
				violationMetrics := opa.ExportViolations(ctx, e.ConstraintViolation, constraints)
				allMetrics = append(allMetrics, violationMetrics...)

				constraintInformationMetrics := opa.ExportConstraintInformation(e.ConstraintInformation, constraints)
				allMetrics = append(allMetrics, constraintInformationMetrics...)

				e.Metrics = allMetrics
			}
		}
	}()
}

func Run(ctx context.Context, cfg config.Config) {
	logger := log.FromContext(ctx)

	exporter := opa.NewExporter(cfg)
	startScheduled(ctx, exporter, cfg)
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.MustRegister(exporter)

	pattern := fmt.Sprintf("/%s", cfg.Path)
	http.Handle(pattern, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
             <head><title>OPA Exporter</title></head>
             <body>
             <h1>OPA Exporter</h1>
             <p><a href='` + cfg.Path + `'>Metrics</a></p>
             </body>
             </html>`))
		if err != nil {
			logger.Error("error handling reseponse", zap.Error(err))
		}
	})
	bind := fmt.Sprintf(":%s", cfg.Port)
	if err := http.ListenAndServe(bind, nil); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}
