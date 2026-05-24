package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	handler Handler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		kmsg := &Message{
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: convertHeaders(msg.Headers),
		}

		err := h.handler.Handle(context.Background(), kmsg)
		if err != nil {

			continue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
