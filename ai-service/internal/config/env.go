package config

import (
	"time"

	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
)

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("AI_LLM_API_KEY", &cfg.LLM.APIKey)
	sharedConfig.EnvOverride("AI_LLM_BASE_URL", &cfg.LLM.BaseURL)
	sharedConfig.EnvOverride("AI_LLM_MODEL", &cfg.LLM.Model)
}

func ApplyDefaults(cfg *Config) {
	if cfg.LLM.BaseURL == "" {
		cfg.LLM.BaseURL = "https://gigachat.devices.sberbank.ru/api/v1"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "GigaChat"
	}
	if cfg.LLM.MaxTokens == 0 {
		cfg.LLM.MaxTokens = 1024
	}
	if cfg.LLM.Timeout == 0 {
		cfg.LLM.Timeout = 90 * time.Second
	}
	if cfg.Server.HTTPPort == 0 {
		cfg.Server.HTTPPort = 8086
	}
	if cfg.Server.GRPCPort == 0 {
		cfg.Server.GRPCPort = 50053
	}
}
