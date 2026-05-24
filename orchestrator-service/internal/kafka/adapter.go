package kafka

import (
	"context"
	"encoding/json"

	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type gitEventKafkaHandler struct {
	handler GitEventHandler
}

func newGitEventKafkaHandler(h GitEventHandler) kafka.Handler {
	return &gitEventKafkaHandler{handler: h}
}

func (h *gitEventKafkaHandler) Handle(
	ctx context.Context,
	msg *kafka.Message,
) error {
	var event events.GitEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.L().Warn().Err(err).Msg("git-consumer decode error")
		return kafka.ErrDecodeFailed
	}

	return h.handler.HandleGitEvent(ctx, &event)
}
