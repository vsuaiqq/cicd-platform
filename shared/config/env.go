package config

import "os"

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func EnvOverride(key string, dst *string) {
	if v := os.Getenv(key); v != "" {
		*dst = v
	}
}
