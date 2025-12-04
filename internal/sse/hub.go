package sse

import (
	
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Event struct {
	Type   string      `json:"type"`
	RID    string      `json:"rid,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	TS     int64       `json:"ts"`
}

type Hub struct {
	log    zerolog.Logger
	mu     sync.RWMutex
	clients map[chan []byte]struct{}
	in     chan []byte
	quit   chan struct{}
}

func NewHub(log zerolog.Logger) *Hub {
	return &Hub{
		log: log, clients: make(map[chan []byte]struct{}),
		in: make(chan []byte, 256), quit: make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case b := <-h.in:
			h.mu.RLock()
			for ch := range h.clients {
				select { case ch <- b: default: }
			}
			h.mu.RUnlock()
		case <-h.quit:
			return
		}
	}
}

func (h *Hub) Broadcast(v any) {
	b, _ := json.Marshal(v)
	h.in <- b
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := make(chan []byte, 32)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
	}()

	// Send an initial comment so some proxies keep the connection alive
	_, _ = w.Write([]byte(": connected\n\n"))
	flusher.Flush()

	// Periodic heartbeat to keep connections alive (15s)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()

	for {
		select {
		case b := <-ch:
			// SSE frame: "data: <json>\n\n"
			if _, err := w.Write([]byte("data: ")); err != nil {
				return
			}
			if _, err := w.Write(b); err != nil {
				return
			}
			if _, err := w.Write([]byte("\n\n")); err != nil {
				return
			}
			flusher.Flush()

		case <-ticker.C:
			// Heartbeat comment (doesn’t surface in EventSource.onmessage)
			if _, err := w.Write([]byte(": ping\n\n")); err != nil {
				return
			}
			flusher.Flush()

		case <-ctx.Done():
			// Client disconnected
			return
		}
	}
}


func (h *Hub) Close() { close(h.quit) }

func JSONEvent(t string, rid string, payload any, ts int64) map[string]any {
	return map[string]any{"type": t, "rid": rid, "payload": payload, "ts": ts}
}
