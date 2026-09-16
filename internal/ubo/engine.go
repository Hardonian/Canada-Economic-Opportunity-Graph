package ubo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// ICARiskLevel defines the Investment Canada Act national security risk posture.
type ICARiskLevel string

const (
	RiskClear           ICARiskLevel = "CLEAR"
	RiskWatchlist       ICARiskLevel = "WATCHLIST"
	RiskMandatoryReview ICARiskLevel = "MANDATORY_REVIEW"
	RiskProhibited      ICARiskLevel = "PROHIBITED"
)

// BeneficialOwner represents a resolved direct or indirect shareholder/controller.
type BeneficialOwner struct {
	EntityID         string  `json:"entity_id"`
	Name             string  `json:"name"`
	Jurisdiction     string  `json:"jurisdiction"`
	OwnershipPercent float64 `json:"ownership_percent"`
	IsSOE            bool    `json:"is_soe"` // State-Owned Enterprise
	IsFTA            bool    `json:"is_fta"` // Canada Free Trade Agreement Partner (USMCA, CPTPP, CETA)
	ControlMechanism string  `json:"control_mechanism"` // Equity, DebtCovenant, OfftakeOption, BoardSeat
}

// SovereignScreeningResult provides audit-grade analysis of foreign control and ICA exposure.
type SovereignScreeningResult struct {
	ProjectID             string            `json:"project_id"`
	ProponentID           string            `json:"proponent_id"`
	ProponentName         string            `json:"proponent_name"`
	UltimateOwners        []BeneficialOwner `json:"ultimate_owners"`
	DomesticControlShare  float64           `json:"domestic_control_share"`  // 0.0 - 1.0 (Canada)
	FTAPartnerShare       float64           `json:"fta_partner_share"`       // 0.0 - 1.0 (USMCA, CETA, etc.)
	NonFTAShare           float64           `json:"non_fta_share"`           // 0.0 - 1.0
	SOEExposurePercent    float64           `json:"soe_exposure_percent"`    // State-Owned Enterprise %
	ICARisk               ICARiskLevel      `json:"ica_risk"`                // CLEAR, WATCHLIST, MANDATORY_REVIEW, PROHIBITED
	CriticalMineralFlag   bool              `json:"critical_mineral_flag"`
	DualUseSovereignty    bool              `json:"dual_use_sovereignty"`
	NationalSecurityNotes []string          `json:"national_security_notes"`
	AuditHash             string            `json:"audit_hash"`
	ScreenedAt            time.Time         `json:"screened_at"`
}

var ftaJurisdictions = map[string]bool{
	"CA": true, "CAN": true, "CANADA": true,
	"US": true, "USA": true, "UNITED STATES": true,
	"MX": true, "MEX": true, "MEXICO": true,
	"GB": true, "GBR": true, "UNITED KINGDOM": true,
	"JP": true, "JPN": true, "JAPAN": true,
	"DE": true, "DEU": true, "GERMANY": true,
	"FR": true, "FRA": true, "FRANCE": true,
	"AU": true, "AUS": true, "AUSTRALIA": true,
	"KR": true, "KOR": true, "SOUTH KOREA": true,
	"NO": true, "NOR": true, "NORWAY": true,
	"SE": true, "SWE": true, "SWEDEN": true,
	"NL": true, "NLD": true, "NETHERLANDS": true,
	"EU": true,
}

var nonFTAHighRiskJurisdictions = map[string]bool{
	"CN": true, "CHN": true, "CHINA": true,
	"RU": true, "RUS": true, "RUSSIA": true,
	"IR": true, "IRN": true, "IRAN": true,
	"KP": true, "PRK": true, "NORTH KOREA": true,
	"BY": true, "BLR": true, "BELARUS": true,
}

// Evaluator performs deterministic UBO screening and Investment Canada Act classification.
type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// ScreenProject evaluates foreign control, state-owned enterprise risk, and ICA review status.
func (e *Evaluator) ScreenProject(project *domain.Project, owners []BeneficialOwner) *SovereignScreeningResult {
	result := &SovereignScreeningResult{
		ProjectID:             project.ID,
		ProponentID:           project.ProponentID,
		ScreenedAt:            time.Now().UTC(),
		NationalSecurityNotes: make([]string, 0),
		UltimateOwners:        owners,
	}

	if project.Proponent != nil {
		result.ProponentName = project.Proponent.LegalName
	} else {
		result.ProponentName = project.ProponentID
	}

	// Classify Sector Sensitivity
	if project.Sector == domain.SectorCriticalMinerals || strings.Contains(strings.ToLower(project.Subsector), "lithium") ||
		strings.Contains(strings.ToLower(project.Subsector), "nickel") || strings.Contains(strings.ToLower(project.Subsector), "rare earth") ||
		strings.Contains(strings.ToLower(project.Subsector), "uranium") {
		result.CriticalMineralFlag = true
	}

	if project.Sector == domain.SectorNuclearEnergy || project.Sector == domain.SectorDefenceArctic ||
		project.Sector == domain.SectorAICompute || strings.Contains(strings.ToLower(project.LocationName), "arctic") ||
		strings.Contains(strings.ToLower(project.Summary), "dual-use") {
		result.DualUseSovereignty = true
	}

	if len(owners) == 0 {
		// Default to pure domestic ownership if no external owners provided
		result.DomesticControlShare = 1.0
		result.ICARisk = RiskClear
		result.NationalSecurityNotes = append(result.NationalSecurityNotes, "100% Canadian domestic proponent control identified.")
		result.AuditHash = calculateHash(result)
		return result
	}

	var domesticShare, ftaShare, nonFTAShare, soeShare float64
	var nonFTASOEFound bool

	for _, o := range owners {
		juris := strings.ToUpper(strings.TrimSpace(o.Jurisdiction))
		isCA := juris == "CA" || juris == "CAN" || juris == "CANADA"

		if isCA {
			domesticShare += o.OwnershipPercent
		} else if ftaJurisdictions[juris] {
			ftaShare += o.OwnershipPercent
		} else {
			nonFTAShare += o.OwnershipPercent
		}

		if o.IsSOE {
			soeShare += o.OwnershipPercent
			if nonFTAHighRiskJurisdictions[juris] {
				nonFTASOEFound = true
			}
		}
	}

	result.DomesticControlShare = domesticShare
	result.FTAPartnerShare = ftaShare
	result.NonFTAShare = nonFTAShare
	result.SOEExposurePercent = soeShare

	// Determine Investment Canada Act (ICA) Review Tier
	switch {
	case nonFTASOEFound && (result.CriticalMineralFlag || result.DualUseSovereignty):
		result.ICARisk = RiskProhibited
		result.NationalSecurityNotes = append(result.NationalSecurityNotes,
			"CRITICAL: Non-FTA State-Owned Enterprise (SOE) stake in strategic mineral or dual-use asset triggers ICA Policy on Foreign SOE Investments in Critical Minerals (Order in Council Prohibition).")
	case result.SOEExposurePercent > 0.10 && result.CriticalMineralFlag:
		result.ICARisk = RiskMandatoryReview
		result.NationalSecurityNotes = append(result.NationalSecurityNotes,
			"Mandatory Part IV.1 National Security Review required under Investment Canada Act (>10% SOE ownership in critical minerals).")
	case result.NonFTAShare > 0.25 || nonFTASOEFound:
		result.ICARisk = RiskWatchlist
		result.NationalSecurityNotes = append(result.NationalSecurityNotes,
			"Elevated non-FTA foreign ownership threshold detected (>25%). Enhanced scrutiny on board appointment rights and off-take restrictions.")
	default:
		result.ICARisk = RiskClear
		result.NationalSecurityNotes = append(result.NationalSecurityNotes,
			"Ownership structure compliant with standard ICA notification thresholds and Five Eyes / allied partner treaties.")
	}

	result.AuditHash = calculateHash(result)
	return result
}

func calculateHash(r *SovereignScreeningResult) string {
	raw := fmt.Sprintf("%s|%s|%.2f|%.2f|%.2f|%.2f|%s",
		r.ProjectID, r.ProponentID, r.DomesticControlShare, r.FTAPartnerShare, r.NonFTAShare, r.SOEExposurePercent, r.ICARisk)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
