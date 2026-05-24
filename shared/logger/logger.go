package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServiceName string

	Level string

	Pretty bool

	Output io.Writer
}

func Init(cfg Config) {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown"
	}

	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	var out io.Writer = cfg.Output
	if cfg.Pretty {
		out = zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: "15:04:05",
			PartsOrder: []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				"service",
				zerolog.MessageFieldName,
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(i.(string))
			},
		}
	}

	log.Logger = zerolog.New(out).
		Level(level).
		With().
		Str("service", cfg.ServiceName).
		Timestamp().
		Logger()
}

func parseLevel(s string) zerolog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return zerolog.DebugLevel
	case "info", "":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func L() *zerolog.Logger {
	return &log.Logger
}

func With() zerolog.Context {
	return log.Logger.With()
}

type contextKey struct{ name string }

var (
	TraceIDKey = contextKey{"trace_id"}
	SpanIDKey  = contextKey{"span_id"}
)

func FromContext(ctx context.Context) *zerolog.Logger {
	l := log.Logger
	if ctx == nil {
		return &l
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		l = l.With().Str("trace_id", traceID).Logger()
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		l = l.With().Str("span_id", spanID).Logger()
	}
	return &l
}
