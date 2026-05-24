package ws

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/vsuaiqq/cicd/shared/logger"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	readLimit  = 512
)

type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,

			CheckOrigin: func(_ *http.Request) bool { return true },
		},
	}
}

func (h *Handler) ServeGlobalWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L().Warn().Err(err).Msg("ws global upgrade error")
		return
	}

	c := h.hub.SubscribeGlobal()
	go writeGlobalPump(conn, c)
	readGlobalPump(conn, c, h.hub)
}

func writeGlobalPump(conn *websocket.Conn, c *globalClient) {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	defer conn.Close()

	for {
		select {
		case msg := <-c.send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.done:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			return
		}
	}
}

func readGlobalPump(conn *websocket.Conn, c *globalClient, hub *Hub) {
	defer hub.UnsubscribeGlobal(c)
	defer conn.Close()

	conn.SetReadLimit(readLimit)
	conn.SetReadDeadline(time.Now().Add(pingPeriod + pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L().Warn().Err(err).Str("run_id", runID).Msg("ws upgrade error")
		return
	}

	c, catchUp := h.hub.subscribeWithLastEvent(runID)

	if len(catchUp) > 0 {
		select {
		case c.send <- catchUp:
		default:
			logger.L().Warn().Str("run_id", runID).Msg("ws catch-up send dropped: buffer full")
		}
	}

	go writePump(conn, c)
	readPump(conn, c, h.hub)
}

func writePump(conn *websocket.Conn, c *client) {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	defer conn.Close()

	for {
		select {
		case msg := <-c.send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			return
		}
	}
}

func readPump(conn *websocket.Conn, c *client, hub *Hub) {
	defer hub.unsubscribe(c)
	defer conn.Close()

	conn.SetReadLimit(readLimit)

	conn.SetReadDeadline(time.Now().Add(pingPeriod + pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
