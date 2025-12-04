package config

import "os"

type Config struct {
	Port        string
	ServiceName string
	Env         string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		ServiceName: env("SERVICE_NAME", "flash-api"),
		Env:         env("ENV", "dev"),
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
