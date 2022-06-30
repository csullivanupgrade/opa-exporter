package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/csullivanupgrade/opa-exporter/internal/config"
	"github.com/csullivanupgrade/opa-exporter/pkg/opa"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startScheduled(e *opa.Exporter, cfg config.Config) {
	done := make(chan bool)
	ticker := time.NewTicker(cfg.Interval)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Tick at", t)
				constraints, err := opa.GetConstraints(&cfg.InCluster)
				if err != nil {
					log.Printf("%+v\n", err)
				}
				log.Printf("Found %v constraints", len(constraints))
				allMetrics := make([]prometheus.Metric, 0)
				violationMetrics := opa.ExportViolations(e.ConstraintViolation, constraints)
				allMetrics = append(allMetrics, violationMetrics...)

				constraintInformationMetrics := opa.ExportConstraintInformation(e.ConstraintInformation, constraints)
				allMetrics = append(allMetrics, constraintInformationMetrics...)

				e.Metrics = allMetrics
			}
		}
	}()
}

func Run(cfg config.Config) {
	exporter := opa.NewExporter(cfg)
	startScheduled(exporter, cfg)
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
			log.Printf("err handling response: %v", err)
		}
	})
	bind := fmt.Sprintf(":%s", cfg.Port)
	log.Fatal(http.ListenAndServe(bind, nil))
}
