package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsCollector_LifecycleAndRender(t *testing.T) {
	c := NewMetricsCollector()

	c.IncHTTP("GET", "/api/v1/projects", http.StatusOK, 12*time.Millisecond)
	c.IncHTTP("GET", "/api/v1/projects", http.StatusOK, 15*time.Millisecond)
	c.IncHTTP("POST", "/api/v1/graphql", http.StatusOK, 45*time.Millisecond)
	c.IncWALWrite()
	c.IncWALWrite()
	c.IncCheckpoint()
	c.IncAttestation(true)
	c.SetAuditorPeers(4)
	c.IncLakehouseQuery()

	output := c.RenderPrometheus()

	expectedSubstrings := []string{
		"cog_http_requests_total",
		"method=\"GET\",path=\"/api/v1/projects\",status=\"200\"} 2",
		"cog_storage_wal_writes_total 2",
		"cog_storage_checkpoints_total 1",
		"cog_verifier_attestations_total 1",
		"cog_verifier_quorums_reached_total 1",
		"cog_active_auditor_peers 4",
		"cog_lakehouse_queries_total 1",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("expected Prometheus output to contain %q, but was not found", sub)
		}
	}

	// Test HTTP handler endpoint
	handler := c.Handler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /metrics, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/plain") {
		t.Errorf("expected text/plain content type, got %s", rec.Header().Get("Content-Type"))
	}
}
