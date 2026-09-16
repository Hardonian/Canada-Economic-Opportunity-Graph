package lobbyist_registry

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

const lobFixture = `{
  "source": "Lobbyist Registry Test",
  "source_url": "https://lobbycanada.gc.ca/",
  "retrieved_at": "2026-09-15T00:00:00Z",
  "effective_at": "2026-09-14T00:00:00Z",
  "dataset_vintage": "2026-Q3",
  "registrations": [
    {
      "registration_id": "LOB-TEST-001",
      "registrant_name": "Test Industry Association",
      "registrant_type": "Industry Association",
      "client_organization": null,
      "subject_matters": ["Policy X", "Policy Y"],
      "subject_matters_fr": ["Politique X", "Politique Y"],
      "designated_public_office_holder": false,
      "active": true,
      "registration_date": "2026-04-15",
      "last_amended": "2026-08-30",
      "sector_focus": "Energy",
      "source_url": "https://example.com/LOB-TEST-001"
    },
    {
      "registration_id": "LOB-TEST-002",
      "registrant_name": "Inactive Test",
      "registrant_type": "Industry Association",
      "client_organization": null,
      "subject_matters": ["Z"],
      "subject_matters_fr": ["Z"],
      "designated_public_office_holder": false,
      "active": false,
      "registration_date": "2026-04-15",
      "last_amended": "2026-08-30",
      "sector_focus": "Energy",
      "source_url": "https://example.com/LOB-TEST-002"
    }
  ]
}`

func TestParseSkipsInactive(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "lob.json")
	if err := os.WriteFile(fp, []byte(lobFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewLobbyistRegistryAdapter(fp)
	data, err := a.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	res, err := a.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	// Only the active registration should be ingested.
	if got := len(res.Entities); got != 1 {
		t.Fatalf("Entities: got %d, want 1 (inactive must be skipped)", got)
	}
	if got := res.Entities[0].LegalName; got != "Test Industry Association" {
		t.Errorf("LegalName: got %q", got)
	}
	if got := res.Entities[0].Jurisdiction; got != "CA:FED" {
		t.Errorf("Jurisdiction: got %q", got)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	a := NewLobbyistRegistryAdapter("")
	if _, err := a.Parse([]byte("nope")); err == nil {
		t.Fatal("expected error")
	}
}

func TestNameAndTier(t *testing.T) {
	a := NewLobbyistRegistryAdapter("")
	if got := a.Name(); got != "lobbyist_registry" {
		t.Errorf("Name: got %q", got)
	}
	if got := a.Tier(); got != domain.SourceTier2 {
		t.Errorf("Tier: got %v", got)
	}
}

func TestFixtureIsValidJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "lob.json")
	if err := os.WriteFile(fp, []byte(lobFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewLobbyistRegistryAdapter(fp)
	data, _ := a.Fetch(context.Background())
	var v map[string]json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("fixture must parse as JSON object: %v", err)
	}
}