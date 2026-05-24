package kafka

import (
	"context"

	"github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type GitEventConsumer struct {
	consumer *kafka.Consumer
	handler  kafka.Handler
}

func NewGitEventConsumer(
	consumer *kafka.Consumer,
	domainHandler GitEventHandler,
) *GitEventConsumer {
	return &GitEventConsumer{
		consumer: consumer,
		handler:  newGitEventKafkaHandler(domainHandler),
	}
}

func (c *GitEventConsumer) Start(
	ctx context.Context,
	topics []string,
) error {
	logger.L().Info().Strs("topics", topics).Msg("git-consumer starting")
	return c.consumer.Consume(ctx, topics, c.handler)
}

func (c *GitEventConsumer) Close() error {
	return c.consumer.Close()
}
