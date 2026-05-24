package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/kafka"
)

var ErrNilEvent = errors.New("git event publisher: event must not be nil")

type GitEventPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewGitEventPublisher(producer *kafka.Producer, topic string) *GitEventPublisher {
	return &GitEventPublisher{producer: producer, topic: topic}
}

func (p *GitEventPublisher) Publish(ctx context.Context, event *events.GitEvent) error {
	if event == nil {
		return ErrNilEvent
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("git event publisher: marshal failed: %w", err)
	}

	msg := kafka.Message{
		Key:   keyForEvent(event),
		Value: payload,
	}

	if err := p.producer.Send(ctx, p.topic, &msg); err != nil {
		return fmt.Errorf("git event publisher: kafka send failed: %w", err)
	}

	return nil
}

func keyForEvent(event *events.GitEvent) []byte {
	if event.Repository != nil && event.Repository.URL != "" {
		return []byte(event.Repository.URL)
	}
	if event.CommitSHA != "" {
		return []byte(event.CommitSHA)
	}
	return nil
}
