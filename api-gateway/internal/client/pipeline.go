package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var ErrNotFound = errors.New("not found")

type PipelineClient struct {
	baseURL     string
	httpClient  *http.Client
	internalKey string
}

func NewPipelineClient(baseURL string, timeout time.Duration, internalAPIKey string) *PipelineClient {
	return &PipelineClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpClient:  &http.Client{Timeout: timeout},
		internalKey: internalAPIKey,
	}
}

func (c *PipelineClient) setInternalHeaders(h http.Header) {
	if c.internalKey != "" {
		h.Set("X-Internal-API-Key", c.internalKey)
	}
}

func (c *PipelineClient) ListRuns(ctx context.Context, projectID string) ([]map[string]any, error) {
	endpoint := c.baseURL + "/api/runs?" + url.Values{"project_id": {projectID}}.Encode()
	var result []map[string]any
	if err := c.get(ctx, endpoint, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *PipelineClient) GetRun(ctx context.Context, runID string) (map[string]any, error) {
	endpoint := c.baseURL + "/api/runs/" + url.PathEscape(runID)
	var result map[string]any
	if err := c.get(ctx, endpoint, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type RunDetails struct {
	ID           string      `json:"id"`
	ProjectID    string      `json:"project_id"`
	Status       string      `json:"status"`
	PipelineYAML string      `json:"pipeline_yaml"`
	Jobs         []JobDetail `json:"jobs"`
}

type JobDetail struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Status      string       `json:"status"`
	Steps       []StepDetail `json:"steps"`
}

type StepDetail struct {
	Index     int    `json:"index"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	ExitCode  int    `json:"exit_code"`
	LogOutput string `json:"log_output"`
}

func (c *PipelineClient) GetRunDetails(ctx context.Context, runID string) (*RunDetails, error) {
	endpoint := c.baseURL + "/api/runs/" + url.PathEscape(runID)
	var result RunDetails
	if err := c.get(ctx, endpoint, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *PipelineClient) DialRunWS(ctx context.Context, runID string) (*websocket.Conn, error) {
	wsURL := c.wsBaseURL() + "/ws/runs/" + url.PathEscape(runID)
	dialer := websocket.Dialer{HandshakeTimeout: 45 * time.Second}
	hdr := http.Header{}
	c.setInternalHeaders(hdr)
	conn, _, err := dialer.DialContext(ctx, wsURL, hdr)
	if err != nil {
		return nil, fmt.Errorf("pipeline client: ws dial: %w", err)
	}
	return conn, nil
}

func (c *PipelineClient) DialGlobalWS(ctx context.Context) (*websocket.Conn, error) {
	wsURL := c.wsBaseURL() + "/ws/events"
	dialer := websocket.Dialer{HandshakeTimeout: 45 * time.Second}
	hdr := http.Header{}
	c.setInternalHeaders(hdr)
	conn, _, err := dialer.DialContext(ctx, wsURL, hdr)
	if err != nil {
		return nil, fmt.Errorf("pipeline client: global ws dial: %w", err)
	}
	return conn, nil
}

func (c *PipelineClient) wsBaseURL() string {
	if s, ok := strings.CutPrefix(c.baseURL, "https://"); ok {
		return "wss://" + s
	}
	if s, ok := strings.CutPrefix(c.baseURL, "http://"); ok {
		return "ws://" + s
	}
	return c.baseURL
}

func (c *PipelineClient) CancelRun(ctx context.Context, runID string) error {
	return c.post(ctx, c.baseURL+"/api/runs/"+url.PathEscape(runID)+"/cancel", nil)
}

func (c *PipelineClient) ApproveJob(ctx context.Context, runID, jobID string) error {
	return c.post(ctx, c.baseURL+"/api/runs/"+url.PathEscape(runID)+"/jobs/"+url.PathEscape(jobID)+"/approve", nil)
}

func (c *PipelineClient) RejectJob(ctx context.Context, runID, jobID string) error {
	return c.post(ctx, c.baseURL+"/api/runs/"+url.PathEscape(runID)+"/jobs/"+url.PathEscape(jobID)+"/reject", nil)
}

func (c *PipelineClient) post(ctx context.Context, endpoint string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return fmt.Errorf("pipeline client: build request: %w", err)
	}
	c.setInternalHeaders(req.Header)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pipeline client: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pipeline client: status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *PipelineClient) get(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("pipeline client: build request: %w", err)
	}
	c.setInternalHeaders(req.Header)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pipeline client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pipeline client: status %d: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
