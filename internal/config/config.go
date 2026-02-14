package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	SchedulerName     string
	Env               string
	K8sAPIServerURL   string
	K8sAuthTokenPath  string // optional; when set, use for in-cluster auth (Bearer token)
	Namespace         string
	PollingInterval   time.Duration
}

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("scheduler.name", "custom-scheduler")
	v.SetDefault("env", "development")
	v.SetDefault("k8s.api-server-url", "http://localhost:8080")
	v.SetDefault("k8s.auth-token-path", "")
	v.SetDefault("polling-interval", 1*time.Second)
	v.SetDefault("namespace", "custom-scheduler-namespace")

	_ = v.BindEnv("k8s.api-server-url", "K8S_API_SERVER_URL")
	_ = v.BindEnv("k8s.auth-token-path", "K8S_AUTH_TOKEN_PATH")
	_ = v.BindEnv("scheduler.name", "SCHEDULER_NAME")
	_ = v.BindEnv("namespace", "NAMESPACE")
}

// Get returns the current configuration.
func Get() *Config {
	return &Config{
		SchedulerName:    v.GetString("scheduler.name"),
		Env:              v.GetString("env"),
		K8sAPIServerURL:  v.GetString("k8s.api-server-url"),
		K8sAuthTokenPath:  v.GetString("k8s.auth-token-path"),
		PollingInterval:  v.GetDuration("polling-interval"),
		Namespace:        v.GetString("namespace"),
	}
}