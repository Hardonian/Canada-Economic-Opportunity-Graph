package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/adaptersandbox"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestAdapterSandboxMutationsRequireConfiguredMatchingSecret(t *testing.T) {
	registry := adaptersandbox.NewRegistry()
	if err := registry.Register(&adaptersandbox.AdapterEntry{
		Name:         "municipal-feed",
		Version:      "1.0.0",
		SourceURL:    "https://data.example.ca/municipal-feed",
		Tier:         domain.SourceTier1,
		RegisteredAt: time.Now().UTC(),
	}); err != nil {
		t.Fatal(err)
	}

	options := testOptions()
	options.AdapterRegistry = registry
	server := mustServer(t, database.NewMemoryStore(), options)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/adapters/municipal-feed/approve", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("unconfigured admin auth status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}

	options.AdapterAdminSecret = "correct-horse-battery-staple-admin-secret"
	server = mustServer(t, database.NewMemoryStore(), options)

	request = httptest.NewRequest(http.MethodPost, "/api/v1/adapters/municipal-feed/approve", nil)
	request.Header.Set("X-Admin-Secret", "any-non-empty-value")
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("wrong admin secret status = %d, want %d", response.Code, http.StatusUnauthorized)
	}

	request = httptest.NewRequest(http.MethodPost, "/api/v1/adapters/municipal-feed/approve", nil)
	request.Header.Set("X-Admin-Secret", options.AdapterAdminSecret)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("matching admin secret status = %d, body = %s", response.Code, response.Body.String())
	}
	if entry := registry.Get("municipal-feed"); entry == nil || !entry.Approved {
		t.Fatal("matching admin secret did not approve the adapter")
	}
}
