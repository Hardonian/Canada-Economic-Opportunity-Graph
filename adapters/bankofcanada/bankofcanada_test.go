package bankofcanada

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// minimalFixture mimics a BoC Valet JSON response shape used for parsing
// tests. The actual fixture file lives at data/fixtures/bank_of_canada_recent.json
// and is loaded by the application at runtime.
const minimalFixture = `{
  "observations": [
    {
      "d": "2026-09-01",
      "CPI_W": {"v": "139.8"},
      "GDPPV": {"v": "2400000.0"}
    },
    {
      "d": "2026-08-01",
      "CPI_W": {"v": "139.5"},
      "GDPPV": {"v": "2380000.0"}
    }
  ]
}`

func TestParseMinimalFixture(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "boc_min.json")
	if err := os.WriteFile(fp, []byte(minimalFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	a := NewBoCAdapter(fp)
	raw, err := a.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	res, err := a.Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if got, want := len(res.TradeMetrics), 4; got != want {
		t.Fatalf("TradeMetrics: got %d, want %d", got, want)
	}
	if got, want := len(res.Evidence), 1; got != want {
		t.Fatalf("Evidence: got %d, want %d", got, want)
	}
	if len(res.Evidence) > 0 && res.Evidence[0].ParserVersion != "boc-valet-json-v1" {
		t.Errorf("ParserVersion: got %q, want boc-valet-json-v1", res.Evidence[0].ParserVersion)
	}
}

func TestParseSkipsEmptyValues(t *testing.T) {
	// BoC Valet uses empty strings ("") for missing observations.
	// Adapter MUST NOT substitute 0 — it should skip them.
	empty := `{"observations": [{"d": "2026-09-01", "CPI_W": {"v": ""}, "GDPPV": {"v": "2400000.0"}}]}`
	dir := t.TempDir()
	fp := filepath.Join(dir, "boc_empty.json")
	if err := os.WriteFile(fp, []byte(empty), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewBoCAdapter(fp)
	raw, err := a.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	res, err := a.Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := len(res.TradeMetrics); got != 1 {
		t.Fatalf("TradeMetrics: got %d, want 1 (only non-empty)", got)
	}
	if res.TradeMetrics[0].MetricCode != "GDPPV" {
		t.Errorf("MetricCode: got %q, want GDPPV", res.TradeMetrics[0].MetricCode)
	}
}

func TestParseFailureOnInvalidJSON(t *testing.T) {
	a := NewBoCAdapter("")
	_, err := a.Parse([]byte("not json"))
	if err == nil {
		t.Fatal("expected error on invalid JSON")
	}
	if got := a.Health().ParseFailures; got != 1 {
		t.Errorf("ParseFailures: got %d, want 1", got)
	}
}

func TestLiveEnvDefaultsToFixture(t *testing.T) {
	t.Setenv("BANKOFCANADA_LIVE", "0")
	a := NewBoCAdapter("")
	h := a.Health()
	if h.Mode != "CURATED_SNAPSHOT" {
		t.Errorf("Mode with BANKOFCANADA_LIVE=0: got %q, want CURATED_SNAPSHOT", h.Mode)
	}
}

func TestNameAndTier(t *testing.T) {
	a := NewBoCAdapter("")
	if got := a.Name(); got != "bank_of_canada_valet" {
		t.Errorf("Name: got %q", got)
	}
	if got := a.Tier(); got != 1 {
		t.Errorf("Tier: got %v, want 1 (SourceTier1)", got)
	}
}

func TestFetchResponseShapeIsJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "boc_shape.json")
	if err := os.WriteFile(fp, []byte(minimalFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	a := NewBoCAdapter(fp)
	raw, _ := a.Fetch(context.Background())
	var v map[string]json.RawMessage
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Errorf("fixture must be valid JSON: %v", err)
	}
}

func TestInferUnit(t *testing.T) {
	cases := map[string]string{
		"EPS_CPI_INFLATION": "percent_yoy",
		"CPI_W":             "percent_yoy",
		"AVG.INTWO":         "percent",
		"CL.CDN.MOST.1DL":   "percent",
		"UNKNOWN":           "index",
	}
	for in, want := range cases {
		if got := inferUnit(in); got != want {
			t.Errorf("inferUnit(%q): got %q, want %q", in, got, want)
		}
	}
}