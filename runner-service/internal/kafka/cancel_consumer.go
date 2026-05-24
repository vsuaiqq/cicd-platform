package kafka

import (
	"context"
	"encoding/json"

	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type CancelHandler interface {
	HandleCancel(evt *events.CancelJobEvent)
}

type cancelKafkaHandler struct {
	handler CancelHandler
}

func (h *cancelKafkaHandler) Handle(_ context.Context, msg *sharedKafka.Message) error {
	var evt events.CancelJobEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		logger.L().Warn().Err(err).Msg("runner-cancel decode error")
		return sharedKafka.ErrDecodeFailed
	}
	h.handler.HandleCancel(&evt)
	return nil
}

type CancelConsumer struct {
	consumer *sharedKafka.Consumer
	handler  sharedKafka.Handler
}

func NewCancelConsumer(consumer *sharedKafka.Consumer, handler CancelHandler) *CancelConsumer {
	return &CancelConsumer{
		consumer: consumer,
		handler:  &cancelKafkaHandler{handler: handler},
	}
}

func (c *CancelConsumer) Start(ctx context.Context, topics []string) error {
	logger.L().Info().Strs("topics", topics).Msg("runner-cancel starting")
	return c.consumer.Consume(ctx, topics, c.handler)
}

func (c *CancelConsumer) Close() error {
	return c.consumer.Close()
}
