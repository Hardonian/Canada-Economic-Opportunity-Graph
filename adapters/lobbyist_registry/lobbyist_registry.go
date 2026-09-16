// Package lobbyist_registry ingests active federal lobbyist registrations from
// the Office of the Commissioner of Lobbying of Canada. Defaults to a reviewed
// fixture; setting LOBBYIST_REGISTRY_LIVE=1 fetches the registry endpoint with
// strict allowlisting, response-size bounds, and a versioned parser. No API
// key required (public registry).
package lobbyist_registry

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/identity"
)

const (
	adapterName      = "lobbyist_registry"
	pipelineVersion  = "lobbyist-registry-v1"
	parserVersion    = "lobbyist-registry-json-v1"
	defaultFixture   = "data/fixtures/lobbyist_registry.json"
	maxResponseBytes = int64(4 << 20)
	maxRedirects     = 2
	requestTimeout   = 10 * time.Second
)

var allowedHosts = map[string]bool{
	"lobbycanada.gc.ca": true,
}

type LobbyistRegistryAdapter struct {
	fixturePath string
	httpClient  *http.Client
	health      adapters.SourceHealth
}

func NewLobbyistRegistryAdapter(fixturePath string) *LobbyistRegistryAdapter {
	if fixturePath == "" {
		fixturePath = defaultFixture
	}
	return &LobbyistRegistryAdapter{
		fixturePath: fixturePath,
		httpClient: &http.Client{
			Timeout: requestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return fmt.Errorf("too many redirects (max %d)", maxRedirects)
				}
				if !allowedHosts[req.URL.Host] {
					return fmt.Errorf("redirect to disallowed host %q", req.URL.Host)
				}
				return nil
			},
		},
		health: adapters.SourceHealth{
			AdapterName: adapterName,
			Tier:        domain.SourceTier2,
			Status:      "UNKNOWN",
			Mode:        "CURATED_SNAPSHOT",
		},
	}
}

func (a *LobbyistRegistryAdapter) Name() string { return adapterName }
func (a *LobbyistRegistryAdapter) Tier() domain.SourceTier { return domain.SourceTier2 }
func (a *LobbyistRegistryAdapter) Health() *adapters.SourceHealth { return &a.health }

func (a *LobbyistRegistryAdapter) Fetch(ctx context.Context) ([]byte, error) {
	a.health.LastAttempt = time.Now().UTC()
	if os.Getenv("LOBBYIST_REGISTRY_LIVE") == "1" {
		data, err := a.fetchLive(ctx)
		if err == nil {
			a.health.LastSuccess = time.Now().UTC()
			a.health.Status = "HEALTHY"
			a.health.Mode = "LIVE"
			return data, nil
		}
		a.health.Status = "DEGRADED"
		a.health.LastError = err.Error()
	}
	data, err := adapters.ReadBoundedFile(a.fixturePath, adapters.MaxFixtureBytes)
	if err != nil {
		a.health.Status = "DEGRADED"
		a.health.LastError = err.Error()
		return nil, fmt.Errorf("read lobbyist fixture: %w", err)
	}
	a.health.LastSuccess = time.Now().UTC()
	a.health.Status = "HEALTHY"
	a.health.Mode = "CURATED_SNAPSHOT"
	return data, nil
}

func (a *LobbyistRegistryAdapter) fetchLive(ctx context.Context) ([]byte, error) {
	url := "https://lobbycanada.gc.ca/app/secure/ocl/lrs/do/vwByLrdsrch?lang=E"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ceo-g-canada-economic-opportunity-graph/"+pipelineVersion)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	limited := io.LimitReader(resp.Body, maxResponseBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxResponseBytes {
		return nil, fmt.Errorf("response exceeds %d byte limit", maxResponseBytes)
	}
	return body, nil
}

type rawRegistration struct {
	RegistrationID              string   `json:"registration_id"`
	RegistrantName              string   `json:"registrant_name"`
	RegistrantType              string   `json:"registrant_type"`
	ClientOrganization          *string  `json:"client_organization"`
	SubjectMatters              []string `json:"subject_matters"`
	SubjectMattersFr            []string `json:"subject_matters_fr"`
	DesignatedPublicOfficeHolder bool    `json:"designated_public_office_holder"`
	Active                      bool     `json:"active"`
	RegistrationDate            string   `json:"registration_date"`
	LastAmended                 string   `json:"last_amended"`
	SectorFocus                 string   `json:"sector_focus"`
	SourceURL                   string   `json:"source_url"`
}

type rawLobbyistFile struct {
	Source        string             `json:"source"`
	SourceURL     string             `json:"source_url"`
	RetrievedAt   string             `json:"retrieved_at"`
	EffectiveAt   string             `json:"effective_at"`
	DatasetVintage string            `json:"dataset_vintage"`
	Registrations []rawRegistration  `json:"registrations"`
}

func (a *LobbyistRegistryAdapter) Parse(data []byte) (*adapters.IngestionResult, error) {
	var f rawLobbyistFile
	if err := json.Unmarshal(data, &f); err != nil {
		a.health.ParseFailures++
		return nil, fmt.Errorf("parse lobbyist JSON: %w", err)
	}

	sum := sha256.Sum256(data)
	contentHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()

	evID := identity.StableID("evidence", "lobbyist-registry", contentHash)
	evidence := &domain.Evidence{
		ID:                 evID,
		SourceURL:          f.SourceURL,
		Publisher:          "Office of the Commissioner of Lobbying of Canada / Commissaire au lobbying du Canada",
		SourceTier:         domain.SourceTier2,
		Visibility:         domain.VisibilityPublicAttribution,
		Publishable:        true,
		RetrievalTimestamp: now,
		Confidence:         domain.ConfidenceReported,
		ExtractionMethod:   "lobbyist_registry_adapter",
		ContentHash:        contentHash,
		HashScope:          "raw_registrations_payload",
		SourceClass:        "POLITICAL_DISCLOSURE",
		ParserVersion:      parserVersion,
	}

	res := &adapters.IngestionResult{Evidence: []*domain.Evidence{evidence}}
	a.health.DocumentsSeen = len(f.Registrations)
	a.health.DocumentsChanged = len(f.Registrations)

	for _, r := range f.Registrations {
		if !r.Active {
			continue
		}
		entityID := identity.StableID("entity", "lobbyist", r.RegistrantName)
		slug := strings.ToLower(strings.ReplaceAll(r.RegistrantName, " ", "-"))
		res.Entities = append(res.Entities, &domain.Entity{
			ID:           entityID,
			Slug:         slug,
			LegalName:    r.RegistrantName,
			CommonName:   r.RegistrantName,
			Aliases:      []string{r.RegistrantName},
			EntityType:   r.RegistrantType,
			Jurisdiction: "CA:FED",
			EvidenceID:   evID,
			Evidence:     evidence,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	return res, nil
}