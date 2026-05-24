package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"

	"github.com/vsuaiqq/cicd/shared/events"
)

type GitHubProcessor struct{}

func NewGitHubProcessor() *GitHubProcessor {
	return &GitHubProcessor{}
}

func (p *GitHubProcessor) ValidateAndParse(r *http.Request, secretKey []byte) (*events.GitEvent, error) {
	payload, err := github.ValidatePayload(r, secretKey)
	if err != nil {
		return nil, fmt.Errorf("github: signature validation failed: %w", err)
	}

	raw, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, fmt.Errorf("github: failed to parse webhook: %w", err)
	}

	switch event := raw.(type) {
	case *github.PushEvent:
		return parsePushEvent(event), nil
	case *github.PingEvent:

		return nil, nil
	default:

		return nil, nil
	}
}

func parsePushEvent(event *github.PushEvent) *events.GitEvent {
	commits := make([]*events.Commit, len(event.Commits))
	for i, c := range event.Commits {
		commits[i] = &events.Commit{
			SHA:     c.GetSHA(),
			Message: c.GetMessage(),
			Author:  &events.User{Login: c.GetAuthor().GetLogin()},
			URL:     c.GetURL(),
		}
	}

	branch, _ := strings.CutPrefix(event.GetRef(), "refs/heads/")

	gitEvent := &events.GitEvent{
		Platform: events.GitHub,
		Event:    events.Push,
		Repository: &events.Repository{
			Name: event.GetRepo().GetName(),
			URL:  event.GetRepo().GetSSHURL(),
		},
		Ref:       event.GetRef(),
		Branch:    branch,
		CommitSHA: event.GetAfter(),
		Commits:   commits,
		Author:    &events.User{Login: event.GetPusher().GetLogin()},
	}

	if ts := event.GetHeadCommit().GetTimestamp(); !ts.IsZero() {
		gitEvent.Timestamp = ts.Unix()
	}

	return gitEvent
}
