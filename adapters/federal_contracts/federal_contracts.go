// Package federal_contracts ingests Open.Canada Federal Contracts (>$10K) as
// Procurement evidence. Defaults to a reviewed fixture so the pipeline stays
// deterministic; setting FEDERAL_CONTRACTS_LIVE=1 fetches the Open.Canada
// contracts search (CSV) with strict allowlisting, response-size bounds, and a
// versioned parser. No API key required.
package federal_contracts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/identity"
)

const (
	adapterName      = "open_canada_federal_contracts"
	pipelineVersion  = "open-canada-contracts-v1"
	parserVersion    = "open-canada-contracts-json-v1"
	defaultFixture   = "data/fixtures/federal_contracts.json"
	maxResponseBytes = int64(4 << 20) // 4 MiB
	maxRedirects     = 2
	requestTimeout   = 10 * time.Second
)

// allowedHosts is the explicit allowlist for the live fetch.
var allowedHosts = map[string]bool{
	"search.open.canada.ca": true,
	"open.canada.ca":        true,
}

type FederalContractsAdapter struct {
	fixturePath string
	httpClient  *http.Client
	health      adapters.SourceHealth
}

func NewFederalContractsAdapter(fixturePath string) *FederalContractsAdapter {
	if fixturePath == "" {
		fixturePath = defaultFixture
	}
	return &FederalContractsAdapter{
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
			Tier:        domain.SourceTier1,
			Status:      "UNKNOWN",
			Mode:        "CURATED_SNAPSHOT",
		},
	}
}

func (a *FederalContractsAdapter) Name() string { return adapterName }
func (a *FederalContractsAdapter) Tier() domain.SourceTier {
	return domain.SourceTier1
}
func (a *FederalContractsAdapter) Health() *adapters.SourceHealth {
	return &a.health
}

func (a *FederalContractsAdapter) Fetch(ctx context.Context) ([]byte, error) {
	a.health.LastAttempt = time.Now().UTC()
	if os.Getenv("FEDERAL_CONTRACTS_LIVE") == "1" {
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
		return nil, fmt.Errorf("read federal contracts fixture: %w", err)
	}
	a.health.LastSuccess = time.Now().UTC()
	a.health.Status = "HEALTHY"
	a.health.Mode = "CURATED_SNAPSHOT"
	return data, nil
}

func (a *FederalContractsAdapter) fetchLive(ctx context.Context) ([]byte, error) {
	url := "https://search.open.canada.ca/contracts/?contract_value_from=10000&format=json"
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

type rawFederalContract struct {
	ContractID          string  `json:"contract_id"`
	VendorName          string  `json:"vendor_name"`
	VendorLocation      string  `json:"vendor_location"`
	BuyerDepartment     string  `json:"buyer_department"`
	BuyerDepartmentFr   string  `json:"buyer_department_fr"`
	Description         string  `json:"description"`
	DescriptionFr       string  `json:"description_fr"`
	ContractValueCad    int64   `json:"contract_value_cad"`
	AwardDate           string  `json:"award_date"`
	Sector              string  `json:"sector"`
	ProcurementCategory string  `json:"procurement_category"`
	SourceURL           string  `json:"source_url"`
}

type rawFederalContractsFile struct {
	Source         string               `json:"source"`
	SourceURL      string               `json:"source_url"`
	RetrievedAt    string               `json:"retrieved_at"`
	EffectiveAt    string               `json:"effective_at"`
	DatasetVintage string               `json:"dataset_vintage"`
	Contracts      []rawFederalContract `json:"contracts"`
}

func parseISODate(s string) *time.Time {
	if s == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
	}
	return nil
}

func (a *FederalContractsAdapter) Parse(data []byte) (*adapters.IngestionResult, error) {
	var f rawFederalContractsFile
	if err := json.Unmarshal(data, &f); err != nil {
		a.health.ParseFailures++
		return nil, fmt.Errorf("parse federal contracts JSON: %w", err)
	}

	sum := sha256.Sum256(data)
	contentHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()

	evID := identity.StableID("evidence", "open-canada-contracts", contentHash)
	evidence := &domain.Evidence{
		ID:                 evID,
		SourceURL:          f.SourceURL,
		Publisher:          "Open.Canada / Government of Canada",
		SourceTier:         domain.SourceTier1,
		Visibility:         domain.VisibilityPublicAttribution,
		Publishable:        true,
		RetrievalTimestamp: now,
		Confidence:         domain.ConfidenceReported,
		ExtractionMethod:   "open_canada_federal_contracts_adapter",
		ContentHash:        contentHash,
		HashScope:          "raw_contracts_payload",
		SourceClass:        "PROCUREMENT_RECORD",
		ParserVersion:      parserVersion,
	}

	res := &adapters.IngestionResult{Evidence: []*domain.Evidence{evidence}}
	a.health.DocumentsSeen = len(f.Contracts)
	a.health.DocumentsChanged = len(f.Contracts)

	for i, c := range f.Contracts {
		procID := identity.StableID("procurement", "open-canada-contracts", c.ContractID)
		proc := &domain.Procurement{
			ID:               procID,
			TenderID:         c.ContractID,
			Title:            c.Description,
			Stage:            "AWARDED",
			Buyer:            c.BuyerDepartment,
			BuyerType:        "Federal",
			EstimatedCAD:     c.ContractValueCad,
			ClosingDate:      parseISODate(c.AwardDate),
			SourceURL:        c.SourceURL,
			Categories:       []string{c.Sector, c.ProcurementCategory},
			RequirementClass: domain.RequirementConfirmed,
			EvidenceID:       evID,
			CreatedAt:        now,
		}
		_ = i
		res.Procurements = append(res.Procurements, proc)
	}
	return res, nil
}