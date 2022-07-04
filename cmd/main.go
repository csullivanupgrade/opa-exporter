package main

import (
	"context"
	"flag"

	"github.com/csullivanupgrade/opa-exporter/internal/config"
	"github.com/csullivanupgrade/opa-exporter/internal/log"
	"github.com/csullivanupgrade/opa-exporter/internal/server"
)

func main() {
	configFile := flag.String("config", "", "The path for the config file")
	flag.Parse()
	cfg := config.New(*configFile)

	logger := log.NewLogger(cfg.LogLevel, cfg.LogMode)

	ctx, cancel := context.WithCancel(log.SetContext(context.Background(), logger))

	server.Run(ctx, *cfg)
	defer cancel()
}
