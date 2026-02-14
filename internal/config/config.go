package config

import "github.com/spf13/viper"

// Config holds application configuration.
type Config struct {
	SchedulerName string
	Env           string
}

var v *viper.Viper

func init() {
	v = viper.New()

	v.SetDefault("scheduler.name", "custom-scheduler")
	v.SetDefault("env", "development")
}

// Get returns the current configuration.
func Get() Config {
	return Config{
		SchedulerName: v.GetString("scheduler.name"),
		Env:           v.GetString("env"),
	}
}