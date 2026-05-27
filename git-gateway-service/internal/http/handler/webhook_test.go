package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vsuaiqq/cicd/shared/events"
)

func signPayload(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestBitbucketProcessor_pushEvent(t *testing.T) {
	secret := []byte("bitbucket-secret")
	payload := []byte(`{
		"push": {
			"changes": [{
				"new": {
					"type": "branch",
					"name": "main",
					"target": {"hash": "abc123def456"}
				}
			}]
		},
		"repository": {
			"name": "demo",
			"full_name": "team/demo",
			"links": {
				"clone": [{"name": "ssh", "href": "git@bitbucket.org:team/demo.git"}]
			}
		},
		"actor": {"username": "alice"}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook/bitbucket/p1", bytes.NewReader(payload))
	req.Header.Set("X-Event-Key", "repo:push")
	req.Header.Set("X-Hub-Signature-256", signPayload(secret, payload))

	ev, err := NewBitbucketProcessor().ValidateAndParse(req, secret)
	if err != nil {
		t.Fatalf("ValidateAndParse: %v", err)
	}
	if ev.Platform != events.Bitbucket {
		t.Fatalf("platform = %q", ev.Platform)
	}
	if ev.Branch != "main" || ev.CommitSHA != "abc123def456" {
		t.Fatalf("branch/sha = %s/%s", ev.Branch, ev.CommitSHA)
	}
	if ev.Repository.URL != "git@bitbucket.org:team/demo.git" {
		t.Fatalf("repo url = %q", ev.Repository.URL)
	}
}

func TestBitbucketProcessor_ping(t *testing.T) {
	secret := []byte("bitbucket-secret")
	payload := []byte(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	req.Header.Set("X-Event-Key", "diagnostics:ping")
	req.Header.Set("X-Hub-Signature-256", signPayload(secret, payload))

	ev, err := NewBitbucketProcessor().ValidateAndParse(req, secret)
	if err != nil {
		t.Fatalf("ValidateAndParse: %v", err)
	}
	if ev != nil {
		t.Fatalf("expected nil event, got %+v", ev)
	}
}

func TestBitbucketProcessor_badSignature(t *testing.T) {
	payload := []byte(`{"push":{"changes":[]}}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	req.Header.Set("X-Event-Key", "repo:push")
	req.Header.Set("X-Hub-Signature-256", "sha256=deadbeef")

	_, err := NewBitbucketProcessor().ValidateAndParse(req, []byte("secret"))
	if err == nil {
		t.Fatal("expected signature error")
	}
}

func TestGitVerseProcessor_pushEvent(t *testing.T) {
	secret := []byte("gitverse-secret")
	payload := []byte(`{
		"ref": "refs/heads/main",
		"after": "ff00112233445566778899aabbccddeeff00112233",
		"repository": {
			"name": "demo",
			"ssh_url": "git@gitverse.ru:org/demo.git"
		},
		"pusher": {"login": "bob"},
		"commits": [{
			"id": "ff00112233445566778899aabbccddeeff00112233",
			"message": "init",
			"author": {"username": "bob"}
		}]
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook/gitverse/p1", bytes.NewReader(payload))
	req.Header.Set("X-GitVerse-Event", "push")
	req.Header.Set("X-GitVerse-Signature", signPayload(secret, payload))

	ev, err := NewGitVerseProcessor().ValidateAndParse(req, secret)
	if err != nil {
		t.Fatalf("ValidateAndParse: %v", err)
	}
	if ev.Platform != events.GitVerse {
		t.Fatalf("platform = %q", ev.Platform)
	}
	if ev.Branch != "main" {
		t.Fatalf("branch = %q", ev.Branch)
	}
	if ev.Author.Login != "bob" {
		t.Fatalf("author = %q", ev.Author.Login)
	}
}

func TestGitVerseProcessor_ping(t *testing.T) {
	secret := []byte("gitverse-secret")
	payload := []byte(`{"zen":"ok"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	req.Header.Set("X-GitVerse-Event", "ping")
	req.Header.Set("X-GitVerse-Signature", signPayload(secret, payload))

	ev, err := NewGitVerseProcessor().ValidateAndParse(req, secret)
	if err != nil {
		t.Fatalf("ValidateAndParse: %v", err)
	}
	if ev != nil {
		t.Fatal("expected nil event")
	}
}
