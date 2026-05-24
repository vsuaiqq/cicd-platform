package kafka

import (
	"context"
	"encoding/json"

	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type JobResultHandler interface {
	HandleJobResult(ctx context.Context, result *events.PipelineJobResult) error
}

type jobResultKafkaHandler struct {
	handler JobResultHandler
}

func (h *jobResultKafkaHandler) Handle(ctx context.Context, msg *sharedKafka.Message) error {
	var result events.PipelineJobResult
	if err := json.Unmarshal(msg.Value, &result); err != nil {
		logger.L().Warn().Err(err).Msg("job-result-consumer decode error")
		return sharedKafka.ErrDecodeFailed
	}
	return h.handler.HandleJobResult(ctx, &result)
}

type JobResultConsumer struct {
	consumer *sharedKafka.Consumer
	handler  sharedKafka.Handler
}

func NewJobResultConsumer(consumer *sharedKafka.Consumer, handler JobResultHandler) *JobResultConsumer {
	return &JobResultConsumer{
		consumer: consumer,
		handler:  &jobResultKafkaHandler{handler: handler},
	}
}

func (c *JobResultConsumer) Start(ctx context.Context, topics []string) error {
	logger.L().Info().Strs("topics", topics).Msg("job-result-consumer starting")
	return c.consumer.Consume(ctx, topics, c.handler)
}

func (c *JobResultConsumer) Close() error {
	return c.consumer.Close()
}
