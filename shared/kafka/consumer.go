package kafka

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/vsuaiqq/cicd/shared/logger"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

func NewConsumer(cfg Config) (*Consumer, error) {
	scfg := sarama.NewConfig()

	scfg.ClientID = cfg.ClientID

	scfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	scfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	scfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, scfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{group: group}, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	topics []string,
	handler Handler,
) error {
	for {
		if err := c.group.Consume(ctx, topics, &consumerGroupHandler{
			handler: handler,
		}); err != nil {
			logger.L().Warn().Err(err).Msg("kafka consume error")
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}
