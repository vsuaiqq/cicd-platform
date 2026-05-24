package kafka

import (
	"context"

	"github.com/vsuaiqq/cicd/shared/events"
)

type GitEventHandler interface {
	HandleGitEvent(ctx context.Context, event *events.GitEvent) error
}
