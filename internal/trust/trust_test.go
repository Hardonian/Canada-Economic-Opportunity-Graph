package trust

import (
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestEvaluate_EmptyEvidence(t *testing.T) {
	proj := &domain.Project{ID: "p-1"}
	rep := Evaluate(proj, nil, nil, time.Now().UTC())
	if rep.EvidenceCount != 0 {
		t.Fatalf("expected 0 evidence count, got %d", rep.EvidenceCount)
	}
	if rep.Status != domain.StatusUnavailable {
		t.Fatalf("expected status unavailable, got %s", rep.Status)
	}
}

func TestEvaluate_WithEvidence(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-10 * 24 * time.Hour)
	proj := &domain.Project{
		ID:       "p-1",
		CapexCAD: 1_000_000,
	}

	ev1 := &domain.Evidence{
		ID:                 "ev-1",
		SourceTier:         domain.SourceTier1,
		Confidence:         domain.ConfidenceVerified,
		RetrievalTimestamp: past,
		EffectiveDate:      &past,
		ContentHash:        "hash1",
	}
	ev2 := &domain.Evidence{
		ID:                 "ev-2",
		SourceTier:         domain.SourceTier2,
		Confidence:         domain.ConfidenceReported,
		RetrievalTimestamp: past,
		EffectiveDate:      &past,
		ContentHash:        "hash2",
	}

	rel := &domain.Relationship{
		ProjectID:      "p-1",
		SourceEntityID: "p-1",
		TargetEntityID: "o-1",
		EvidenceID:     "ev-1",
	}

	rep := Evaluate(proj, []*domain.Evidence{ev1, ev2}, []*domain.Relationship{rel}, now)
	if rep.EvidenceCount != 2 {
		t.Fatalf("expected 2 evidence count, got %d", rep.EvidenceCount)
	}
	if rep.PrimarySourceCoverage != 50.0 {
		t.Fatalf("expected 50%% primary coverage, got %.1f", rep.PrimarySourceCoverage)
	}
	if rep.Staleness != "LOW" {
		t.Fatalf("expected LOW staleness, got %s", rep.Staleness)
	}
}

func TestMerkleTreeAndProof(t *testing.T) {
	now := time.Now().UTC()
	evidence := []*domain.Evidence{
		{ID: "e1", ContentHash: "hash-001", RetrievalTimestamp: now},
		{ID: "e2", ContentHash: "hash-002", RetrievalTimestamp: now},
		{ID: "e3", ContentHash: "hash-003", RetrievalTimestamp: now},
		{ID: "e4", ContentHash: "hash-004", RetrievalTimestamp: now},
	}

	tree := BuildMerkleTreeFromEvidence(evidence)
	if tree == nil || tree.RootHash == "" {
		t.Fatal("expected non-empty Merkle tree")
	}
	if tree.LeafCount != 4 {
		t.Fatalf("expected 4 leaves, got %d", tree.LeafCount)
	}

	proof, err := tree.GenerateProof("e3")
	if err != nil {
		t.Fatalf("unexpected proof error: %v", err)
	}
	if proof.Leaf.EvidenceID != "e3" {
		t.Fatalf("expected proof leaf e3, got %s", proof.Leaf.EvidenceID)
	}

	valid := VerifyProof(proof)
	if !valid {
		t.Fatal("expected proof verification to succeed")
	}

	pub := PublishMerkleRoot(tree)
	if pub == "" {
		t.Fatal("expected non-empty published root hash")
	}
}
