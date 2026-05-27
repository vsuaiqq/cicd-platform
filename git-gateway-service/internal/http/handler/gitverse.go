package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vsuaiqq/cicd/shared/events"
)

type GitVerseProcessor struct{}

func NewGitVerseProcessor() *GitVerseProcessor {
	return &GitVerseProcessor{}
}

func (p *GitVerseProcessor) ValidateAndParse(r *http.Request, secretKey []byte) (*events.GitEvent, error) {
	payload, err := readRequestBody(r)
	if err != nil {
		return nil, fmt.Errorf("gitverse: read body: %w", err)
	}

	sig := r.Header.Get("X-GitVerse-Signature")
	if sig == "" {
		sig = r.Header.Get("X-GitVerse-Signature-256")
	}
	if err := verifyHMACSHA256(payload, secretKey, sig, "sha256="); err != nil {
		return nil, fmt.Errorf("gitverse: signature validation failed: %w", err)
	}

	switch strings.ToLower(r.Header.Get("X-GitVerse-Event")) {
	case "push":
		return parseGitVersePush(payload)
	case "ping":
		return nil, nil
	default:
		return nil, nil
	}
}

type gitVersePush struct {
	Ref    string `json:"ref"`
	After  string `json:"after"`
	Before string `json:"before"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Username string `json:"username"`
			Name     string `json:"name"`
		} `json:"author"`
	} `json:"commits"`
	Repository struct {
		Name    string `json:"name"`
		SSHURL  string `json:"ssh_url"`
		CloneURL string `json:"clone_url"`
		Owner   struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Pusher struct {
		Login string `json:"login"`
	} `json:"pusher"`
}

func parseGitVersePush(payload []byte) (*events.GitEvent, error) {
	var body gitVersePush
	if err := json.Unmarshal(payload, &body); err != nil {
		return nil, fmt.Errorf("gitverse: parse push payload: %w", err)
	}

	branch, ok := strings.CutPrefix(body.Ref, "refs/heads/")
	if !ok || branch == "" {
		return nil, fmt.Errorf("gitverse: unsupported ref %q", body.Ref)
	}
	if body.After == "" || strings.HasPrefix(body.After, "0000000") {
		return nil, nil
	}

	commits := make([]*events.Commit, len(body.Commits))
	for i, c := range body.Commits {
		login := c.Author.Username
		if login == "" {
			login = c.Author.Name
		}
		commits[i] = &events.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author:  &events.User{Login: login},
			URL:     c.URL,
		}
	}

	repoURL := body.Repository.SSHURL
	if repoURL == "" {
		repoURL = body.Repository.CloneURL
	}

	author := body.Pusher.Login
	if author == "" {
		author = body.Repository.Owner.Login
	}

	return &events.GitEvent{
		Platform: events.GitVerse,
		Event:    events.Push,
		Repository: &events.Repository{
			Name: body.Repository.Name,
			URL:  repoURL,
		},
		Ref:       body.Ref,
		Branch:    branch,
		CommitSHA: body.After,
		Commits:   commits,
		Author:    &events.User{Login: author},
	}, nil
}
