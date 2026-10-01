package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/adaptersandbox"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/verifier"
)

func TestABACAndVerifierWiring(t *testing.T) {
	store := database.NewMemoryStore()
	testProj := &domain.Project{
		ID:        "proj-defence-test",
		Slug:      "arctic-radar-base",
		Name:      "Arctic Deep Defense Array",
		Sector:    "Clean Energy & Grid",
		Subsector: "Military Radar Support",
		Summary:   "Classified radar outpost near NORAD_SITE_TEST99.",
		Province:  "NU",
		CapexCAD:  450000000,
	}
	_ = store.SaveProject(context.Background(), testProj)

	verifierNet := verifier.NewNetwork(5)
	attStore := verifier.NewAttestationStore()
	sandboxReg := adaptersandbox.NewRegistry()

	server, err := NewServerWithOptions(store, Options{
		AllowedOrigins:      []string{"*"},
		MaxRequestBodyBytes: 4096,
		RateLimitPerMinute:  600,
		RateLimitBurst:      100,
		RequestTimeout:      1000000000,
		ReadinessTimeout:    1000000000,
		AdapterRegistry:     sandboxReg,
		VerifierStore:       attStore,
		VerifierNetwork:     verifierNet,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// 1. Test GET /api/v1/projects/{id} with unclassified user
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/proj-defence-test", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("X-Security-Classification") != "UNCLASSIFIED" {
		t.Errorf("expected UNCLASSIFIED classification header, got %s", rec.Header().Get("X-Security-Classification"))
	}

	var resp struct {
		Project *domain.Project `json:"project"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Project.Summary != "[SUMMARY RESTRICTED - REQUIRES CANADIAN SECRET CLEARANCE]" {
		t.Errorf("expected redacted summary for unclassified user, got %q", resp.Project.Summary)
	}

	// 2. Test GET /api/v1/projects/{id} with SECRET user
	reqSecret := httptest.NewRequest(http.MethodGet, "/api/v1/projects/proj-defence-test", nil)
	reqSecret.Header.Set("X-Clearance-Level", "SECRET")
	recSecret := httptest.NewRecorder()
	server.ServeHTTP(recSecret, reqSecret)

	if recSecret.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for secret user, got %d", recSecret.Code)
	}
	if recSecret.Header().Get("X-Security-Classification") != "SECRET" {
		t.Errorf("expected SECRET classification header, got %s", recSecret.Header().Get("X-Security-Classification"))
	}
	var respSecret struct {
		Project *domain.Project `json:"project"`
	}
	_ = json.Unmarshal(recSecret.Body.Bytes(), &respSecret)
	if respSecret.Project.Summary == "[SUMMARY RESTRICTED - REQUIRES CANADIAN SECRET CLEARANCE]" {
		t.Errorf("expected unredacted summary for SECRET user")
	}

	// 3. Test GET /api/v1/export/project/{id}?format=markdown with DLP scrub
	reqExport := httptest.NewRequest(http.MethodGet, "/api/v1/export/project/proj-defence-test?format=markdown", nil)
	recExport := httptest.NewRecorder()
	server.ServeHTTP(recExport, reqExport)

	if recExport.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recExport.Code)
	}
	if recExport.Header().Get("X-Security-Classification") != "UNCLASSIFIED" {
		t.Errorf("expected UNCLASSIFIED export header")
	}

	// 4. Test Verifier Attestations endpoint (was returning 503 prior to wiring)
	reqVerifier := httptest.NewRequest(http.MethodGet, "/api/v1/verifier/attestations", nil)
	recVerifier := httptest.NewRecorder()
	server.ServeHTTP(recVerifier, reqVerifier)

	if recVerifier.Code != http.StatusOK {
		t.Errorf("expected 200 OK from verifier attestations endpoint, got %d", recVerifier.Code)
	}

	// 5. Test Sandbox Adapters endpoint (was returning 503 prior to wiring)
	reqAdapters := httptest.NewRequest(http.MethodGet, "/api/v1/adapters", nil)
	recAdapters := httptest.NewRecorder()
	server.ServeHTTP(recAdapters, reqAdapters)

	if recAdapters.Code != http.StatusOK {
		t.Errorf("expected 200 OK from adapters endpoint, got %d", recAdapters.Code)
	}
}
