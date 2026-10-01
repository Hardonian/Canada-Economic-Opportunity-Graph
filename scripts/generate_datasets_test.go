package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestValidateHistoricalArtifacts(t *testing.T) {
	releaseRoot := t.TempDir()
	files := map[string][]byte{
		"public/manifest.json": []byte(`{"dataset_version":"test"}\n`),
		"cegs/projects.jsonl":  []byte(`{"id":"project-1"}\n`),
	}

	if err := validateHistoricalArtifacts(releaseRoot, files); err != nil {
		t.Fatalf("missing release should be valid: %v", err)
	}
	for name, data := range files {
		path := filepath.Join(releaseRoot, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := validateHistoricalArtifacts(releaseRoot, files); err != nil {
		t.Fatalf("matching release should be valid: %v", err)
	}

	files["public/manifest.json"] = []byte(`{"dataset_version":"changed"}\n`)
	err := validateHistoricalArtifacts(releaseRoot, files)
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite historical release") {
		t.Fatalf("mismatched release error = %v", err)
	}
}

func TestIsSnapshotGeneratedEvidence(t *testing.T) {
	generated := &domain.Evidence{Publisher: "Bank of Canada / Banque du Canada"}
	if !isSnapshotGeneratedEvidence(generated) {
		t.Fatal("Bank of Canada fixture evidence should use the collection timestamp")
	}
	primary := &domain.Evidence{Publisher: "Natural Resources Canada", RetrievalTimestamp: time.Now()}
	if isSnapshotGeneratedEvidence(primary) {
		t.Fatal("primary-source fixture evidence must retain its source timestamp")
	}
	if isSnapshotGeneratedEvidence(nil) {
		t.Fatal("nil evidence cannot be snapshot-generated")
	}
}
