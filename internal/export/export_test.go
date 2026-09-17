package export

import (
	"context"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestExportProjectBundle(t *testing.T) {
	store := database.NewMemoryStore()
	ctx := context.Background()
	now := time.Now().UTC()

	proj := &domain.Project{
		ID:        "test-proj-01",
		Name:      "Test Project Bundle",
		Sector:    domain.SectorCriticalMinerals,
		Province:  "ON",
		CapexCAD:  2_500_000_000,
		UpdatedAt: now,
	}
	if err := store.SaveProject(ctx, proj); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}

	bundle, err := ExportProjectBundle(ctx, store, "test-proj-01")
	if err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}
	if bundle.Project.ID != "test-proj-01" {
		t.Fatalf("expected project ID test-proj-01, got %s", bundle.Project.ID)
	}
	if bundle.ExportedAt.IsZero() {
		t.Fatal("expected non-zero ExportedAt")
	}

	md := bundle.ToMarkdown()
	if len(md) == 0 {
		t.Fatal("expected non-empty markdown export")
	}
}
