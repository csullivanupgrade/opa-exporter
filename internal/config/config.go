// Package config handles getting configuration from the environment and/or from file
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	port         = "METRICS_PORT"
	path         = "METRICS_PATH"
	inCluster    = "INCLUSTER"
	interval     = "INTERVAL"
	logLevel     = "LOG_LEVEL"
	logMode      = "LOG_MODE"
	namespace    = "NAMESPACE"
	writeTimeout = "WRITE_TIMEOUT"
	readTimeout  = "READ_TIMEOUT"
)

type Config struct {
	Port         string        `mapstructure:"port"`
	Path         string        `mapstructure:"path"`
	InCluster    bool          `mapstructure:"inCluster"`
	Interval     time.Duration `mapstructure:"interval"`
	LogLevel     string        `mapstructure:"logLevel"`
	LogMode      string        `mapstructure:"logMode"`
	Namespace    string        `mapstructure:"namespace"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

// New returns a new Config struct
func New(configFile string) *Config {
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("OPA_EXPORTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetDefault(port, "9141")
	viper.SetDefault(path, "metrics")
	viper.SetDefault(inCluster, true)
	viper.SetDefault(interval, "60s")
	viper.SetDefault(logLevel, "info")
	viper.SetDefault(logMode, "production")
	viper.SetDefault(namespace, "default")
	viper.SetDefault(readTimeout, "5s")
	viper.SetDefault(writeTimeout, "10s")

	c := &Config{
		Port:      viper.GetString(port),
		Path:      viper.GetString(path),
		InCluster: viper.GetBool(inCluster),
		Interval:  viper.GetDuration(interval),
	}

	if configFile != "" {
		fmt.Println("Initializing server using config file:", configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Could not open config file:", err)
			os.Exit(1)
		}
		if err := viper.Unmarshal(c); err != nil {
			fmt.Println("Could not read config file:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Using configurations:")
	fmt.Println("-- pram:\t", c.Port)
	fmt.Println("-- path:\t", c.Path)
	fmt.Println("-- incluster:\t", c.InCluster)
	fmt.Println("-- interval:\t", c.Interval)
	fmt.Println("-- logLevel:\t", c.LogLevel)
	fmt.Println("-- logMode:\t", c.LogMode)
	fmt.Println("-- namespace:\t", c.Namespace)
	fmt.Println("-- readTimeout:\t", c.ReadTimeout)
	fmt.Println("-- writeTimeout:\t", c.WriteTimeout)

	return c
}
