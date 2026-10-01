// Package eventsse implements high-throughput Server-Sent Events (SSE) streaming
// for real-time project milestone transitions, regulatory notices, and momentum signals.
package eventsse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/google/uuid"
)

// Subscriber represents an active SSE client stream.
type Subscriber struct {
	ID        string
	Channel   chan *domain.Event
	Filter    func(*domain.Event) bool
	CreatedAt time.Time
}

// Hub manages active SSE subscriptions and broadcasts domain events.
type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]*Subscriber
}

// NewHub initializes an event streaming hub.
func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[string]*Subscriber),
	}
}

// DefaultHub is the singleton instance used across the API server.
var DefaultHub = NewHub()

// Subscribe registers a new subscriber channel with an optional filter.
func (h *Hub) Subscribe(filter func(*domain.Event) bool) *Subscriber {
	h.mu.Lock()
	defer h.mu.Unlock()

	sub := &Subscriber{
		ID:        uuid.NewString(),
		Channel:   make(chan *domain.Event, 64),
		Filter:    filter,
		CreatedAt: time.Now().UTC(),
	}
	h.subscribers[sub.ID] = sub
	return sub
}

// Unsubscribe removes and closes a subscriber channel.
func (h *Hub) Unsubscribe(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sub, ok := h.subscribers[id]; ok {
		delete(h.subscribers, id)
		close(sub.Channel)
	}
}

// Broadcast distributes a new event to all active matching subscribers.
// Non-blocking write: if a subscriber's channel buffer is full, it drops the event
// to prevent lagging clients from stalling other consumers.
func (h *Hub) Broadcast(event *domain.Event) {
	if event == nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, sub := range h.subscribers {
		if sub.Filter != nil && !sub.Filter(event) {
			continue
		}
		select {
		case sub.Channel <- event:
		default:
			// Buffer full, drop to protect memory
		}
	}
}

// SubscriberCount returns the current number of active listeners.
func (h *Hub) SubscriberCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.subscribers)
}

// StreamHandler handles incoming HTTP GET requests for text/event-stream.
func (h *Hub) StreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by client", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	projectFilter := r.URL.Query().Get("project_id")
	typeFilter := r.URL.Query().Get("event_type")

	sub := h.Subscribe(func(ev *domain.Event) bool {
		if projectFilter != "" && ev.ProjectID != projectFilter {
			return false
		}
		if typeFilter != "" && ev.EventType != typeFilter {
			return false
		}
		return true
	})
	defer h.Unsubscribe(sub.ID)

	// Send initial connection establishment event
	fmt.Fprintf(w, "event: connected\ndata: {\"subscriber_id\":\"%s\",\"status\":\"STREAMING\"}\n\n", sub.ID)
	flusher.Flush()

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-heartbeat.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case ev, open := <-sub.Channel:
			if !open {
				return
			}
			data, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: project_event\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// BroadcastSynthetic generates a test event for demonstration or E2E validation.
func BroadcastSynthetic(projectID, title, eventType string) {
	DefaultHub.Broadcast(&domain.Event{
		ID:         "ev-" + uuid.NewString()[:8],
		ProjectID:  projectID,
		Title:      title,
		EventType:  eventType,
		EventDate:  time.Now().UTC(),
		EvidenceID: "ev-stream-provenance",
		CreatedAt:  time.Now().UTC(),
	})
}
