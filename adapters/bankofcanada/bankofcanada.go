// Package bankofcanada ingests official Bank of Canada Valet observations for
// CPI (national + provincial), GDP, and the policy rate without requiring any
// API key. Defaults to a reviewed fixture so the pipeline stays deterministic;
// setting BANKOFCANADA_LIVE=1 fetches the same series from the public Valet
// endpoint over HTTPS with strict allowlisting, response-size bounds, and a
// versioned parser.
package bankofcanada

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
	adapterName      = "bank_of_canada_valet"
	pipelineVersion  = "boc-valet-v1"
	parserVersion    = "boc-valet-json-v1"
	defaultFixture   = "data/fixtures/bank_of_canada_recent.json"
	maxResponseBytes = int64(2 << 20)
	maxRedirects     = 2
	requestTimeout   = 8 * time.Second
)

const valetAllowedHost = "www.bankofcanada.ca"

var defaultSeries = []string{
	"CPI_W",
	"CPI_QC",
	"CPI_ON",
	"GDPPV",
	"IRST_FF01",
}

type BoCAdapter struct {
	fixturePath string
	httpClient  *http.Client
	health      adapters.SourceHealth
}

func NewBoCAdapter(fixturePath string) *BoCAdapter {
	if fixturePath == "" {
		fixturePath = defaultFixture
	}
	return &BoCAdapter{
		fixturePath: fixturePath,
		httpClient: &http.Client{
			Timeout: requestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return fmt.Errorf("too many redirects (max %d)", maxRedirects)
				}
				if req.URL.Host != valetAllowedHost {
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

func (a *BoCAdapter) Name() string { return adapterName }
func (a *BoCAdapter) Tier() domain.SourceTier {
	return domain.SourceTier1
}
func (a *BoCAdapter) Health() *adapters.SourceHealth { return &a.health }

func (a *BoCAdapter) Fetch(ctx context.Context) ([]byte, error) {
	a.health.LastAttempt = time.Now().UTC()
	if os.Getenv("BANKOFCANADA_LIVE") == "1" {
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
		return nil, fmt.Errorf("read BoC fixture: %w", err)
	}
	a.health.LastSuccess = time.Now().UTC()
	a.health.Status = "HEALTHY"
	a.health.Mode = "CURATED_SNAPSHOT"
	return data, nil
}

func (a *BoCAdapter) fetchLive(ctx context.Context) ([]byte, error) {
	url := "https://" + valetAllowedHost + "/valet/observations/JSON/" +
		strings.Join(defaultSeries, ",") + "/recent"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ceo-g-canada-economic-opportunity-graph/"+pipelineVersion)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("live fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("live fetch HTTP %d", resp.StatusCode)
	}
	limited := io.LimitReader(resp.Body, maxResponseBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read live body: %w", err)
	}
	if int64(len(body)) > maxResponseBytes {
		return nil, fmt.Errorf("live response exceeds %d byte limit", maxResponseBytes)
	}
	return body, nil
}

// rawValetResponse mirrors the BoC Valet JSON shape (subset).
type rawValetResponse struct {
	Observations []map[string]json.RawMessage `json:"observations"`
}

func (a *BoCAdapter) Parse(data []byte) (*adapters.IngestionResult, error) {
	var resp rawValetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		a.health.ParseFailures++
		return nil, fmt.Errorf("parse BoC Valet JSON: %w", err)
	}

	sum := sha256.Sum256(data)
	contentHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()

	evID := identity.StableID("evidence", "boc-valet", contentHash)
	evidence := &domain.Evidence{
		ID:                 evID,
		SourceURL:          "https://www.bankofcanada.ca/valet/observations",
		Publisher:          "Bank of Canada / Banque du Canada",
		SourceTier:         domain.SourceTier1,
		Visibility:         domain.VisibilityPublicAttribution,
		Publishable:        true,
		RetrievalTimestamp: now,
		Confidence:         domain.ConfidenceReported,
		ExtractionMethod:   "bank_of_canada_valet_adapter",
		ContentHash:        contentHash,
		HashScope:          "raw_observation_payload",
		SourceClass:        "OFFICIAL_STATISTICS",
		ParserVersion:      parserVersion,
	}

	res := &adapters.IngestionResult{
		Evidence: []*domain.Evidence{evidence},
	}
	a.health.DocumentsSeen = len(resp.Observations)
	a.health.DocumentsChanged = len(resp.Observations)

	for _, obs := range resp.Observations {
		dateRaw, ok := obs["d"]
		if !ok {
			continue
		}
		var date string
		_ = json.Unmarshal(dateRaw, &date)
		var periodTime time.Time
		if t, err := time.Parse("2006-01-02", date); err == nil {
			periodTime = t
		} else {
			periodTime = now
		}
		for series, raw := range obs {
			if series == "d" {
				continue
			}
			var point struct {
				V string `json:"v"`
			}
			if err := json.Unmarshal(raw, &point); err != nil {
				continue
			}
			if point.V == "" {
				continue
			}
			value := parseFloat(point.V)
			metricID := identity.StableID(
				"tradecap",
				"boc-valet",
				series+"-"+date,
			)
			res.TradeMetrics = append(res.TradeMetrics, &domain.TradeMetric{
				ID:              metricID,
				Geography:       "CA",
				MetricCode:      series,
				MetricName:      series + " (BoC Valet)",
				ReferencePeriod: date,
				Value:           value,
				Unit:            inferUnit(series),
				EvidenceID:      evID,
				ObservedAt:      periodTime,
			})
		}
	}
	return res, nil
}

func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}

func inferUnit(series string) string {
	switch {
	case strings.HasPrefix(series, "CPI"):
		return "index_2002=100"
	case series == "GDPPV":
		return "chained_2017_CAD_millions"
	case strings.HasPrefix(series, "IRST"):
		return "percent"
	default:
		return "index"
	}
}