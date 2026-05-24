package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
)

type ResultPublisher struct {
	producer *sharedKafka.Producer
	topic    string
}

func NewResultPublisher(producer *sharedKafka.Producer, topic string) *ResultPublisher {
	return &ResultPublisher{producer: producer, topic: topic}
}

func (p *ResultPublisher) Publish(ctx context.Context, result *events.PipelineJobResult) error {
	if result == nil {
		return fmt.Errorf("result publisher: result is nil")
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("result publisher: marshal: %w", err)
	}
	msg := &sharedKafka.Message{
		Key:   []byte(result.JobID),
		Value: payload,
	}
	return p.producer.Send(ctx, p.topic, msg)
}
