package kafka

import "time"

type Config struct {
	Brokers  []string
	ClientID string

	GroupID string

	RetryMax     int
	RetryBackoff time.Duration
}
