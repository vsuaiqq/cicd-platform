package config

import sharedConfig "github.com/vsuaiqq/cicd/shared/config"

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("POSTGRES_DSN", &cfg.Postgres.DSN)
	sharedConfig.EnvOverride("PROJECTS_ENCRYPTION_KEY", &cfg.Projects.EncryptionKey)
	sharedConfig.EnvOverride("WEBHOOK_BASE_URL", &cfg.Projects.WebhookBaseURL)
}

func Validate(cfg *Config) error {
	return sharedConfig.RequireNonEmpty(map[string]string{
		"POSTGRES_DSN":            cfg.Postgres.DSN,
		"PROJECTS_ENCRYPTION_KEY": cfg.Projects.EncryptionKey,
	})
}
