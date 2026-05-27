package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vsuaiqq/cicd/shared/events"
)

type BitbucketProcessor struct{}

func NewBitbucketProcessor() *BitbucketProcessor {
	return &BitbucketProcessor{}
}

func (p *BitbucketProcessor) ValidateAndParse(r *http.Request, secretKey []byte) (*events.GitEvent, error) {
	payload, err := readRequestBody(r)
	if err != nil {
		return nil, fmt.Errorf("bitbucket: read body: %w", err)
	}

	sig := r.Header.Get("X-Hub-Signature-256")
	if sig == "" {
		sig = r.Header.Get("X-Hub-Signature")
	}
	if err := verifyHMACSHA256(payload, secretKey, sig, "sha256="); err != nil {
		return nil, fmt.Errorf("bitbucket: signature validation failed: %w", err)
	}

	eventKey := r.Header.Get("X-Event-Key")
	switch eventKey {
	case "repo:push":
		return parseBitbucketPush(payload)
	case "diagnostics:ping":
		return nil, nil
	default:
		return nil, nil
	}
}

type bitbucketPush struct {
	Push struct {
		Changes []struct {
			New struct {
				Type   string `json:"type"`
				Name   string `json:"name"`
				Target struct {
					Hash string `json:"hash"`
				} `json:"target"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Links    struct {
			Clone []struct {
				Name string `json:"name"`
				Href string `json:"href"`
			} `json:"clone"`
		} `json:"links"`
	} `json:"repository"`
	Actor struct {
		DisplayName string `json:"display_name"`
		Username    string `json:"username"`
	} `json:"actor"`
}

func parseBitbucketPush(payload []byte) (*events.GitEvent, error) {
	var body bitbucketPush
	if err := json.Unmarshal(payload, &body); err != nil {
		return nil, fmt.Errorf("bitbucket: parse push payload: %w", err)
	}

	var branch, commitSHA string
	for _, change := range body.Push.Changes {
		if change.New.Type != "branch" {
			continue
		}
		branch = change.New.Name
		commitSHA = change.New.Target.Hash
		break
	}
	if branch == "" || commitSHA == "" {
		return nil, fmt.Errorf("bitbucket: push payload missing branch commit")
	}

	repoURL := bitbucketSSHURL(body.Repository.Links.Clone)
	if repoURL == "" && body.Repository.FullName != "" {
		repoURL = "git@bitbucket.org:" + body.Repository.FullName + ".git"
	}

	author := body.Actor.Username
	if author == "" {
		author = body.Actor.DisplayName
	}

	return &events.GitEvent{
		Platform: events.Bitbucket,
		Event:    events.Push,
		Repository: &events.Repository{
			Name: body.Repository.Name,
			URL:  repoURL,
		},
		Ref:       "refs/heads/" + branch,
		Branch:    branch,
		CommitSHA: commitSHA,
		Author:    &events.User{Login: author},
	}, nil
}

func bitbucketSSHURL(clones []struct {
	Name string `json:"name"`
	Href string `json:"href"`
}) string {
	for _, c := range clones {
		if c.Name == "ssh" && c.Href != "" {
			return c.Href
		}
	}
	for _, c := range clones {
		if strings.HasPrefix(c.Href, "git@") {
			return c.Href
		}
	}
	return ""
}
