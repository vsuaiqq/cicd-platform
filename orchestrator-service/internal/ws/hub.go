package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/vsuaiqq/cicd/shared/logger"
)

type EventType string

const (
	EventRunUpdated          EventType = "run_updated"
	EventJobUpdated          EventType = "job_updated"
	EventJobAwaitingApproval EventType = "job_awaiting_approval"

	EventHeartbeat EventType = "heartbeat"
)

const heartbeatInterval = time.Second

const terminalEventTTL = 24 * time.Hour

const transitionalEventTTL = time.Hour

type StepEvent struct {
	Index      int    `json:"index"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	ExitCode   int    `json:"exit_code"`
	LogOutput  string `json:"log_output"`
	StartedAt  int64  `json:"started_at"`
	FinishedAt int64  `json:"finished_at"`
}

type Event struct {
	Type  EventType `json:"type"`
	RunID string    `json:"run_id"`

	ProjectID string `json:"project_id,omitempty"`
	Status    string `json:"status,omitempty"`
	JobName   string `json:"job_name,omitempty"`
	JobID     string `json:"job_id,omitempty"`

	Steps []StepEvent `json:"steps,omitempty"`

	ServerTimeMs int64 `json:"server_time_ms,omitempty"`

	RunFinishedAtMs int64 `json:"run_finished_at_ms,omitempty"`

	JobStartedAtMs int64 `json:"job_started_at_ms,omitempty"`

	JobFinishedAtMs int64 `json:"job_finished_at_ms,omitempty"`
}

type client struct {
	runID string
	send  chan []byte
	done  chan struct{}
}

func newClient(runID string) *client {
	return &client{
		runID: runID,
		send:  make(chan []byte, 64),
		done:  make(chan struct{}),
	}
}

type globalClient struct {
	send chan []byte
	done chan struct{}
}

func newGlobalClient() *globalClient {
	return &globalClient{
		send: make(chan []byte, 64),
		done: make(chan struct{}),
	}
}

type cachedEvent struct {
	data      []byte
	expiresAt time.Time
}

var terminalStatuses = map[string]bool{
	"success":   true,
	"failed":    true,
	"cancelled": true,
}

type Hub struct {
	mu            sync.RWMutex
	clients       map[string]map[*client]struct{}
	lastEvent     map[string]cachedEvent
	activeRuns    map[string]struct{}
	globalClients map[*globalClient]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[string]map[*client]struct{}),
		lastEvent:     make(map[string]cachedEvent),
		activeRuns:    make(map[string]struct{}),
		globalClients: make(map[*globalClient]struct{}),
	}
}

func (h *Hub) SubscribeGlobal() *globalClient {
	c := newGlobalClient()
	h.mu.Lock()
	h.globalClients[c] = struct{}{}
	h.mu.Unlock()
	return c
}

func (h *Hub) UnsubscribeGlobal(c *globalClient) {
	h.mu.Lock()
	delete(h.globalClients, c)
	h.mu.Unlock()
	close(c.done)
}

func (h *Hub) Start(ctx context.Context) {
	go h.heartbeatLoop(ctx)
	go h.gcLoop(ctx)
}

func (h *Hub) MarkActive(runID string) {
	h.mu.Lock()
	h.activeRuns[runID] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) subscribeWithLastEvent(runID string) (*client, []byte) {
	c := newClient(runID)
	h.mu.Lock()
	if h.clients[runID] == nil {
		h.clients[runID] = make(map[*client]struct{})
	}
	h.clients[runID][c] = struct{}{}
	var catchUp []byte
	if entry, ok := h.lastEvent[runID]; ok {
		catchUp = entry.data
	}
	h.mu.Unlock()
	return c, catchUp
}

func (h *Hub) unsubscribe(c *client) {
	h.mu.Lock()
	delete(h.clients[c.runID], c)
	if len(h.clients[c.runID]) == 0 {
		delete(h.clients, c.runID)
	}
	h.mu.Unlock()
	close(c.done)
}

func (h *Hub) Broadcast(runID string, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		logger.L().Warn().Err(err).Msg("ws marshal error")
		return
	}

	isTerminal := event.Type == EventRunUpdated && terminalStatuses[event.Status]
	ttl := transitionalEventTTL
	if isTerminal {
		ttl = terminalEventTTL
	}

	h.mu.Lock()
	h.lastEvent[runID] = cachedEvent{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
	if isTerminal {
		delete(h.activeRuns, runID)
	}

	recipients := make([]*client, 0, len(h.clients[runID]))
	for c := range h.clients[runID] {
		recipients = append(recipients, c)
	}
	var globalRecipients []*globalClient
	if event.Type == EventRunUpdated {
		globalRecipients = make([]*globalClient, 0, len(h.globalClients))
		for c := range h.globalClients {
			globalRecipients = append(globalRecipients, c)
		}
	}
	h.mu.Unlock()

	for _, c := range recipients {
		select {
		case c.send <- data:
		case <-c.done:
		default:
			logger.L().Warn().Str("run_id", runID).Msg("ws slow client, dropping message")
		}
	}
	for _, c := range globalRecipients {
		select {
		case c.send <- data:
		case <-c.done:
		default:

		}
	}
}

func (h *Hub) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			h.broadcastHeartbeats(now)
		}
	}
}

func (h *Hub) broadcastHeartbeats(now time.Time) {
	type runClients struct {
		runID   string
		clients []*client
	}

	h.mu.RLock()
	targets := make([]runClients, 0, len(h.activeRuns))
	for runID := range h.activeRuns {
		subs := h.clients[runID]
		if len(subs) == 0 {
			continue
		}
		cs := make([]*client, 0, len(subs))
		for c := range subs {
			cs = append(cs, c)
		}
		targets = append(targets, runClients{runID: runID, clients: cs})
	}
	h.mu.RUnlock()

	if len(targets) == 0 {
		return
	}

	tsMs := now.UnixMilli()
	for _, t := range targets {
		data, _ := json.Marshal(Event{
			Type:         EventHeartbeat,
			RunID:        t.runID,
			ServerTimeMs: tsMs,
		})
		for _, c := range t.clients {
			select {
			case c.send <- data:
			case <-c.done:
			default:

			}
		}
	}
}

func (h *Hub) gcLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			h.evictStaleEvents(now)
		}
	}
}

func (h *Hub) evictStaleEvents(now time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for runID, entry := range h.lastEvent {
		if now.After(entry.expiresAt) && len(h.clients[runID]) == 0 {
			delete(h.lastEvent, runID)
		}
	}
}
