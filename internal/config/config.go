package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	port      = "METRICS_PORT"
	path      = "METRICS_PATH"
	inCluster = "INCLUSTER"
	interval  = "INTERVAL"
	namespace = "NAMESPACE"
)

type Config struct {
	Port      string        `mapstructure:"port"`
	Path      string        `mapstructure:"path"`
	InCluster bool          `mapstructure:"inCluster"`
	Interval  time.Duration `mapstructure:"interval"`
	Namespace string        `mapstructure:"namespace"`
}

// New returns a new Config struct
func New(configFile string) *Config {
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("OPA_EXPORTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// viper.AutomaticEnv()

	viper.SetDefault(port, "9141")
	viper.SetDefault(path, "metrics")
	viper.SetDefault(inCluster, true)
	viper.SetDefault(interval, "60s")
	viper.SetDefault(namespace, "default")

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
	fmt.Println("-- namepsace:\t", c.Namespace)

	return c
}
