package config

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type ServiceLoader[T any] struct {
	ServiceName    string
	ConfigPath     string
	ApplyEnv       func(*T)
	ApplyDefaults  func(*T)
	Validate       func(*T) error
}

func LoadService[T any](ctx context.Context, loader ServiceLoader[T]) (*T, error) {
	if loader.ServiceName != "" {
		if err := BootstrapSecrets(ctx, loader.ServiceName); err != nil {
			return nil, fmt.Errorf("vault bootstrap: %w", err)
		}
	}

	cfg, err := Load[T](loader.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}

	if loader.ApplyEnv != nil {
		loader.ApplyEnv(cfg)
	}
	if loader.ApplyDefaults != nil {
		loader.ApplyDefaults(cfg)
	}
	if loader.Validate != nil {
		if err := loader.Validate(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func EnvOverrideBrokers(key string, brokers *[]string) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) > 0 {
		*brokers = out
	}
}
