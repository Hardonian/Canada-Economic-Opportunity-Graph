package cmhc_housing

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

const cmhcFixture = `{
  "source": "CMHC Test",
  "source_url": "https://www.cmhc-schl.gc.ca/",
  "retrieved_at": "2026-09-15T00:00:00Z",
  "effective_at": "2026-09-14T00:00:00Z",
  "dataset_vintage": "2026-Q3",
  "indicators": [
    {
      "indicator_code": "STARTS-CAN",
      "indicator_name": "Total Housing Starts — Canada",
      "indicator_name_fr": "Mises en chantier totales — Canada",
      "geography": "CAN",
      "value": 244000,
      "unit": "units_saar",
      "reference_period": "2026-08",
      "source_url": "https://example.com/STARTS-CAN"
    },
    {
      "indicator_code": "VACANCY-CAN",
      "indicator_name": "Rental Vacancy Rate",
      "indicator_name_fr": "Taux d'inoccupation",
      "geography": "CAN",
      "value": 2.1,
      "unit": "percent",
      "reference_period": "2026-08",
      "source_url": "https://example.com/VACANCY-CAN"
    }
  ]
}`

func TestParseAllIndicators(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "cmhc.json")
	if err := os.WriteFile(fp, []byte(cmhcFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewCMHCHousingAdapter(fp)
	data, err := a.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	res, err := a.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := len(res.TradeMetrics); got != 2 {
		t.Fatalf("TradeMetrics: got %d, want 2", got)
	}
	if got := len(res.Evidence); got != 1 {
		t.Errorf("Evidence: got %d, want 1", got)
	}
	m := res.TradeMetrics[0]
	if m.Geography != "CAN" {
		t.Errorf("Geography: got %q, want CAN", m.Geography)
	}
	if m.MetricCode != "STARTS-CAN" {
		t.Errorf("MetricCode: got %q", m.MetricCode)
	}
	if m.Unit != "units_saar" {
		t.Errorf("Unit: got %q", m.Unit)
	}
	if m.ReferencePeriod != "2026-08" {
		t.Errorf("ReferencePeriod: got %q", m.ReferencePeriod)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	a := NewCMHCHousingAdapter("")
	if _, err := a.Parse([]byte("nope")); err == nil {
		t.Fatal("expected error")
	}
}

func TestNameAndTier(t *testing.T) {
	a := NewCMHCHousingAdapter("")
	if got := a.Name(); got != "cmhc_housing" {
		t.Errorf("Name: got %q", got)
	}
	if got := a.Tier(); got != domain.SourceTier1 {
		t.Errorf("Tier: got %v", got)
	}
}

func TestFixtureIsValidJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "cmhc.json")
	if err := os.WriteFile(fp, []byte(cmhcFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewCMHCHousingAdapter(fp)
	data, _ := a.Fetch(context.Background())
	var v map[string]json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("fixture must parse as JSON object: %v", err)
	}
}