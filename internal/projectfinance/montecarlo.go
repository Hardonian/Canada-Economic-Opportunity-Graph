package projectfinance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// MetricDistribution contains empirical quantiles from a Monte Carlo simulation.
type MetricDistribution struct {
	Mean float64 `json:"mean"`
	P10  float64 `json:"p10"` // Downside / conservative
	P50  float64 `json:"p50"` // Median expectation
	P90  float64 `json:"p90"` // Upside tail
}

// MonteCarloSimulationResult summarizes 10,000 randomized project finance paths.
type MonteCarloSimulationResult struct {
	ProjectID              string             `json:"project_id"`
	ProjectName            string             `json:"project_name"`
	IterationsRun          int                `json:"iterations_run"`
	BaselineCapexCAD       int64              `json:"baseline_capex_cad"`
	ProjectIRRPercent      MetricDistribution `json:"project_irr_percent"`
	EquityIRRPercent       MetricDistribution `json:"equity_irr_percent"`
	MinDSCR                MetricDistribution `json:"min_dscr"` // Minimum Debt Service Coverage Ratio
	AvgDSCR                MetricDistribution `json:"avg_dscr"`
	LoanLifeCoverageRatio  MetricDistribution `json:"loan_life_coverage_ratio"`
	ProbabilityOfDefaultPct float64           `json:"probability_of_default_pct"` // % of runs where DSCR < 1.05
	SyntheticCreditRating  string             `json:"synthetic_credit_rating"`    // "AAA", "AA", "A", "BBB+", "BBB-", "BB", "B"
	InvestmentGrade        bool               `json:"investment_grade"`
	AuditHash              string             `json:"audit_hash"`
}

// FinanceSimulator runs stochastic project finance cash flow modeling.
type FinanceSimulator struct{}

// NewFinanceSimulator creates an instance of the finance simulator.
func NewFinanceSimulator() *FinanceSimulator {
	return &FinanceSimulator{}
}

// RunSimulation executes a fast, deterministic pseudo-random Monte Carlo model.
func (s *FinanceSimulator) RunSimulation(project *domain.Project, iterations int) *MonteCarloSimulationResult {
	if iterations <= 0 {
		iterations = 1000
	}
	if iterations > 10000 {
		iterations = 10000
	}

	capex := project.CapexCAD
	if capex <= 0 {
		capex = 750_000_000
	}

	// Base sector cash flow yield and volatility parameters
	baseEBITDAMargin := 0.28
	baseYield := 0.115
	volatility := 0.20

	switch project.Sector {
	case domain.SectorNuclearEnergy:
		baseYield = 0.095
		volatility = 0.12 // Long-term rate-base stability
	case domain.SectorCriticalMinerals:
		baseYield = 0.155
		volatility = 0.32 // Commodity cyclicality
	case domain.SectorAICompute:
		baseYield = 0.140
		volatility = 0.18
	case domain.SectorTransportation:
		baseYield = 0.100
		volatility = 0.14
	}

	r := rand.New(rand.NewSource(42)) // Deterministic seed for reproducible audit hashes

	projectIRRs := make([]float64, iterations)
	equityIRRs := make([]float64, iterations)
	minDSCRs := make([]float64, iterations)
	avgDSCRs := make([]float64, iterations)
	llcrs := make([]float64, iterations)

	defaultsCount := 0

	// Capital stack assumption: 60% senior debt, 40% equity
	debtAmount := float64(capex) * 0.60
	equityAmount := float64(capex) * 0.40
	annualDebtService := (debtAmount / 20.0) + (debtAmount * 0.052) // 20-year amort at 5.2% blended interest

	for i := 0; i < iterations; i++ {
		// Box-Muller normal deviates for commodity price shock and interest rate drift
		u1 := r.Float64()
		u2 := r.Float64()
		if u1 <= 0 {
			u1 = 0.0001
		}
		z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)

		shockFactor := math.Exp(z * volatility)
		effectiveYield := baseYield * shockFactor
		annualRevenue := float64(capex) * effectiveYield
		annualEBITDA := annualRevenue * baseEBITDAMargin

		dscr := annualEBITDA / annualDebtService
		if dscr < 0.1 {
			dscr = 0.1
		}

		projIRR := (effectiveYield * 100.0) - 2.5
		eqIRR := ((annualEBITDA - annualDebtService) / equityAmount) * 100.0
		llcr := (annualEBITDA * 15.0) / debtAmount

		if dscr < 1.05 {
			defaultsCount++
		}

		projectIRRs[i] = projIRR
		equityIRRs[i] = eqIRR
		minDSCRs[i] = dscr * 0.88 // stress trough
		avgDSCRs[i] = dscr
		llcrs[i] = llcr
	}

	sort.Float64s(projectIRRs)
	sort.Float64s(equityIRRs)
	sort.Float64s(minDSCRs)
	sort.Float64s(avgDSCRs)
	sort.Float64s(llcrs)

	calcDist := func(arr []float64) MetricDistribution {
		sum := 0.0
		for _, v := range arr {
			sum += v
		}
		return MetricDistribution{
			Mean: math.Round((sum/float64(len(arr)))*100) / 100,
			P10:  math.Round(arr[int(float64(len(arr))*0.10)]*100) / 100,
			P50:  math.Round(arr[int(float64(len(arr))*0.50)]*100) / 100,
			P90:  math.Round(arr[int(float64(len(arr))*0.90)]*100) / 100,
		}
	}

	probDefault := (float64(defaultsCount) / float64(iterations)) * 100.0
	avgDSCRMetric := calcDist(avgDSCRs)

	rating := "BBB"
	isIG := true
	if avgDSCRMetric.P50 >= 1.85 && probDefault < 1.0 {
		rating = "AA"
	} else if avgDSCRMetric.P50 >= 1.55 && probDefault < 3.0 {
		rating = "A"
	} else if avgDSCRMetric.P50 >= 1.30 && probDefault < 8.0 {
		rating = "BBB"
	} else if avgDSCRMetric.P50 >= 1.15 {
		rating = "BB+"
		isIG = false
	} else {
		rating = "B"
		isIG = false
	}

	res := &MonteCarloSimulationResult{
		ProjectID:               project.ID,
		ProjectName:             project.Name,
		IterationsRun:           iterations,
		BaselineCapexCAD:        capex,
		ProjectIRRPercent:       calcDist(projectIRRs),
		EquityIRRPercent:        calcDist(equityIRRs),
		MinDSCR:                 calcDist(minDSCRs),
		AvgDSCR:                 avgDSCRMetric,
		LoanLifeCoverageRatio:   calcDist(llcrs),
		ProbabilityOfDefaultPct: math.Round(probDefault*10) / 10,
		SyntheticCreditRating:   rating,
		InvestmentGrade:         isIG,
	}

	auditData := fmt.Sprintf("%s|%d|%d|%.2f|%.2f|%s",
		res.ProjectID, res.BaselineCapexCAD, res.IterationsRun, res.AvgDSCR.P50, res.ProbabilityOfDefaultPct, res.SyntheticCreditRating)
	h := sha256.Sum256([]byte(auditData))
	res.AuditHash = hex.EncodeToString(h[:])

	return res
}
