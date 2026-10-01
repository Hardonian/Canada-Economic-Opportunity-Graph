package ingestion

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// DeadLetterEntry represents an unparseable, schema-drifted, or invalid ingested record.
type DeadLetterEntry struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	RecordType  string    `json:"record_type"`
	Timestamp   time.Time `json:"timestamp"`
	ErrorReason string    `json:"error_reason"`
	RawPayload  string    `json:"raw_payload,omitempty"`
	RetryCount  int       `json:"retry_count"`
}

// DeadLetterQueue provides thread-safe quarantine for failed pipeline records.
type DeadLetterQueue struct {
	mu       sync.RWMutex
	capacity int
	entries  []*DeadLetterEntry
}

// Global default dead letter queue.
var DefaultDLQ = NewDeadLetterQueue(1000)

// NewDeadLetterQueue creates a bounded dead letter queue.
func NewDeadLetterQueue(capacity int) *DeadLetterQueue {
	if capacity <= 0 {
		capacity = 500
	}
	return &DeadLetterQueue{
		capacity: capacity,
		entries:  make([]*DeadLetterEntry, 0, capacity),
	}
}

// Enqueue adds a quarantined record to the DLQ.
func (q *DeadLetterQueue) Enqueue(source, recordType, reason, rawPayload string) *DeadLetterEntry {
	q.mu.Lock()
	defer q.mu.Unlock()

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", source, recordType, reason, time.Now().UnixNano())))
	id := "dlq-" + hex.EncodeToString(hash[:8])

	entry := &DeadLetterEntry{
		ID:          id,
		Source:      source,
		RecordType:  recordType,
		Timestamp:   time.Now().UTC(),
		ErrorReason: reason,
		RawPayload:  rawPayload,
		RetryCount:  0,
	}

	if len(q.entries) >= q.capacity {
		// Evict oldest
		q.entries = q.entries[1:]
	}
	q.entries = append(q.entries, entry)
	return entry
}

// List returns the most recent quarantined records up to limit.
func (q *DeadLetterQueue) List(limit int) []*DeadLetterEntry {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if limit <= 0 || limit > len(q.entries) {
		limit = len(q.entries)
	}

	result := make([]*DeadLetterEntry, limit)
	start := len(q.entries) - limit
	copy(result, q.entries[start:])
	return result
}

// Count returns the total number of currently quarantined records.
func (q *DeadLetterQueue) Count() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.entries)
}

// Clear flushes the dead letter queue.
func (q *DeadLetterQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = q.entries[:0]
}
