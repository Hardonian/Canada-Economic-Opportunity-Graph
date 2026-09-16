package federal_contracts

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

const federalFixture = `{
  "source": "Open.Canada Federal Contracts",
  "source_url": "https://search.open.canada.ca/contracts/",
  "retrieved_at": "2026-09-15T00:00:00Z",
  "effective_at": "2026-09-14T00:00:00Z",
  "dataset_vintage": "2026-Q3",
  "contracts": [
    {
      "contract_id": "OC-2026-001",
      "vendor_name": "TestVendor",
      "buyer_department": "TestDept",
      "description": "Test contract",
      "contract_value_cad": 1000000,
      "award_date": "2026-08-01",
      "sector": "Defence",
      "procurement_category": "Services",
      "source_url": "https://example.com/OC-2026-001"
    }
  ]
}`

func TestFetchFixture(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "fc.json")
	if err := os.WriteFile(fp, []byte(federalFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewFederalContractsAdapter(fp)
	data, err := a.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	res, err := a.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := len(res.Procurements); got != 1 {
		t.Fatalf("Procurements: got %d, want 1", got)
	}
	if got := res.Procurements[0].Buyer; got != "TestDept" {
		t.Errorf("Buyer: got %q", got)
	}
	if got := res.Procurements[0].EstimatedCAD; got != 1000000 {
		t.Errorf("EstimatedCAD: got %d", got)
	}
	if got := res.Procurements[0].RequirementClass; got != domain.RequirementConfirmed {
		t.Errorf("RequirementClass: got %v, want %v", got, domain.RequirementConfirmed)
	}
	if len(res.Evidence) != 1 {
		t.Errorf("Evidence: got %d, want 1", len(res.Evidence))
	}
}

func TestParseInvalidJSON(t *testing.T) {
	a := NewFederalContractsAdapter("")
	if _, err := a.Parse([]byte("not json")); err == nil {
		t.Fatal("expected error")
	}
	if a.Health().ParseFailures != 1 {
		t.Errorf("ParseFailures: got %d, want 1", a.Health().ParseFailures)
	}
}

func TestNameAndTier(t *testing.T) {
	a := NewFederalContractsAdapter("")
	if got := a.Name(); got != "open_canada_federal_contracts" {
		t.Errorf("Name: got %q", got)
	}
	if got := a.Tier(); got != domain.SourceTier1 {
		t.Errorf("Tier: got %v", got)
	}
}

func TestParseISODate(t *testing.T) {
	if d := parseISODate(""); d != nil {
		t.Errorf("empty: got %v, want nil", d)
	}
	if d := parseISODate("not-a-date"); d != nil {
		t.Errorf("invalid: got %v, want nil", d)
	}
	if d := parseISODate("2026-08-01"); d == nil {
		t.Error("valid date: got nil, want non-nil")
	}
}

func TestFixturesAreValidJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "fc.json")
	if err := os.WriteFile(fp, []byte(federalFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewFederalContractsAdapter(fp)
	data, _ := a.Fetch(context.Background())
	var v map[string]json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("fixture must parse as JSON object: %v", err)
	}
}