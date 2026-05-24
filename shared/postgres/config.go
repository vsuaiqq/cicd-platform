package postgres

import "time"

type Config struct {
	DSN string

	MaxOpenConns int
	MaxIdleConns int

	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	ConnectTimeout time.Duration
}
