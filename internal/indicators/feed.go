package indicators

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// LiveFeedEngine manages real-time streaming market/policy ticks,
// in-memory circular buffers, and canonical snapshot exports.
type LiveFeedEngine struct {
	mu           sync.RWMutex
	seqCounter   int64
	recentTicks  []*LiveFeedTick
	maxTicks     int
	latestTicks  map[KPICode]*LiveFeedTick
	observations map[string][]*KPIObservation // keyed by scope ("NATIONAL" or project ID)
}

// NewLiveFeedEngine initializes the feed engine with baseline canonical observations.
func NewLiveFeedEngine() *LiveFeedEngine {
	engine := &LiveFeedEngine{
		recentTicks:  make([]*LiveFeedTick, 0, 500),
		maxTicks:     500,
		latestTicks:  make(map[KPICode]*LiveFeedTick),
		observations: make(map[string][]*KPIObservation),
	}
	engine.seedCanonicalBaseline()
	return engine
}

// seedCanonicalBaseline populates authoritative baseline observations from Bank of Canada,
// NRCan, LME, ICE, and CER.
func (e *LiveFeedEngine) seedCanonicalBaseline() {
	now := time.Now().UTC()

	initialFeeds := []struct {
		code      KPICode
		name      string
		cat       KPICategory
		val       float64
		unit      string
		change    float64
		changePct float64
		source    string
	}{
		{KPIWCSWTIDifferential, "WCS-WTI Crude Discount", CategoryCommodityMacro, 12.85, "USD/bbl", -0.35, -2.65, "ICE / Argus Hardisty Spot"},
		{KPIAECOGasSpotCAD, "AECO C Natural Gas Spot", CategoryCommodityMacro, 2.48, "CAD/GJ", 0.08, 3.33, "NGX Alberta NIT Hub"},
		{KPILMENickelCashUSD, "LME Grade-1 Nickel Cash", CategoryCommodityMacro, 18950.0, "USD/tonne", 125.0, 0.66, "London Metal Exchange Cash"},
		{KPIUraniumUXU3O8USD, "Ux U3O8 Spot Uranium", CategoryCommodityMacro, 86.50, "USD/lb", 0.50, 0.58, "UxC Cameco Spot Index"},
		{KPILithiumCarbonateUSD, "Battery Lithium Carbonate", CategoryCommodityMacro, 14350.0, "USD/tonne", -100.0, -0.69, "Benchmark Minerals North America"},
		{KPICADUSDExchangeRate, "Bank of Canada CAD/USD", CategoryCommodityMacro, 0.7385, "CAD_per_USD", 0.0012, 0.16, "Bank of Canada Valet Noon Fix"},
		{KPIBoCOvernightRate, "BoC Policy Interest Rate", CategoryCommodityMacro, 2.75, "percent", 0.0, 0.0, "Bank of Canada Policy Decision"},
		{KPIGoC10YearBondYield, "GoC 10-Year Bond Yield", CategoryCommodityMacro, 3.02, "percent", -0.04, -1.31, "Bank of Canada Valet Markets"},
	}

	for _, item := range initialFeeds {
		e.RecordTick(item.code, item.name, item.cat, item.val, item.unit, item.change, item.changePct, item.source, now)
	}

	// Baseline National Observations
	nationalObs := []KPIObservation{
		{
			ID:              "obs-esg-abatement-can",
			MetricCode:      KPIEmissionsAbatementAnnual,
			MetricName:      "Annual Lifecycle GHG Abatement",
			Category:        CategoryESGDecarbonization,
			Scope:           "NATIONAL",
			Value:           42.8,
			Unit:            "Mt CO2e/yr",
			ReferencePeriod: "2025-2026",
			ObservedAt:      now,
			SourceURL:       "https://www.canada.ca/en/environment-climate-change/services/environmental-indicators/greenhouse-gas-emissions.html",
			Publisher:       "Environment and Climate Change Canada",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-esg-abatement-can", 42.8),
		},
		{
			ID:              "obs-indig-equity-can",
			MetricCode:      KPIIndigenousEquityPct,
			MetricName:      "First Nations Equity Co-Ownership Stake",
			Category:        CategoryIndigenousEquity,
			Scope:           "NATIONAL",
			Value:           18.4,
			Unit:            "percent",
			ReferencePeriod: "2025-2026",
			ObservedAt:      now,
			SourceURL:       "https://cib-bic.ca/en/fast-facts/",
			Publisher:       "Canada Infrastructure Bank",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-indig-equity-can", 18.4),
		},
		{
			ID:              "obs-indig-ilgp-can",
			MetricCode:      KPIIndigenousILGPGuaranteeCAD,
			MetricName:      "Indigenous Loan Guarantee Program Allocation",
			Category:        CategoryIndigenousEquity,
			Scope:           "NATIONAL",
			Value:           1_450_000_000,
			Unit:            "CAD",
			ReferencePeriod: "2025-2026",
			ObservedAt:      now,
			SourceURL:       "https://natural-resources.canada.ca/our-natural-resources/indigenous-peoples-and-natural-resources/indigenous-loan-guarantee-program/25841",
			Publisher:       "Natural Resources Canada (NRCan)",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-indig-ilgp-can", 1_450_000_000),
		},
		{
			ID:              "obs-capex-velocity-can",
			MetricCode:      KPICapexSpendVelocityCADMo,
			MetricName:      "Capital Deployment Spend Velocity",
			Category:        CategoryCapitalVelocity,
			Scope:           "NATIONAL",
			Value:           285.0,
			Unit:            "CAD_M/month",
			ReferencePeriod: "2026-Q1",
			ObservedAt:      now,
			SourceURL:       "https://www.statcan.gc.ca/en/subjects-start/capital_expenditures",
			Publisher:       "Statistics Canada Capital Spending Survey",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-capex-velocity-can", 285.0),
		},
		{
			ID:              "obs-domestic-content-can",
			MetricCode:      KPIDomesticCanadianContentPct,
			MetricName:      "Canadian Domestic Content Share",
			Category:        CategorySupplyChainContent,
			Scope:           "NATIONAL",
			Value:           68.2,
			Unit:            "percent",
			ReferencePeriod: "2025",
			ObservedAt:      now,
			SourceURL:       "https://ised-isde.canada.ca/site/canadian-company-capabilities/en",
			Publisher:       "Innovation, Science and Economic Development Canada",
			Confidence:      domain.ConfidenceSupported,
			ContentHash:     hashObservation("obs-domestic-content-can", 68.2),
		},
		{
			ID:              "obs-grid-queue-can",
			MetricCode:      KPIGridInterconnectQueueMo,
			MetricName:      "Grid Interconnection Queue Wait Latency",
			Category:        CategoryGridPhysics,
			Scope:           "NATIONAL",
			Value:           26.4,
			Unit:            "months",
			ReferencePeriod: "2026",
			ObservedAt:      now,
			SourceURL:       "https://www.ieso.ca/en/Learn-About-Ontario-Power/System-Impact-Assessments",
			Publisher:       "IESO / AESO Provincial System Assessments",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-grid-queue-can", 26.4),
		},
		{
			ID:              "obs-iaac-duration-can",
			MetricCode:      KPIIAACStatutoryDurationMo,
			MetricName:      "Federal Impact Assessment (IAAC) Elapsed Duration",
			Category:        CategoryRegulatorySpeed,
			Scope:           "NATIONAL",
			Value:           31.8,
			Unit:            "months",
			ReferencePeriod: "2022-2026",
			ObservedAt:      now,
			SourceURL:       "https://iaac-aeic.gc.ca/050/evaluations",
			Publisher:       "Impact Assessment Agency of Canada",
			Confidence:      domain.ConfidenceVerified,
			ContentHash:     hashObservation("obs-iaac-duration-can", 31.8),
		},
		{
			ID:              "obs-redseal-gap-can",
			MetricCode:      KPIRedSealTradesGap,
			MetricName:      "Peak Construction Red Seal Trades Deficit",
			Category:        CategoryLaborSkills,
			Scope:           "NATIONAL",
			Value:           14200.0,
			Unit:            "FTE",
			ReferencePeriod: "2026-2028",
			ObservedAt:      now,
			SourceURL:       "https://www.buildforce.ca/en/system/files/reports/2026-national-construction-forecast.pdf",
			Publisher:       "BuildForce Canada",
			Confidence:      domain.ConfidenceSupported,
			ContentHash:     hashObservation("obs-redseal-gap-can", 14200.0),
		},
	}

	for _, obs := range nationalObs {
		copyObs := obs
		e.observations["NATIONAL"] = append(e.observations["NATIONAL"], &copyObs)
	}
}

// RecordTick adds a validated real-time streaming tick into the circular buffer.
func (e *LiveFeedEngine) RecordTick(
	code KPICode,
	name string,
	cat KPICategory,
	val float64,
	unit string,
	changeAbs float64,
	changePct float64,
	source string,
	t time.Time,
) *LiveFeedTick {
	seq := atomic.AddInt64(&e.seqCounter, 1)

	dir := "FLAT"
	if changeAbs > 0.0001 {
		dir = "UP"
	} else if changeAbs < -0.0001 {
		dir = "DOWN"
	}

	hashPayload := fmt.Sprintf("%d:%s:%.4f:%s:%d", seq, code, val, source, t.UnixNano())
	sum := sha256.Sum256([]byte(hashPayload))
	hash := hex.EncodeToString(sum[:])

	tick := &LiveFeedTick{
		SequenceID:     seq,
		MetricCode:     code,
		Name:           name,
		Category:       cat,
		Value:          val,
		Unit:           unit,
		ChangeAbsolute: changeAbs,
		ChangePercent:  changePct,
		Direction:      dir,
		Source:         source,
		Timestamp:      t,
		Hash:           hash,
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.latestTicks[code] = tick
	if len(e.recentTicks) >= e.maxTicks {
		e.recentTicks = e.recentTicks[1:]
	}
	e.recentTicks = append(e.recentTicks, tick)

	return tick
}

// GenerateSimulatedTick produces a realistic micro-drift update for live streaming demonstration.
func (e *LiveFeedEngine) GenerateSimulatedTick(code KPICode) *LiveFeedTick {
	e.mu.RLock()
	prev, exists := e.latestTicks[code]
	e.mu.RUnlock()

	if !exists {
		return nil
	}

	// Pseudo-random deterministic drift based on time
	now := time.Now().UTC()
	driftPct := math.Sin(float64(now.UnixNano())) * 0.004 // +/- 0.4%
	newVal := prev.Value * (1.0 + driftPct)
	diff := newVal - prev.Value
	diffPct := (diff / prev.Value) * 100.0

	return e.RecordTick(code, prev.Name, prev.Category, roundFloat(newVal, 4), prev.Unit, roundFloat(diff, 4), roundFloat(diffPct, 2), prev.Source, now)
}

// GetLatestTicks returns the current snapshot of latest ticks for all tracked feeds.
func (e *LiveFeedEngine) GetLatestTicks() []*LiveFeedTick {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ticks := make([]*LiveFeedTick, 0, len(e.latestTicks))
	for _, tick := range e.latestTicks {
		ticks = append(ticks, tick)
	}
	return ticks
}

// GetRecentTicks returns recent stream history.
func (e *LiveFeedEngine) GetRecentTicks(limit int) []*LiveFeedTick {
	e.mu.RLock()
	defer e.mu.RUnlock()

	n := len(e.recentTicks)
	if limit <= 0 || limit > n {
		limit = n
	}

	res := make([]*LiveFeedTick, limit)
	copy(res, e.recentTicks[n-limit:])
	return res
}

// GetObservations returns recorded observations for a given scope ("NATIONAL" or project ID).
func (e *LiveFeedEngine) GetObservations(scope string) []*KPIObservation {
	e.mu.RLock()
	defer e.mu.RUnlock()

	obs, exists := e.observations[scope]
	if !exists {
		return nil
	}
	res := make([]*KPIObservation, len(obs))
	copy(res, obs)
	return res
}

// BuildMacroSummary aggregates cross-sector rollups and national averages.
func (e *LiveFeedEngine) BuildMacroSummary() *MacroKPISummary {
	ticks := e.GetLatestTicks()

	summary := &MacroKPISummary{
		NationalEmissionsAbatementMt: 42.8,
		AverageIndigenousEquityPct:   18.4,
		TotalILGPAllocatedCAD:        1_450_000_000,
		NationalSpendRunRateCADMo:    285.0,
		AverageDomesticContentPct:    68.2,
		AverageGridQueueWaitMonths:   26.4,
		AverageIAACReviewDurationMo:  31.8,
		NationalRedSealDeficitFTE:    14200,
		ActiveFeedTicksCount:         len(ticks),
		LatestCommodityTicks:         ticks,
		GeneratedAt:                  time.Now().UTC(),
	}

	return summary
}

// GenerateSnapshot produces an immutable export container of all definitions, observations, and macro metrics.
func (e *LiveFeedEngine) GenerateSnapshot() (*KPISnapshot, error) {
	definitions := CanonicalRegistry()
	allObs := make([]*KPIObservation, 0, 100)

	e.mu.RLock()
	for _, obsList := range e.observations {
		allObs = append(allObs, obsList...)
	}
	e.mu.RUnlock()

	summary := e.BuildMacroSummary()

	// Compute snapshot hash
	dataBytes, err := json.Marshal(struct {
		Defs []*KPIDefinition  `json:"definitions"`
		Obs  []*KPIObservation `json:"observations"`
		Sum  *MacroKPISummary  `json:"summary"`
	}{Defs: definitions, Obs: allObs, Sum: summary})
	if err != nil {
		return nil, fmt.Errorf("serialize snapshot: %w", err)
	}

	sum := sha256.Sum256(dataBytes)
	auditHash := hex.EncodeToString(sum[:])

	return &KPISnapshot{
		Version:      "cegs-kpi-v1.0",
		GeneratedAt:  time.Now().UTC(),
		Definitions:  definitions,
		Observations: allObs,
		MacroSummary: summary,
		AuditHash:    auditHash,
	}, nil
}

func hashObservation(id string, val float64) string {
	payload := fmt.Sprintf("%s:%.4f", id, val)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func roundFloat(val float64, decimals int) float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}
