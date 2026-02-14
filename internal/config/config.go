package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	SchedulerName   string
	Env             string
	K8sAPIServerURL string
	Namespace       string

	PollingInterval time.Duration
}

var v *viper.Viper

func init() {
	v = viper.New()

	v.SetDefault("scheduler.name", "custom-scheduler")
	v.SetDefault("env", "development")
	v.SetDefault("k8s.api-server-url", "http://localhost:8080")
	v.SetDefault("polling-interval", 1*time.Second)
	v.SetDefault("namespace", "custom-scheduler-namespace")
}

// Get returns the current configuration.
func Get() *Config {
	return &Config{
		SchedulerName: v.GetString("scheduler.name"),
		Env:           v.GetString("env"),
		K8sAPIServerURL: v.GetString("k8s.api-server-url"),
		PollingInterval: v.GetDuration("polling-interval"),
		Namespace: v.GetString("namespace"),
	}
}