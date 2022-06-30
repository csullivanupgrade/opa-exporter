package main

import (
	"flag"

	"github.com/csullivanupgrade/opa-exporter/internal/config"
	"github.com/csullivanupgrade/opa-exporter/internal/server"
)

func main() {
	configFile := flag.String("config", "", "The path for the config file")
	flag.Parse()
	cfg := config.New(*configFile)
	server.Run(*cfg)
}
