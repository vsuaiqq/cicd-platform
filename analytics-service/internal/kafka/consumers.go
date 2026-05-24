package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vsuaiqq/cicd/analytics-service/internal/db"
	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type RunEventConsumer struct {
	consumer *sharedKafka.Consumer
	repo     *db.Repository
}

func NewRunEventConsumer(consumer *sharedKafka.Consumer, repo *db.Repository) *RunEventConsumer {
	return &RunEventConsumer{consumer: consumer, repo: repo}
}

func (c *RunEventConsumer) Start(ctx context.Context, topics []string) error {
	return c.consumer.Consume(ctx, topics, c)
}

func (c *RunEventConsumer) Close() error {
	return c.consumer.Close()
}

func (c *RunEventConsumer) Handle(ctx context.Context, msg *sharedKafka.Message) error {
	var evt events.RunCompletedEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		logger.L().Warn().Err(err).Msg("run event decode error")
		return sharedKafka.ErrDecodeFailed
	}

	createdAt := time.Unix(evt.FinishedAt, 0)
	if evt.FinishedAt == 0 {
		createdAt = time.Now()
	}

	dur := evt.DurationSec
	if dur < 0 {
		dur = 0
	}
	if err := c.repo.InsertRunEvent(ctx, &db.RunEvent{
		ProjectID:   evt.ProjectID,
		RunID:       evt.RunID,
		Status:      evt.Status,
		Branch:      evt.Branch,
		DurationSec: uint32(dur),
		CreatedAt:   createdAt,
	}); err != nil {
		logger.L().Error().Err(err).Msg("insert run event error")
		return err
	}
	logger.L().Info().
		Str("run_id", evt.RunID).Str("status", evt.Status).Str("project_id", evt.ProjectID).
		Msg("recorded run")
	return nil
}

type JobResultConsumer struct {
	consumer *sharedKafka.Consumer
	repo     *db.Repository
}

func NewJobResultConsumer(consumer *sharedKafka.Consumer, repo *db.Repository) *JobResultConsumer {
	return &JobResultConsumer{consumer: consumer, repo: repo}
}

func (c *JobResultConsumer) Start(ctx context.Context, topics []string) error {
	return c.consumer.Consume(ctx, topics, c)
}

func (c *JobResultConsumer) Close() error {
	return c.consumer.Close()
}

func (c *JobResultConsumer) Handle(ctx context.Context, msg *sharedKafka.Message) error {
	var result events.PipelineJobResult
	if err := json.Unmarshal(msg.Value, &result); err != nil {
		logger.L().Warn().Err(err).Msg("job result decode error")
		return sharedKafka.ErrDecodeFailed
	}

	if result.Status != "success" && result.Status != "failed" && result.Status != "cancelled" {
		return nil
	}

	dur := 0
	if result.FinishedAt > result.StartedAt && result.StartedAt > 0 {
		dur = int(result.FinishedAt - result.StartedAt)
	}

	attempt := result.AttemptNumber
	if attempt < 1 {
		attempt = 1
	}

	createdAt := time.Unix(result.FinishedAt, 0)
	if result.FinishedAt == 0 {
		createdAt = time.Now()
	}

	if dur < 0 {
		dur = 0
	}
	if err := c.repo.InsertJobEvent(ctx, &db.JobEvent{
		RunID:       result.RunID,
		JobID:       result.JobID,
		ProjectID:   result.ProjectID,
		JobName:     result.JobName,
		Status:      result.Status,
		DurationSec: uint32(dur),
		Attempt:     uint8(attempt),
		CreatedAt:   createdAt,
	}); err != nil {
		logger.L().Error().Err(err).Msg("insert job event error")
		return err
	}
	return nil
}
