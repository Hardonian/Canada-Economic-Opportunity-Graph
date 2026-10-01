package eventsse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestHub_SubscribeBroadcastUnsubscribe(t *testing.T) {
	hub := NewHub()

	sub := hub.Subscribe(func(ev *domain.Event) bool {
		return ev.ProjectID == "crawford-nickel"
	})

	if hub.SubscriberCount() != 1 {
		t.Fatalf("expected 1 subscriber, got %d", hub.SubscriberCount())
	}

	evMatching := &domain.Event{
		ID:        "ev-1",
		ProjectID: "crawford-nickel",
		Title:     "Permit Approved",
		EventType: "regulatory_filing",
		EventDate: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	evNonMatching := &domain.Event{
		ID:        "ev-2",
		ProjectID: "other-project",
		Title:     "Other Event",
		EventType: "stage_change",
		EventDate: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	hub.Broadcast(evMatching)
	hub.Broadcast(evNonMatching)

	select {
	case received := <-sub.Channel:
		if received.ID != "ev-1" {
			t.Fatalf("expected ev-1, got %s", received.ID)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for event")
	}

	// Non-matching event should not be in channel
	select {
	case unexpected := <-sub.Channel:
		t.Fatalf("received unexpected filtered event: %v", unexpected)
	default:
		// success
	}

	hub.Unsubscribe(sub.ID)
	if hub.SubscriberCount() != 0 {
		t.Fatalf("expected 0 subscribers, got %d", hub.SubscriberCount())
	}
}

func TestHub_StreamHandler(t *testing.T) {
	hub := NewHub()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/events?project_id=crawford-nickel", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		hub.StreamHandler(rec, req)
	}()

	// Wait for subscription to register
	time.Sleep(50 * time.Millisecond)
	if hub.SubscriberCount() != 1 {
		t.Fatalf("expected 1 subscriber, got %d", hub.SubscriberCount())
	}

	hub.Broadcast(&domain.Event{
		ID:        "ev-crawford-stream",
		ProjectID: "crawford-nickel",
		Title:     "Construction Notice",
		EventType: "stage_change",
		EventDate: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	})

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := rec.Body.String()
	if !strings.Contains(body, "event: connected") {
		t.Fatalf("expected connected event in stream: %s", body)
	}
	if !strings.Contains(body, "ev-crawford-stream") {
		t.Fatalf("expected ev-crawford-stream in stream: %s", body)
	}
}
