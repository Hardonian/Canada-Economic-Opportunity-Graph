package nationalplanning

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// IndigenousCoInvestmentModel details equity ownership, concessionary debt, and long-term community cash distributions.
type IndigenousCoInvestmentModel struct {
	ProjectID                   string    `json:"project_id"`
	ProjectName                 string    `json:"project_name"`
	TotalCapexCAD               int64     `json:"total_capex_cad"`
	EquityOwnershipPercent      float64   `json:"equity_ownership_percent"` // e.g. 25.0% - 50.0%
	IndigenousEquityCAD         int64     `json:"indigenous_equity_cad"`
	CIBLoanAmountCAD            int64     `json:"cib_loan_amount_cad"`
	BandDirectEquityCAD         int64     `json:"band_direct_equity_cad"`
	AnnualGrossProjectYieldCAD  int64     `json:"annual_gross_project_yield_cad"`
	IndigenousShareGrossYield   int64     `json:"indigenous_share_gross_yield_cad"`
	AnnualCIBDebtServiceCAD     int64     `json:"annual_cib_debt_service_cad"`
	AnnualNetCommunityCashYield int64     `json:"annual_net_community_cash_yield_cad"`
	ThirtyYearCumulativeWealth  int64     `json:"thirty_year_cumulative_wealth_cad"`
	FinancialSelfDetermination  string    `json:"financial_self_determination_rating"` // MODEST, STRONG, GENERATIONAL_TRANSFORMATIVE
	AuditHash                   string    `json:"audit_hash"`
	ModelledAt                  time.Time `json:"modelled_at"`
}

// IndigenousModeler evaluates First Nations equity participation scenarios.
type IndigenousModeler struct{}

func NewIndigenousModeler() *IndigenousModeler {
	return &IndigenousModeler{}
}

// ModelEquityParticipation computes debt service and multi-generational sovereign wealth flows.
func (im *IndigenousModeler) ModelEquityParticipation(project *domain.Project, equityPct float64) *IndigenousCoInvestmentModel {
	if equityPct <= 0 {
		equityPct = 25.0 // Default 25% co-ownership benchmark
	}

	capex := float64(project.CapexCAD)
	equityShare := capex * (equityPct / 100.0)

	// Under Canada Infrastructure Bank (CIB) Indigenous Equity Loan Program:
	// Up to 90% of equity purchase can be financed via concessionary GoC borrowing rate + 50 bps (e.g. 4.0%)
	cibLoan := equityShare * 0.90
	bandCash := equityShare * 0.10

	// Typical Canadian industrial project EBITDA yield: 11% on total CAPEX
	annualGrossYield := capex * 0.11
	indigGrossYield := annualGrossYield * (equityPct / 100.0)

	// CIB 25-year amortizing facility @ 4.0%
	rate := 0.040
	n := 25.0
	annualDebtService := cibLoan * (rate / (1.0 - math.Pow(1.0+rate, -n)))

	netCommunityCash := indigGrossYield - annualDebtService
	if netCommunityCash < 0 {
		netCommunityCash = 0
	}

	// 30-year cumulative return (25 years paying debt + net cash, 5 years with 0 debt service)
	cumulative30Year := (netCommunityCash * 25) + (indigGrossYield * 5)

	rating := "MODEST"
	switch {
	case cumulative30Year >= 1_000_000_000:
		rating = "GENERATIONAL_TRANSFORMATIVE"
	case cumulative30Year >= 250_000_000:
		rating = "STRONG"
	}

	model := &IndigenousCoInvestmentModel{
		ProjectID:                   project.ID,
		ProjectName:                 project.Name,
		TotalCapexCAD:               project.CapexCAD,
		EquityOwnershipPercent:      equityPct,
		IndigenousEquityCAD:         int64(math.Round(equityShare)),
		CIBLoanAmountCAD:            int64(math.Round(cibLoan)),
		BandDirectEquityCAD:         int64(math.Round(bandCash)),
		AnnualGrossProjectYieldCAD:  int64(math.Round(annualGrossYield)),
		IndigenousShareGrossYield:   int64(math.Round(indigGrossYield)),
		AnnualCIBDebtServiceCAD:     int64(math.Round(annualDebtService)),
		AnnualNetCommunityCashYield: int64(math.Round(netCommunityCash)),
		ThirtyYearCumulativeWealth:  int64(math.Round(cumulative30Year)),
		FinancialSelfDetermination:  rating,
		ModelledAt:                  time.Now().UTC(),
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%.1f|%d|%d", model.ProjectID, model.EquityOwnershipPercent, model.IndigenousEquityCAD, model.ThirtyYearCumulativeWealth)))
	model.AuditHash = hex.EncodeToString(h[:])

	return model
}
