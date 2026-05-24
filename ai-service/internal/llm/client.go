package llm

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	model      string
	maxTokens  int
	httpClient *http.Client

	tokenMu     sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

const oauthBaseURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

func New(baseURL, apiKey, model string, maxTokens int, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		baseURL:   strings.TrimRight(baseURL, "/"),
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) Complete(ctx context.Context, messages []Message) (string, error) {
	token, err := c.bearerToken(ctx)
	if err != nil {
		return "", fmt.Errorf("llm: get token: %w", err)
	}

	type request struct {
		Model     string    `json:"model"`
		Messages  []Message `json:"messages"`
		MaxTokens int       `json:"max_tokens"`
		Stream    bool      `json:"stream"`
	}
	type choice struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	}
	type apiError struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	type response struct {
		Choices []choice  `json:"choices"`
		Error   *apiError `json:"error,omitempty"`
		Status  int       `json:"status,omitempty"`
		Message string    `json:"message,omitempty"`
	}

	payload := request{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: c.maxTokens,
		Stream:    false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("llm: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llm: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: api call failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm: read response: %w", err)
	}

	var result response
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("llm: parse response (status %d): %w — body: %s",
			resp.StatusCode, err, truncateStr(string(data), 300))
	}

	if resp.StatusCode != http.StatusOK {
		msg := result.Message
		if result.Error != nil && result.Error.Message != "" {
			msg = result.Error.Message
		}
		if msg == "" {
			msg = truncateStr(string(data), 200)
		}
		return "", fmt.Errorf("llm: status %d: %s", resp.StatusCode, msg)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("llm: empty response (no choices returned)")
	}

	content := strings.TrimSpace(result.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("llm: model returned empty content (finish_reason: %s)",
			result.Choices[0].FinishReason)
	}

	return content, nil
}

func (c *Client) bearerToken(ctx context.Context) (string, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-30*time.Second)) {
		return c.accessToken, nil
	}

	token, expiry, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.accessToken = token
	c.tokenExpiry = expiry
	return c.accessToken, nil
}

func (c *Client) fetchToken(ctx context.Context) (string, time.Time, error) {
	form := url.Values{"scope": {"GIGACHAT_API_PERS"}}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthBaseURL,
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("build oauth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+c.apiKey)
	req.Header.Set("RqUID", newUUID())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("oauth request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("read oauth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("oauth status %d: %s", resp.StatusCode, truncateStr(string(data), 200))
	}

	var payload struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", time.Time{}, fmt.Errorf("parse oauth response: %w", err)
	}
	if payload.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("oauth returned empty access_token")
	}

	expiry := time.UnixMilli(payload.ExpiresAt)
	return payload.AccessToken, expiry, nil
}

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
