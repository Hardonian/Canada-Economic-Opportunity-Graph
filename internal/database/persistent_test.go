package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestPersistentStore_CrashRecoveryAndWAL(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cog-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := context.Background()

	// 1. Open store, write records, then close
	store1, err := OpenPersistentStore(tempDir, true)
	if err != nil {
		t.Fatalf("failed to open persistent store: %v", err)
	}

	p1 := &domain.Project{
		ID:        "proj-persist-01",
		Name:      "Chibougamau Critical Lithium",
		Sector:    domain.SectorCriticalMinerals,
		Province:  "QC",
		CapexCAD:  1_800_000_000,
		UpdatedAt: time.Now().UTC(),
	}
	if err := store1.SaveProject(ctx, p1); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}

	ent1 := &domain.Entity{
		ID:        "ent-persist-01",
		LegalName: "Lithium Nordique Inc.",
	}
	if err := store1.SaveEntity(ctx, ent1); err != nil {
		t.Fatalf("failed to save entity: %v", err)
	}

	if err := store1.Close(); err != nil {
		t.Fatalf("failed to close store: %v", err)
	}

	// 2. Re-open store and assert crash recovery via WAL
	store2, err := OpenPersistentStore(tempDir, true)
	if err != nil {
		t.Fatalf("failed to re-open store: %v", err)
	}
	defer store2.Close()

	pRecovered, err := store2.GetProject(ctx, "proj-persist-01")
	if err != nil {
		t.Fatalf("failed to recover project from WAL: %v", err)
	}
	if pRecovered.Name != "Chibougamau Critical Lithium" {
		t.Errorf("expected project name %q, got %q", "Chibougamau Critical Lithium", pRecovered.Name)
	}

	entRecovered, err := store2.GetEntity(ctx, "ent-persist-01")
	if err != nil {
		t.Fatalf("failed to recover entity from WAL: %v", err)
	}
	if entRecovered.LegalName != "Lithium Nordique Inc." {
		t.Errorf("expected entity legal name %q, got %q", "Lithium Nordique Inc.", entRecovered.LegalName)
	}

	// 3. Test Checkpoint (snapshot creation + WAL truncation)
	if err := store2.Checkpoint(ctx); err != nil {
		t.Fatalf("failed to checkpoint: %v", err)
	}

	snapFile := filepath.Join(tempDir, "snapshot.json.gz")
	if fi, err := os.Stat(snapFile); err != nil || fi.Size() == 0 {
		t.Fatalf("expected non-empty snapshot file, err=%v", err)
	}

	// 4. Open a 3rd time to test recovery from snapshot
	store3, err := OpenPersistentStore(tempDir, true)
	if err != nil {
		t.Fatalf("failed to open store from snapshot: %v", err)
	}
	defer store3.Close()

	p3, err := store3.GetProject(ctx, "proj-persist-01")
	if err != nil || p3.ID != "proj-persist-01" {
		t.Fatalf("failed to load project from snapshot: %v", err)
	}
}
