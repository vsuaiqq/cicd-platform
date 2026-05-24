package config

import sharedConfig "github.com/vsuaiqq/cicd/shared/config"

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("POSTGRES_DSN", &cfg.Postgres.DSN)
	sharedConfig.EnvOverride("JWT_ACCESS_SECRET", &cfg.JWT.AccessTokenSecret)
	sharedConfig.EnvOverride("JWT_REFRESH_SECRET", &cfg.JWT.RefreshTokenSecret)
	sharedConfig.EnvOverride("REDIS_ADDR", &cfg.Redis.Addr)
	sharedConfig.EnvOverride("REDIS_PASSWORD", &cfg.Redis.Password)
}

func Validate(cfg *Config) error {
	return sharedConfig.RequireNonEmpty(map[string]string{
		"POSTGRES_DSN":         cfg.Postgres.DSN,
		"JWT_ACCESS_SECRET":    cfg.JWT.AccessTokenSecret,
		"JWT_REFRESH_SECRET":   cfg.JWT.RefreshTokenSecret,
	})
}
