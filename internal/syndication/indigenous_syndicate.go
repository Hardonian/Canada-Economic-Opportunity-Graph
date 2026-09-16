package syndication

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// FirstNationParticipant models a community equity partner in a regional infrastructure asset.
type FirstNationParticipant struct {
	BandCouncilName     string  `json:"band_council_name"`
	TreatyOrTerritory   string  `json:"treaty_or_territory"`
	CorridorKilometers  float64 `json:"corridor_kilometers"`
	EquitySharePercent  float64 `json:"equity_share_percent"`
	GuaranteedDebtCAD   int64   `json:"guaranteed_debt_cad"`
	AnnualDividendCAD   int64   `json:"annual_dividend_cad"`
	Cumulative30YrCAD   int64   `json:"cumulative_30yr_cad"`
}

// MultiNationSyndicate models co-ownership distribution across an infrastructure corridor.
type MultiNationSyndicate struct {
	CorridorProjectName    string                   `json:"corridor_project_name"`
	TotalCorridorKM        float64                  `json:"total_corridor_km"`
	TotalEquityValueCAD    int64                    `json:"total_equity_value_cad"`
	FederalILGPGreaterCAD  int64                    `json:"federal_ilgp_guarantee_cad"`
	Participants           []FirstNationParticipant `json:"participants"`
	BlendedInterestSpreadBps int                    `json:"blended_interest_spread_bps"` // e.g. 235 bps savings
	TotalAnnualDividendsCAD int64                   `json:"total_annual_dividends_cad"`
	Total30YearWealthCAD   int64                    `json:"total_30_year_wealth_cad"`
	AuditHash              string                   `json:"audit_hash"`
}

// BuildMultiNationSyndicate computes deterministic syndication shares proportional to territory traversed.
func BuildMultiNationSyndicate(projectName string, totalEquityCAD int64, communities []struct {
	Name      string
	Territory string
	KM        float64
}) *MultiNationSyndicate {
	var totalKM float64
	for _, c := range communities {
		totalKM += c.KM
	}
	if totalKM == 0 {
		totalKM = 1.0
	}

	participants := make([]FirstNationParticipant, len(communities))
	var totalAnnualDiv int64
	var total30YrWealth int64
	ilgpGuarantee := int64(float64(totalEquityCAD) * 0.95) // 95% guaranteed tranche

	for i, c := range communities {
		sharePct := (c.KM / totalKM) * 100.0
		allocatedDebt := int64(float64(ilgpGuarantee) * (c.KM / totalKM))
		
		// Return assumptions: 9.5% gross yield, debt service amortized over 20 years at ~3.75% sovereign base
		grossReturn := float64(totalEquityCAD) * (c.KM / totalKM) * 0.095
		debtService := (float64(allocatedDebt) / 20.0) + (float64(allocatedDebt) * 0.0375)
		netAnnual := int64(grossReturn - debtService)
		if netAnnual < 0 {
			netAnnual = int64(grossReturn * 0.25) // minimum floor distribution
		}
		cum30Yr := (netAnnual * 20) + int64(grossReturn * 10) // 10 years unencumbered post-debt amortization

		participants[i] = FirstNationParticipant{
			BandCouncilName:    c.Name,
			TreatyOrTerritory:  c.Territory,
			CorridorKilometers: c.KM,
			EquitySharePercent: sharePct,
			GuaranteedDebtCAD:  allocatedDebt,
			AnnualDividendCAD:  netAnnual,
			Cumulative30YrCAD:  cum30Yr,
		}

		totalAnnualDiv += netAnnual
		total30YrWealth += cum30Yr
	}

	syndicate := &MultiNationSyndicate{
		CorridorProjectName:      projectName,
		TotalCorridorKM:          totalKM,
		TotalEquityValueCAD:      totalEquityCAD,
		FederalILGPGreaterCAD:    ilgpGuarantee,
		Participants:             participants,
		BlendedInterestSpreadBps: 235,
		TotalAnnualDividendsCAD:  totalAnnualDiv,
		Total30YearWealthCAD:     total30YrWealth,
	}

	auditData := fmt.Sprintf("%s|%d|%d|%.2f|%d", projectName, totalEquityCAD, ilgpGuarantee, totalKM, len(communities))
	h := sha256.Sum256([]byte(auditData))
	syndicate.AuditHash = hex.EncodeToString(h[:])

	return syndicate
}
