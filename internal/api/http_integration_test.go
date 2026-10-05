package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// TestPublicHTTPContract exercises mounted routes together rather than calling
// individual handlers, protecting the browser/API boundary.
func TestPublicHTTPContract(t *testing.T) {
	store := database.NewMemoryStore()
	if err := store.SaveProject(context.Background(), &domain.Project{
		ID: "integration-project", Slug: "integration-project", Name: "Integration Project",
		Sector: "Clean Energy & Grid", Province: "ON", CapexCAD: 1_000_000,
	}); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(mustServer(t, store, testOptions()))
	defer server.Close()

	for _, path := range []string{"/health", "/ready", "/api/v1/projects?limit=10", "/api/v1/radar", "/api/v1/system/status"} {
		response, err := http.Get(server.URL + path) // #nosec G107 -- httptest server URL.
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status = %d, want %d", path, response.StatusCode, http.StatusOK)
		}
	}

	statusResponse, err := http.Get(server.URL + "/api/v1/system/status") // #nosec G107 -- httptest server URL.
	if err != nil {
		t.Fatal(err)
	}
	defer statusResponse.Body.Close()
	var systemStatus struct {
		Status       string              `json:"status"`
		Capabilities []productCapability `json:"capabilities"`
		Data         database.RadarStats `json:"data"`
	}
	if err := json.NewDecoder(statusResponse.Body).Decode(&systemStatus); err != nil {
		t.Fatal(err)
	}
	if systemStatus.Status != "OPERATIONAL" || len(systemStatus.Capabilities) < 7 || systemStatus.Data.TotalProjects != 1 {
		t.Fatalf("system status = %#v", systemStatus)
	}

	response, err := http.Get(server.URL + "/api/v1/filings/recent") // #nosec G107 -- httptest server URL.
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("filings status = %d", response.StatusCode)
	}
	var payload struct {
		DataMode string `json:"data_mode"`
		Notice   string `json:"notice"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.DataMode != "DEMONSTRATION_SNAPSHOT" || payload.Notice == "" {
		t.Fatalf("filings source disclosure = %#v, want labelled demonstration snapshot", payload)
	}
}
