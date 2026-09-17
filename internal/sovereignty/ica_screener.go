package sovereignty

import (
	"fmt"
	"strings"
)

// ICAScreeningRequest parameters for Investment Canada Act national security screening.
type ICAScreeningRequest struct {
	ProjectID                string  `json:"project_id"`
	Sector                   string  `json:"sector"`
	AcquiringEntity          string  `json:"acquiring_entity"`
	AcquiringJurisdiction    string  `json:"acquiring_jurisdiction"`
	IsStateOwned             bool    `json:"is_state_owned"`
	AcquisitionSharePct      float64 `json:"acquisition_share_pct"`
	IncludesCriticalMinerals bool    `json:"includes_critical_minerals"`
	IncludesDualUseTech      bool    `json:"includes_dual_use_tech"`
}

// ICAScreeningVerdict holds the statutory national security evaluation.
type ICAScreeningVerdict struct {
	Verdict          string   `json:"verdict"` // "CLEAR", "MANDATORY_NATIONAL_SECURITY_REVIEW", "PROHIBITED"
	StatutorySection string   `json:"statutory_section"`
	RiskFactors      []string `json:"risk_factors"`
	SovereigntyScore float64  `json:"sovereignty_score"` // 0.0 - 100.0 (Higher = safer for Canada)
	Recommendation   string   `json:"recommendation"`
}

// ICAScreener executes statutory reviews under Section 25.3 of the Investment Canada Act.
type ICAScreener struct{}

// NewICAScreener creates an ICA screener.
func NewICAScreener() *ICAScreener {
	return &ICAScreener{}
}

// ScreenAcquisition determines if a proposed transaction triggers a mandatory Order in Council review.
func (s *ICAScreener) ScreenAcquisition(req ICAScreeningRequest) *ICAScreeningVerdict {
	var risks []string
	juris := strings.ToUpper(req.AcquiringJurisdiction)

	isHostileOrNonFTA := juris == "CN" || juris == "CHINA" || juris == "RU" || juris == "RUSSIA" || juris == "IR" || juris == "IRAN"

	if req.IsStateOwned {
		risks = append(risks, "Acquiring entity is a foreign State-Owned Enterprise (SOE) governed by foreign state directives")
	}

	if req.IncludesCriticalMinerals && req.AcquisitionSharePct >= 10.0 {
		risks = append(risks, "Critical Minerals Policy Directive: Any non-FTA SOE investment >=10% triggers automatic national security review")
	}

	if req.IncludesDualUseTech {
		risks = append(risks, "Dual-use technology, advanced AI compute, or quantum intellectual property transfer risk")
	}

	if isHostileOrNonFTA {
		risks = append(risks, fmt.Sprintf("Acquiring jurisdiction %s has active foreign interference or espionage advisories with CSIS", req.AcquiringJurisdiction))
	}

	verdict := "CLEAR"
	statutorySection := "Section 11 (General Notification)"
	score := 85.0
	recommendation := "Standard net benefit and general notification."

	if req.IsStateOwned && req.IncludesCriticalMinerals && isHostileOrNonFTA {
		verdict = "PROHIBITED"
		statutorySection = "Section 25.4 (Order in Council Divestment / Blocking Order)"
		score = 12.0
		recommendation = "Recommend Cabinet Minister issue formal prohibition order blocking acquisition."
	} else if len(risks) > 0 {
		verdict = "MANDATORY_NATIONAL_SECURITY_REVIEW"
		statutorySection = "Section 25.3 (National Security Review Directive)"
		score = 45.0
		recommendation = "Refer to Governor in Council for 200-day national security review with CSIS/CSE."
	}

	return &ICAScreeningVerdict{
		Verdict:          verdict,
		StatutorySection: statutorySection,
		RiskFactors:      risks,
		SovereigntyScore: score,
		Recommendation:   recommendation,
	}
}
