package official

import (
	"context"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestOfficialAdapter_Lifecycle(t *testing.T) {
	adapter := NewAdapter("nonexistent.json")
	if adapter.Name() != adapterName {
		t.Errorf("expected adapter name %s, got %s", adapterName, adapter.Name())
	}
	if adapter.Tier() != domain.SourceTier1 {
		t.Errorf("expected tier 1, got %v", adapter.Tier())
	}

	h := adapter.Health()
	if h.Mode != "CURATED_SNAPSHOT" {
		t.Errorf("expected mode CURATED_SNAPSHOT, got %s", h.Mode)
	}

	_, err := adapter.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error reading nonexistent fixture")
	}

	// Test valid JSON parse
	sampleJSON := []byte(`[
		{
			"external_id": "OFF-001",
			"name": "Darlington SMR Nuclear",
			"summary": "Small modular reactor deployment in Clarington.",
			"sector": "CLEAN_ENERGY",
			"subsector": "Nuclear SMR",
			"province": "ON",
			"location_name": "Clarington",
			"stage": "APPROVED",
			"capex_cad": 1000000000,
			"capex_status": "HIGH",
			"confidence": "HIGH",
			"proponent": {
				"legal_name": "Ontario Power Generation Inc.",
				"common_name": "OPG",
				"entity_type": "Crown Corporation",
				"jurisdiction": "Ontario"
			},
			"sources": [
				{
					"source_id": "SRC-001",
					"publisher": "Ontario Power Generation",
					"source_url": "https://www.opg.com/smr",
					"publication_date": "2024-01-01T00:00:00Z",
					"effective_date": "2024-01-01T00:00:00Z",
					"retrieved_at": "2024-01-01T00:00:00Z",
					"source_class": "OFFICIAL_RECORD",
					"source_tier": 1,
					"locator": "Press Release",
					"excerpt": "OPG announces deployment of small modular reactor at Darlington.",
					"confidence": "VERIFIED"
				}
			]
		}
	]`)

	result, err := adapter.Parse(sampleJSON)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(result.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(result.Projects))
	}
	if result.Projects[0].Name != "Darlington SMR Nuclear" {
		t.Errorf("expected project name Darlington SMR Nuclear, got %s", result.Projects[0].Name)
	}
	if len(result.Entities) != 1 {
		t.Errorf("expected 1 proponent, got %d", len(result.Entities))
	}
	if len(result.Evidence) != 1 {
		t.Errorf("expected 1 evidence, got %d", len(result.Evidence))
	}
}
