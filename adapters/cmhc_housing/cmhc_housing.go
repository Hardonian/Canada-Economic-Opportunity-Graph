// Package cmhc_housing ingests housing market indicators from the Canada
// Mortgage and Housing Corporation (CMHC). Defaults to a reviewed fixture;
// setting CMHC_HOUSING_LIVE=1 fetches the public CMHC housing data tables
// with strict allowlisting, response-size bounds, and a versioned parser.
// No API key required.
package cmhc_housing

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
	adapterName      = "cmhc_housing"
	pipelineVersion  = "cmhc-housing-v1"
	parserVersion    = "cmhc-housing-json-v1"
	defaultFixture   = "data/fixtures/cmhc_housing.json"
	maxResponseBytes = int64(4 << 20)
	maxRedirects     = 2
	requestTimeout   = 10 * time.Second
)

var allowedHosts = map[string]bool{
	"www.cmhc-schl.gc.ca": true,
	"cmhc-schl.gc.ca":     true,
}

type CMHCHousingAdapter struct {
	fixturePath string
	httpClient  *http.Client
	health      adapters.SourceHealth
}

func NewCMHCHousingAdapter(fixturePath string) *CMHCHousingAdapter {
	if fixturePath == "" {
		fixturePath = defaultFixture
	}
	return &CMHCHousingAdapter{
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

func (a *CMHCHousingAdapter) Name() string { return adapterName }
func (a *CMHCHousingAdapter) Tier() domain.SourceTier { return domain.SourceTier1 }
func (a *CMHCHousingAdapter) Health() *adapters.SourceHealth { return &a.health }

func (a *CMHCHousingAdapter) Fetch(ctx context.Context) ([]byte, error) {
	a.health.LastAttempt = time.Now().UTC()
	if os.Getenv("CMHC_HOUSING_LIVE") == "1" {
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
		return nil, fmt.Errorf("read CMHC fixture: %w", err)
	}
	a.health.LastSuccess = time.Now().UTC()
	a.health.Status = "HEALTHY"
	a.health.Mode = "CURATED_SNAPSHOT"
	return data, nil
}

func (a *CMHCHousingAdapter) fetchLive(ctx context.Context) ([]byte, error) {
	url := "https://www.cmhc-schl.gc.ca/professionals/housing-markets-data-and-research/housing-data/tables"
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

type rawIndicator struct {
	IndicatorCode     string  `json:"indicator_code"`
	IndicatorName     string  `json:"indicator_name"`
	IndicatorNameFr   string  `json:"indicator_name_fr"`
	Geography         string  `json:"geography"`
	Value             float64 `json:"value"`
	Unit              string  `json:"unit"`
	ReferencePeriod   string  `json:"reference_period"`
	SourceURL         string  `json:"source_url"`
}

type rawCMHCFile struct {
	Source         string         `json:"source"`
	SourceURL      string         `json:"source_url"`
	RetrievedAt    string         `json:"retrieved_at"`
	EffectiveAt    string         `json:"effective_at"`
	DatasetVintage string         `json:"dataset_vintage"`
	Indicators     []rawIndicator `json:"indicators"`
}

func (a *CMHCHousingAdapter) Parse(data []byte) (*adapters.IngestionResult, error) {
	var f rawCMHCFile
	if err := json.Unmarshal(data, &f); err != nil {
		a.health.ParseFailures++
		return nil, fmt.Errorf("parse CMHC JSON: %w", err)
	}

	sum := sha256.Sum256(data)
	contentHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()

	evID := identity.StableID("evidence", "cmhc-housing", contentHash)
	evidence := &domain.Evidence{
		ID:                 evID,
		SourceURL:          f.SourceURL,
		Publisher:          "Canada Mortgage and Housing Corporation / Société canadienne d'hypothèques et de logement",
		SourceTier:         domain.SourceTier1,
		Visibility:         domain.VisibilityPublicAttribution,
		Publishable:        true,
		RetrievalTimestamp: now,
		Confidence:         domain.ConfidenceReported,
		ExtractionMethod:   "cmhc_housing_adapter",
		ContentHash:        contentHash,
		HashScope:          "raw_indicators_payload",
		SourceClass:        "OFFICIAL_STATISTICS",
		ParserVersion:      parserVersion,
	}

	res := &adapters.IngestionResult{Evidence: []*domain.Evidence{evidence}}
	a.health.DocumentsSeen = len(f.Indicators)
	a.health.DocumentsChanged = len(f.Indicators)

	for _, ind := range f.Indicators {
		metricID := identity.StableID("tradecap", "cmhc-housing", ind.IndicatorCode+"-"+ind.Geography+"-"+ind.ReferencePeriod)
		geo := strings.ToUpper(ind.Geography)
		period := ind.ReferencePeriod + "-01T00:00:00Z"
		var observedAt time.Time
		if t, err := time.Parse("2006-01-02T15:04:05Z", period); err == nil {
			observedAt = t
		} else {
			observedAt = now
		}
		res.TradeMetrics = append(res.TradeMetrics, &domain.TradeMetric{
			ID:              metricID,
			Geography:       geo,
			MetricCode:      ind.IndicatorCode,
			MetricName:      ind.IndicatorName,
			ReferencePeriod: ind.ReferencePeriod,
			Value:           ind.Value,
			Unit:            ind.Unit,
			EvidenceID:      evID,
			ObservedAt:      observedAt,
		})
	}
	return res, nil
}