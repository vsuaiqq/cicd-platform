package kafka

import "errors"

var (
	ErrHandlerFailed = errors.New("kafka handler failed")
	ErrDecodeFailed  = errors.New("kafka decode failed")
)
