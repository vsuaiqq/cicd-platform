package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type WSProxyHandler struct {
	authClient     *client.AuthClient
	pipelineClient *client.PipelineClient
	projectsClient *client.ProjectsClient
	upgrader       websocket.Upgrader
}

func NewWSProxyHandler(pipelineClient *client.PipelineClient, authClient *client.AuthClient, projectsClient *client.ProjectsClient, allowedOrigins []string) *WSProxyHandler {
	return &WSProxyHandler{
		authClient:     authClient,
		pipelineClient: pipelineClient,
		projectsClient: projectsClient,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     makeOriginChecker(allowedOrigins),
		},
	}
}

func makeOriginChecker(allowedOrigins []string) func(*http.Request) bool {
	for _, o := range allowedOrigins {
		if o == "*" {
			return func(_ *http.Request) bool { return true }
		}
	}
	if len(allowedOrigins) == 0 {
		return func(_ *http.Request) bool { return true }
	}
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}
	return func(r *http.Request) bool {
		_, ok := allowed[r.Header.Get("Origin")]
		return ok
	}
}

type wsRunEvent struct {
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`
}

func (h *WSProxyHandler) ServeGlobalWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	resp, err := h.authClient.ValidateToken(r.Context(), token)
	if err != nil || !resp.Valid || resp.UserId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	lp, err := h.projectsClient.ListProjects(r.Context(), resp.UserId)
	if err != nil {
		logger.L().Warn().Err(err).Msg("ws-proxy list projects failed")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	allowed := make(map[string]struct{}, len(lp.Projects))
	for _, p := range lp.Projects {
		if p.Id != "" {
			allowed[p.Id] = struct{}{}
		}
	}

	upstreamConn, err := h.pipelineClient.DialGlobalWS(r.Context())
	if err != nil {
		logger.L().Warn().Err(err).Msg("ws-proxy global upstream dial error")
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}
	defer upstreamConn.Close()

	clientConn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L().Warn().Err(err).Msg("ws-proxy global client upgrade error")
		return
	}
	defer clientConn.Close()

	errc := make(chan error, 2)

	go func() {
		for {
			msgType, msg, err := upstreamConn.ReadMessage()
			if err != nil {
				errc <- err
				return
			}
			if msgType == websocket.TextMessage {
				var ev wsRunEvent
				if json.Unmarshal(msg, &ev) == nil && ev.Type == "run_updated" {
					if ev.ProjectID == "" {
						continue
					}
					if _, ok := allowed[ev.ProjectID]; !ok {
						continue
					}
				}
			}
			if err := clientConn.WriteMessage(msgType, msg); err != nil {
				errc <- err
				return
			}
		}
	}()

	go func() {
		for {
			msgType, msg, err := clientConn.ReadMessage()
			if err != nil {
				errc <- err
				return
			}
			if err := upstreamConn.WriteMessage(msgType, msg); err != nil {
				errc <- err
				return
			}
		}
	}()

	<-errc
}

func (h *WSProxyHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	resp, err := h.authClient.ValidateToken(r.Context(), token)
	if err != nil || !resp.Valid || resp.UserId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	run, err := h.pipelineClient.GetRun(r.Context(), runID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	projID, ok := projectIDFromRunMap(run)
	if !ok {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, resp.UserId, projID) {
		return
	}

	upstreamConn, err := h.pipelineClient.DialRunWS(r.Context(), runID)
	if err != nil {
		logger.L().Warn().Err(err).Str("run_id", runID).Msg("ws-proxy upstream dial error")
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}
	defer upstreamConn.Close()

	clientConn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L().Warn().Err(err).Str("run_id", runID).Msg("ws-proxy client upgrade error")
		return
	}
	defer clientConn.Close()

	errc := make(chan error, 2)

	go func() {
		for {
			msgType, msg, err := upstreamConn.ReadMessage()
			if err != nil {
				errc <- err
				return
			}
			if err := clientConn.WriteMessage(msgType, msg); err != nil {
				errc <- err
				return
			}
		}
	}()

	go func() {
		for {
			msgType, msg, err := clientConn.ReadMessage()
			if err != nil {
				errc <- err
				return
			}
			if err := upstreamConn.WriteMessage(msgType, msg); err != nil {
				errc <- err
				return
			}
		}
	}()

	<-errc
}
