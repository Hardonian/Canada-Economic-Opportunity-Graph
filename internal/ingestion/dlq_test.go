package ingestion

import (
	"testing"
)

func TestDeadLetterQueue(t *testing.T) {
	dlq := NewDeadLetterQueue(3)

	e1 := dlq.Enqueue("canadabuys", "TENDER", "missing field tender_id", "{'foo': 'bar'}")
	if e1 == nil || e1.ID == "" {
		t.Fatalf("expected valid entry, got nil")
	}

	dlq.Enqueue("iaac", "PROJECT", "invalid coordinates", "{'lat': 'north'}")
	dlq.Enqueue("cmhc", "HOUSING", "schema version mismatch", "{'units': -5}")

	if dlq.Count() != 3 {
		t.Fatalf("expected count 3, got %d", dlq.Count())
	}

	// Overfill: test eviction of oldest
	e4 := dlq.Enqueue("nrcan", "MINING", "malformed JSON", "bad json")
	if dlq.Count() != 3 {
		t.Fatalf("expected bounded count 3, got %d", dlq.Count())
	}

	list := dlq.List(10)
	if len(list) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(list))
	}

	if list[2].ID != e4.ID {
		t.Errorf("expected newest entry to be e4, got %s", list[2].ID)
	}

	dlq.Clear()
	if dlq.Count() != 0 {
		t.Errorf("expected 0 entries after Clear(), got %d", dlq.Count())
	}
}
