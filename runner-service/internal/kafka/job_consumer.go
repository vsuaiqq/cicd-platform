package kafka

import (
	"context"
	"encoding/json"

	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type JobHandler interface {
	HandleJob(ctx context.Context, task *events.PipelineJobTask) error
}

type jobKafkaHandler struct {
	handler JobHandler
}

func (h *jobKafkaHandler) Handle(ctx context.Context, msg *sharedKafka.Message) error {
	var task events.PipelineJobTask
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		logger.L().Warn().Err(err).Msg("runner-consumer decode error")
		return sharedKafka.ErrDecodeFailed
	}
	return h.handler.HandleJob(ctx, &task)
}

type JobConsumer struct {
	consumer *sharedKafka.Consumer
	handler  sharedKafka.Handler
}

func NewJobConsumer(consumer *sharedKafka.Consumer, handler JobHandler) *JobConsumer {
	return &JobConsumer{
		consumer: consumer,
		handler:  &jobKafkaHandler{handler: handler},
	}
}

func (c *JobConsumer) Start(ctx context.Context, topics []string) error {
	logger.L().Info().Strs("topics", topics).Msg("runner-consumer starting")
	return c.consumer.Consume(ctx, topics, c.handler)
}

func (c *JobConsumer) Close() error {
	return c.consumer.Close()
}
