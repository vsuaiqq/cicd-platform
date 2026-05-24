package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/kafka"
)

var ErrNilTask = errors.New("job publisher: task must not be nil")

type CancelPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewCancelPublisher(producer *kafka.Producer, topic string) *CancelPublisher {
	return &CancelPublisher{producer: producer, topic: topic}
}

func (p *CancelPublisher) PublishCancel(ctx context.Context, runID, jobID string) error {
	payload, err := json.Marshal(events.CancelJobEvent{RunID: runID, JobID: jobID})
	if err != nil {
		return fmt.Errorf("cancel publisher: marshal: %w", err)
	}
	key := runID
	if jobID != "" {
		key = jobID
	}
	return p.producer.Send(ctx, p.topic, &kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
}

type RunEventPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewRunEventPublisher(producer *kafka.Producer, topic string) *RunEventPublisher {
	return &RunEventPublisher{producer: producer, topic: topic}
}

func (p *RunEventPublisher) PublishRunCompleted(ctx context.Context, evt *events.RunCompletedEvent) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("run event publisher: marshal: %w", err)
	}
	return p.producer.Send(ctx, p.topic, &kafka.Message{
		Key:   []byte(evt.RunID),
		Value: payload,
	})
}

type JobPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewJobPublisher(producer *kafka.Producer, topic string) *JobPublisher {
	return &JobPublisher{producer: producer, topic: topic}
}

func (p *JobPublisher) PublishJob(ctx context.Context, task *events.PipelineJobTask) error {
	if task == nil {
		return ErrNilTask
	}
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("job publisher: marshal: %w", err)
	}
	return p.producer.Send(ctx, p.topic, &kafka.Message{
		Key:   []byte(task.JobID),
		Value: payload,
	})
}
