package config

import "time"

type Config struct {
	Server              ServerConfig              `yaml:"server"`
	AuthService         AuthServiceConfig         `yaml:"auth_service"`
	ProjectsService     ProjectsServiceConfig     `yaml:"projects_service"`
	OrchestratorService OrchestratorServiceConfig `yaml:"orchestrator_service"`
	AIService           AIServiceConfig           `yaml:"ai_service"`
	AnalyticsService    AnalyticsServiceConfig    `yaml:"analytics_service"`
	InternalAPIKey      string                    `yaml:"-"`
}

type ServerConfig struct {
	Port           int      `yaml:"port"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type AuthServiceConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}

type ProjectsServiceConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}

type OrchestratorServiceConfig struct {
	BaseURL string        `yaml:"base_url"`
	Timeout time.Duration `yaml:"timeout"`
}

type AIServiceConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}

type AnalyticsServiceConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}
