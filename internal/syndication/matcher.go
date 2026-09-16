package syndication

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// InvestorClass categorizes institutional capital sources.
type InvestorClass string

const (
	InvestorClassMapleEight      InvestorClass = "CANADIAN_PUBLIC_PENSION" // CPPIB, CDPQ, OTPP, OMERS, etc.
	InvestorClassGlobalSWF       InvestorClass = "SOVEREIGN_WEALTH_FUND"   // GIC, Temasek, Norges Bank, ADIA
	InvestorClassCrownConcession InvestorClass = "FEDERAL_CROWN_FUND"      // CIB, SIF, Canada Growth Fund (CGF)
	InvestorClassPrivateInfra    InvestorClass = "INFRASTRUCTURE_PE"       // Brookfield, Northleaf, Fengate
	InvestorClassIndigenous      InvestorClass = "INDIGENOUS_DEVELOPMENT"  // FN Development Corporations + ILGP
)

// InstitutionalInvestor represents an institutional capital allocator with strict mandates.
type InstitutionalInvestor struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Class            InvestorClass `json:"class"`
	Jurisdiction     string        `json:"jurisdiction"` // CA, SG, NO, AE, etc.
	AUMBillionsCAD   float64       `json:"aum_billions_cad"`
	MinTicketCAD     int64         `json:"min_ticket_cad"`
	MaxTicketCAD     int64         `json:"max_ticket_cad"`
	TargetReturnHurdle float64     `json:"target_return_hurdle"` // e.g. 7.5%
	TargetSectors    []domain.Sector `json:"target_sectors"`
	PrefersGreenfield bool         `json:"prefers_greenfield"`
	RequiresDomestic bool          `json:"requires_domestic"`
}

// DefaultInstitutionalInvestors returns the canonical directory of Canadian and global allocators.
func DefaultInstitutionalInvestors() []InstitutionalInvestor {
	return []InstitutionalInvestor{
		{
			ID:                 "cppib",
			Name:               "CPP Investments (CPPIB)",
			Class:              InvestorClassMapleEight,
			Jurisdiction:       "CA",
			AUMBillionsCAD:     632.0,
			MinTicketCAD:       250_000_000,
			MaxTicketCAD:       2_500_000_000,
			TargetReturnHurdle: 8.5,
			TargetSectors:      []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorTransportation, domain.SectorAICompute},
			PrefersGreenfield:  false,
			RequiresDomestic:   false,
		},
		{
			ID:                 "cdpq",
			Name:               "Caisse de dépôt et placement du Québec (CDPQ)",
			Class:              InvestorClassMapleEight,
			Jurisdiction:       "CA",
			AUMBillionsCAD:     434.0,
			MinTicketCAD:       150_000_000,
			MaxTicketCAD:       1_500_000_000,
			TargetReturnHurdle: 8.0,
			TargetSectors:      []domain.Sector{domain.SectorCleanEnergy, domain.SectorTransportation, domain.SectorCriticalMinerals},
			PrefersGreenfield:  true,
			RequiresDomestic:   false,
		},
		{
			ID:                 "otpp",
			Name:               "Ontario Teachers' Pension Plan (OTPP)",
			Class:              InvestorClassMapleEight,
			Jurisdiction:       "CA",
			AUMBillionsCAD:     255.0,
			MinTicketCAD:       100_000_000,
			MaxTicketCAD:       1_000_000_000,
			TargetReturnHurdle: 8.2,
			TargetSectors:      []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorAICompute},
			PrefersGreenfield:  false,
			RequiresDomestic:   false,
		},
		{
			ID:                 "cib",
			Name:               "Canada Infrastructure Bank (CIB)",
			Class:              InvestorClassCrownConcession,
			Jurisdiction:       "CA",
			AUMBillionsCAD:     35.0,
			MinTicketCAD:       50_000_000,
			MaxTicketCAD:       1_500_000_000,
			TargetReturnHurdle: 4.5, // Sub-commercial concessionary debt
			TargetSectors:      []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorTransportation, domain.SectorCriticalMinerals},
			PrefersGreenfield:  true,
			RequiresDomestic:   true,
		},
		{
			ID:                 "cgf",
			Name:               "Canada Growth Fund (CGF)",
			Class:              InvestorClassCrownConcession,
			Jurisdiction:       "CA",
			AUMBillionsCAD:     15.0,
			MinTicketCAD:       50_000_000,
			MaxTicketCAD:       1_000_000_000,
			TargetReturnHurdle: 5.5,
			TargetSectors:      []domain.Sector{domain.SectorCriticalMinerals, domain.SectorCleanEnergy, domain.SectorAICompute},
			PrefersGreenfield:  true,
			RequiresDomestic:   true,
		},
		{
			ID:                 "gic",
			Name:               "GIC Private Limited",
			Class:              InvestorClassGlobalSWF,
			Jurisdiction:       "SG",
			AUMBillionsCAD:     1050.0,
			MinTicketCAD:       300_000_000,
			MaxTicketCAD:       3_000_000_000,
			TargetReturnHurdle: 9.0,
			TargetSectors:      []domain.Sector{domain.SectorTransportation, domain.SectorAICompute, domain.SectorCleanEnergy},
			PrefersGreenfield:  false,
			RequiresDomestic:   false,
		},
	}
}

// InvestorMatch evaluates the fit between an investor and a capital opportunity.
type InvestorMatch struct {
	Investor      InstitutionalInvestor `json:"investor"`
	MatchScore    float64               `json:"match_score"` // 0 - 100
	RecommendedTranche string           `json:"recommended_tranche"` // "Lead Equity", "Concessionary Subordinated Debt", "Co-Investment Equity"
	ProposedTicketCAD  int64            `json:"proposed_ticket_cad"`
	Rationale          string           `json:"rationale"`
}

// SyndicationConsortium models an optimal multi-party capital stack solution for a project.
type SyndicationConsortium struct {
	ProjectID          string          `json:"project_id"`
	ProjectName        string          `json:"project_name"`
	TotalCapexCAD      int64           `json:"total_capex_cad"`
	Matches            []InvestorMatch `json:"matches"`
	EquityTrancheCAD   int64           `json:"equity_tranche_cad"`
	DebtTrancheCAD     int64           `json:"debt_tranche_cad"`
	CrownConcessionCAD int64           `json:"crown_concession_cad"`
	IndigenousEquityCAD int64          `json:"indigenous_equity_cad"`
	PrivateCrowdingInRatio float64     `json:"private_crowding_in_ratio"`
	AuditHash          string          `json:"audit_hash"`
}

// Matcher runs institutional capital matching algorithms.
type Matcher struct {
	investors []InstitutionalInvestor
}

// NewMatcher creates an institutional deal matcher.
func NewMatcher(investors []InstitutionalInvestor) *Matcher {
	if len(investors) == 0 {
		investors = DefaultInstitutionalInvestors()
	}
	return &Matcher{investors: investors}
}

// MatchProject identifies the highest-affinity institutional syndicate for a project.
func (m *Matcher) MatchProject(project *domain.Project) *SyndicationConsortium {
	matches := make([]InvestorMatch, 0)
	capex := project.CapexCAD
	if capex == 0 {
		capex = 500_000_000 // Default reference baseline
	}

	for _, inv := range m.investors {
		score := 50.0

		// Sector match check
		sectorMatch := false
		for _, s := range inv.TargetSectors {
			if s == project.Sector {
				sectorMatch = true
				break
			}
		}
		if sectorMatch {
			score += 30.0
		} else {
			score -= 20.0
		}

		// Stage match check
		if inv.PrefersGreenfield && (project.CurrentStage == domain.StageConstruction || project.CurrentStage == domain.StageFEED) {
			score += 10.0
		} else if !inv.PrefersGreenfield && project.CurrentStage == domain.StageOperating {
			score += 15.0
		}

		// Ticket size suitability
		targetTicket := int64(float64(capex) * 0.25)
		if targetTicket < inv.MinTicketCAD {
			score -= 15.0
			targetTicket = inv.MinTicketCAD
		} else if targetTicket > inv.MaxTicketCAD {
			targetTicket = inv.MaxTicketCAD
		}

		if score > 100.0 {
			score = 100.0
		} else if score < 0.0 {
			score = 0.0
		}

		tranche := "Co-Investment Equity"
		rationale := fmt.Sprintf("High sector alignment with %s mandate.", inv.Name)
		if inv.Class == InvestorClassCrownConcession {
			tranche = "Concessionary Subordinated Debt"
			rationale = "Provides catalytic sovereign de-risking debt tranche to crowd in institutional capital."
		} else if inv.Class == InvestorClassMapleEight && score >= 80 {
			tranche = "Lead Sponsor Equity"
			rationale = "Anchors long-term domestic institutional ownership."
		}

		matches = append(matches, InvestorMatch{
			Investor:           inv,
			MatchScore:         score,
			RecommendedTranche: tranche,
			ProposedTicketCAD:  targetTicket,
			Rationale:          rationale,
		})
	}

	// Sort matches descending by affinity score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchScore > matches[j].MatchScore
	})

	equityTranche := int64(float64(capex) * 0.35)
	crownConcession := int64(float64(capex) * 0.20)
	indigenousEquity := int64(float64(capex) * 0.10)
	debtTranche := capex - equityTranche - crownConcession - indigenousEquity
	if debtTranche < 0 {
		debtTranche = 0
	}

	multiplier := 0.0
	if crownConcession > 0 {
		multiplier = float64(equityTranche+debtTranche) / float64(crownConcession)
	}

	consortium := &SyndicationConsortium{
		ProjectID:              project.ID,
		ProjectName:            project.Name,
		TotalCapexCAD:          capex,
		Matches:                matches,
		EquityTrancheCAD:       equityTranche,
		DebtTrancheCAD:         debtTranche,
		CrownConcessionCAD:     crownConcession,
		IndigenousEquityCAD:    indigenousEquity,
		PrivateCrowdingInRatio: multiplier,
	}

	auditData := fmt.Sprintf("%s|%d|%d|%d|%d|%.2f",
		consortium.ProjectID, consortium.TotalCapexCAD, consortium.EquityTrancheCAD,
		consortium.DebtTrancheCAD, consortium.CrownConcessionCAD, consortium.PrivateCrowdingInRatio)
	auditHash := sha256.Sum256([]byte(auditData))
	consortium.AuditHash = hex.EncodeToString(auditHash[:])

	return consortium
}
