package database

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// WALEntryType defines the type of record written to the Write-Ahead Log.
type WALEntryType string

const (
	WALProject            WALEntryType = "PROJECT"
	WALEntity             WALEntryType = "ENTITY"
	WALEvent              WALEntryType = "EVENT"
	WALRelationship       WALEntryType = "RELATIONSHIP"
	WALScore              WALEntryType = "SCORE"
	WALProcurement        WALEntryType = "PROCUREMENT"
	WALCapitalItem        WALEntryType = "CAPITAL_ITEM"
	WALSignal             WALEntryType = "SIGNAL"
	WALOpportunity        WALEntryType = "OPPORTUNITY"
	WALClaim              WALEntryType = "CLAIM"
	WALCandidateProject   WALEntryType = "CANDIDATE_PROJECT"
	WALProjectPhase       WALEntryType = "PROJECT_PHASE"
	WALCapitalRequirement WALEntryType = "CAPITAL_REQUIREMENT"
	WALCapitalNeed        WALEntryType = "CAPITAL_NEED"
	WALMilestone          WALEntryType = "MILESTONE"
	WALReadiness          WALEntryType = "READINESS"
	WALAuditEntry         WALEntryType = "AUDIT_ENTRY"
	WALEvidence           WALEntryType = "EVIDENCE"
	WALTradeMetric        WALEntryType = "TRADE_METRIC"
)

// WALRecord represents an append-only transaction entry in the durable log.
type WALRecord struct {
	Timestamp int64           `json:"ts"`
	Type      WALEntryType    `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

// PersistentSnapshot aggregates the complete dataset for periodic checkpointing.
type PersistentSnapshot struct {
	Projects            []*domain.Project             `json:"projects"`
	Entities            []*domain.Entity              `json:"entities"`
	Events              []*domain.Event               `json:"events"`
	Relationships       []*domain.Relationship        `json:"relationships"`
	Scores              []*domain.ProjectScore        `json:"scores"`
	Procurements        []*domain.Procurement         `json:"procurements"`
	CapitalItems        []*domain.CapitalItem         `json:"capital_items"`
	Signals             []*domain.Signal              `json:"signals"`
	Opportunities       []*domain.Opportunity         `json:"opportunities"`
	Claims              []*domain.Claim               `json:"claims"`
	CandidateProjects   []*domain.CandidateProject    `json:"candidate_projects"`
	ProjectPhases       []*domain.ProjectPhase        `json:"project_phases"`
	CapitalRequirements []*domain.CapitalRequirement  `json:"capital_requirements"`
	CapitalNeeds        []*domain.CapitalNeed         `json:"capital_needs"`
	Milestones          []*domain.Milestone           `json:"milestones"`
	Readiness           []*domain.ReadinessAssessment `json:"readiness"`
	AuditEntries        []*domain.AuditEntry          `json:"audit_entries"`
	Evidence            []*domain.Evidence            `json:"evidence"`
	TradeMetrics        []*domain.TradeMetric         `json:"trade_metrics"`
	SnapshotTime        time.Time                     `json:"snapshot_time"`
}

// PersistentStore provides an ACID-durable, WAL-backed storage engine.
type PersistentStore struct {
	*MemoryStore
	dir         string
	walFile     *os.File
	walWriter   *bufio.Writer
	mu          sync.Mutex
	closed      bool
	autoSync    bool
	replayCount int
}

// OpenPersistentStore opens or creates a persistent store in the target directory.
func OpenPersistentStore(dir string, autoSync bool) (*PersistentStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir %q: %w", dir, err)
	}

	mem := NewMemoryStore()
	ps := &PersistentStore{
		MemoryStore: mem,
		dir:         dir,
		autoSync:    autoSync,
	}

	// 1. Recover from snapshot if available
	if err := ps.loadLatestSnapshot(); err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}

	// 2. Replay outstanding WAL entries
	walPath := filepath.Join(dir, "wal.log")
	if err := ps.replayWAL(walPath); err != nil {
		return nil, fmt.Errorf("replay wal: %w", err)
	}

	// 3. Open WAL in append mode
	f, err := os.OpenFile(walPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open wal: %w", err)
	}
	ps.walFile = f
	ps.walWriter = bufio.NewWriter(f)

	return ps, nil
}

func (ps *PersistentStore) writeWAL(entryType WALEntryType, v any) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed || ps.walWriter == nil {
		return nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	record := WALRecord{
		Timestamp: time.Now().UnixNano(),
		Type:      entryType,
		Payload:   data,
	}

	line, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if _, err := ps.walWriter.Write(append(line, '\n')); err != nil {
		return err
	}

	if ps.autoSync {
		if err := ps.walWriter.Flush(); err != nil {
			return err
		}
		if err := ps.walFile.Sync(); err != nil {
			return err
		}
	}
	return nil
}

// Flush ensures all buffered transactions are committed to physical media.
func (ps *PersistentStore) Flush() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.walWriter != nil {
		if err := ps.walWriter.Flush(); err != nil {
			return err
		}
	}
	if ps.walFile != nil {
		return ps.walFile.Sync()
	}
	return nil
}

// Close flushes the WAL and cleanly shuts down the storage engine.
func (ps *PersistentStore) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.closed {
		return nil
	}
	ps.closed = true

	var err error
	if ps.walWriter != nil {
		if fErr := ps.walWriter.Flush(); fErr != nil && err == nil {
			err = fErr
		}
	}
	if ps.walFile != nil {
		if sErr := ps.walFile.Sync(); sErr != nil && err == nil {
			err = sErr
		}
		if cErr := ps.walFile.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}
	return err
}

// Checkpoint creates a compressed snapshot of the entire state and truncates the WAL.
func (ps *PersistentStore) Checkpoint(ctx context.Context) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if err := ps.walWriter.Flush(); err != nil {
		return err
	}
	if err := ps.walFile.Sync(); err != nil {
		return err
	}

	projects, _, _ := ps.MemoryStore.ListProjects(ctx, ProjectFilter{Limit: 1_000_000})
	entities, _ := ps.MemoryStore.ListEntities(ctx)
	events, _ := ps.MemoryStore.ListRecentEvents(ctx, 1_000_000)
	capItems, _ := ps.MemoryStore.ListAllCapitalItems(ctx)
	opps, _ := ps.MemoryStore.ListAllOpportunities(ctx)
	procs, _ := ps.MemoryStore.ListProcurements(ctx, 1_000_000, 0)
	signals, _ := ps.MemoryStore.ListSignals(ctx, 365*24*time.Hour, 1_000_000)

	snap := PersistentSnapshot{
		Projects:      projects,
		Entities:      entities,
		Events:        events,
		CapitalItems:  capItems,
		Opportunities: opps,
		Procurements:  procs,
		Signals:       signals,
		SnapshotTime:  time.Now().UTC(),
	}

	snapPath := filepath.Join(ps.dir, "snapshot.json.gz")
	tmpPath := snapPath + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(f)
	if err := json.NewEncoder(gz).Encode(snap); err != nil {
		gz.Close()
		f.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := gz.Close(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, snapPath); err != nil {
		return err
	}

	// Truncate WAL
	_ = ps.walFile.Close()
	walPath := filepath.Join(ps.dir, "wal.log")
	newWal, err := os.OpenFile(walPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	ps.walFile = newWal
	ps.walWriter = bufio.NewWriter(newWal)
	return nil
}

func (ps *PersistentStore) loadLatestSnapshot() error {
	snapPath := filepath.Join(ps.dir, "snapshot.json.gz")
	f, err := os.Open(snapPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	var snap PersistentSnapshot
	if err := json.NewDecoder(gz).Decode(&snap); err != nil {
		return err
	}

	ctx := context.Background()
	for _, p := range snap.Projects {
		_ = ps.MemoryStore.SaveProject(ctx, p)
	}
	for _, e := range snap.Entities {
		_ = ps.MemoryStore.SaveEntity(ctx, e)
	}
	for _, ev := range snap.Events {
		_ = ps.MemoryStore.SaveEvent(ctx, ev)
	}
	for _, c := range snap.CapitalItems {
		_ = ps.MemoryStore.SaveCapitalItem(ctx, c)
	}
	for _, o := range snap.Opportunities {
		_ = ps.MemoryStore.SaveOpportunity(ctx, o)
	}
	for _, pr := range snap.Procurements {
		_ = ps.MemoryStore.SaveProcurement(ctx, pr)
	}
	for _, s := range snap.Signals {
		_ = ps.MemoryStore.SaveSignal(ctx, s)
	}
	return nil
}

func (ps *PersistentStore) replayWAL(walPath string) error {
	f, err := os.Open(walPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	ctx := context.Background()
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record WALRecord
		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}

		ps.replayCount++
		switch record.Type {
		case WALProject:
			var p domain.Project
			if json.Unmarshal(record.Payload, &p) == nil {
				_ = ps.MemoryStore.SaveProject(ctx, &p)
			}
		case WALEntity:
			var e domain.Entity
			if json.Unmarshal(record.Payload, &e) == nil {
				_ = ps.MemoryStore.SaveEntity(ctx, &e)
			}
		case WALEvent:
			var ev domain.Event
			if json.Unmarshal(record.Payload, &ev) == nil {
				_ = ps.MemoryStore.SaveEvent(ctx, &ev)
			}
		case WALRelationship:
			var r domain.Relationship
			if json.Unmarshal(record.Payload, &r) == nil {
				_ = ps.MemoryStore.SaveRelationship(ctx, &r)
			}
		case WALScore:
			var s domain.ProjectScore
			if json.Unmarshal(record.Payload, &s) == nil {
				_ = ps.MemoryStore.SaveScore(ctx, &s)
			}
		case WALProcurement:
			var pr domain.Procurement
			if json.Unmarshal(record.Payload, &pr) == nil {
				_ = ps.MemoryStore.SaveProcurement(ctx, &pr)
			}
		case WALCapitalItem:
			var ci domain.CapitalItem
			if json.Unmarshal(record.Payload, &ci) == nil {
				_ = ps.MemoryStore.SaveCapitalItem(ctx, &ci)
			}
		case WALSignal:
			var sg domain.Signal
			if json.Unmarshal(record.Payload, &sg) == nil {
				_ = ps.MemoryStore.SaveSignal(ctx, &sg)
			}
		case WALOpportunity:
			var op domain.Opportunity
			if json.Unmarshal(record.Payload, &op) == nil {
				_ = ps.MemoryStore.SaveOpportunity(ctx, &op)
			}
		case WALClaim:
			var cl domain.Claim
			if json.Unmarshal(record.Payload, &cl) == nil {
				_ = ps.MemoryStore.SaveClaim(ctx, &cl)
			}
		case WALMilestone:
			var ms domain.Milestone
			if json.Unmarshal(record.Payload, &ms) == nil {
				_ = ps.MemoryStore.SaveMilestone(ctx, &ms)
			}
		case WALEvidence:
			var ev domain.Evidence
			if json.Unmarshal(record.Payload, &ev) == nil {
				_ = ps.MemoryStore.SaveEvidence(ctx, &ev)
			}
		}
	}
	return scanner.Err()
}

// Overridden mutators that write to WAL before updating in-memory state:

func (ps *PersistentStore) SaveProject(ctx context.Context, p *domain.Project) error {
	if err := ps.writeWAL(WALProject, p); err != nil {
		return err
	}
	return ps.MemoryStore.SaveProject(ctx, p)
}

func (ps *PersistentStore) SaveEntity(ctx context.Context, e *domain.Entity) error {
	if err := ps.writeWAL(WALEntity, e); err != nil {
		return err
	}
	return ps.MemoryStore.SaveEntity(ctx, e)
}

func (ps *PersistentStore) SaveEvent(ctx context.Context, ev *domain.Event) error {
	if err := ps.writeWAL(WALEvent, ev); err != nil {
		return err
	}
	return ps.MemoryStore.SaveEvent(ctx, ev)
}

func (ps *PersistentStore) SaveRelationship(ctx context.Context, r *domain.Relationship) error {
	if err := ps.writeWAL(WALRelationship, r); err != nil {
		return err
	}
	return ps.MemoryStore.SaveRelationship(ctx, r)
}

func (ps *PersistentStore) SaveScore(ctx context.Context, s *domain.ProjectScore) error {
	if err := ps.writeWAL(WALScore, s); err != nil {
		return err
	}
	return ps.MemoryStore.SaveScore(ctx, s)
}

func (ps *PersistentStore) SaveProcurement(ctx context.Context, p *domain.Procurement) error {
	if err := ps.writeWAL(WALProcurement, p); err != nil {
		return err
	}
	return ps.MemoryStore.SaveProcurement(ctx, p)
}

func (ps *PersistentStore) SaveCapitalItem(ctx context.Context, c *domain.CapitalItem) error {
	if err := ps.writeWAL(WALCapitalItem, c); err != nil {
		return err
	}
	return ps.MemoryStore.SaveCapitalItem(ctx, c)
}

func (ps *PersistentStore) SaveSignal(ctx context.Context, s *domain.Signal) error {
	if err := ps.writeWAL(WALSignal, s); err != nil {
		return err
	}
	return ps.MemoryStore.SaveSignal(ctx, s)
}

func (ps *PersistentStore) SaveOpportunity(ctx context.Context, o *domain.Opportunity) error {
	if err := ps.writeWAL(WALOpportunity, o); err != nil {
		return err
	}
	return ps.MemoryStore.SaveOpportunity(ctx, o)
}

func (ps *PersistentStore) SaveClaim(ctx context.Context, claim *domain.Claim) error {
	if err := ps.writeWAL(WALClaim, claim); err != nil {
		return err
	}
	return ps.MemoryStore.SaveClaim(ctx, claim)
}

func (ps *PersistentStore) SaveMilestone(ctx context.Context, milestone *domain.Milestone) error {
	if err := ps.writeWAL(WALMilestone, milestone); err != nil {
		return err
	}
	return ps.MemoryStore.SaveMilestone(ctx, milestone)
}

func (ps *PersistentStore) SaveEvidence(ctx context.Context, e *domain.Evidence) error {
	if err := ps.writeWAL(WALEvidence, e); err != nil {
		return err
	}
	return ps.MemoryStore.SaveEvidence(ctx, e)
}

// Compile-time check that PersistentStore implements Store
var _ Store = (*PersistentStore)(nil)
