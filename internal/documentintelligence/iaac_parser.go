package documentintelligence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DutyToConsultLevel classifies depth of constitutional obligation under Haida Nation.
type DutyToConsultLevel string

const (
	ConsultNoticeOnly            DutyToConsultLevel = "NOTICE_ONLY"
	ConsultSeriousConsultation   DutyToConsultLevel = "SERIOUS_CONSULTATION"
	ConsultDeepAccommodation     DutyToConsultLevel = "DEEP_ACCOMMODATION_AND_CONSENT"
)

// LitigationRiskRating evaluates vulnerability to Federal Court judicial review injunctions.
type LitigationRiskRating string

const (
	LitigationLow      LitigationRiskRating = "LOW"
	LitigationModerate LitigationRiskRating = "MODERATE"
	LitigationHigh     LitigationRiskRating = "HIGH"
	LitigationCritical LitigationRiskRating = "CRITICAL"
)

// IAACDecisionStatement represents parsed statutory approval terms and environmental covenants.
type IAACDecisionStatement struct {
	ProjectName                      string               `json:"project_name"`
	MinisterName                     string               `json:"minister_name"`
	Decision                         string               `json:"decision"` // e.g. "APPROVAL_WITH_CONDITIONS"
	StatutoryAuthority               string               `json:"statutory_authority"`
	TotalConditionsCount             int                  `json:"total_conditions_count"`
	IndigenousConsultationConditions int                  `json:"indigenous_consultation_conditions"`
	AquaticHabitatConditions         int                  `json:"aquatic_habitat_conditions"`
	WildlifeCaribouConditions        int                  `json:"wildlife_caribou_conditions"`
	AirAndGHGConditions              int                  `json:"air_and_ghg_conditions"`
	FinancialAssuranceBondCAD        int64                `json:"financial_assurance_bond_cad"`
	DutyToConsultDepth               DutyToConsultLevel   `json:"duty_to_consult_depth"`
	InjunctionVulnerabilityScore     float64              `json:"injunction_vulnerability_score"` // 0.0 - 100.0
	LitigationExposure               LitigationRiskRating `json:"litigation_exposure"`
	LegalRiskSummary                 string               `json:"legal_risk_summary"`
	AuditHash                        string               `json:"audit_hash"`
	ExtractedAt                      time.Time            `json:"extracted_at"`
}

var (
	condTotalRegex    = regexp.MustCompile(`(?i)(?:total\s+)?(?:conditions|enforceable conditions)[:\s]+([0-9]+)`)
	indigCondRegex    = regexp.MustCompile(`(?i)(?:indigenous|first nations?|consultation)(?:\s+conditions)?[:\s]+([0-9]+)`)
	aquaticCondRegex  = regexp.MustCompile(`(?i)(?:aquatic|fish|water quality)(?:\s+conditions)?[:\s]+([0-9]+)`)
	wildlifeCondRegex = regexp.MustCompile(`(?i)(?:wildlife|caribou|migratory bird)(?:\s+conditions)?[:\s]+([0-9]+)`)
	bondRegex        = regexp.MustCompile(`(?i)(?:bond|financial assurance|security deposit)[:\s]+(?:C\$|CAD|\$)?\s*([0-9]+(?:\.[0-9]+)?)\s*([BM])`)
)

// ExtractIAACDecisionStatement parses an Impact Assessment Agency of Canada decision text.
func ExtractIAACDecisionStatement(text string) (*IAACDecisionStatement, error) {
	lines := strings.Split(text, "\n")
	fields := make(map[string]string)
	for _, l := range lines {
		if parts := strings.SplitN(strings.TrimSpace(l), ":", 2); len(parts) == 2 {
			fields[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
	}

	stmt := &IAACDecisionStatement{
		ProjectName:        first(fields, "project", "project name", "designation"),
		MinisterName:       first(fields, "minister", "issuing authority"),
		Decision:           first(fields, "decision", "determination"),
		StatutoryAuthority: "Section 54, Impact Assessment Act (S.C. 2019, c. 28, s. 1)",
		ExtractedAt:        time.Now().UTC(),
	}

	if stmt.Decision == "" {
		stmt.Decision = "APPROVAL_WITH_CONDITIONS"
	}

	if match := condTotalRegex.FindStringSubmatch(text); len(match) > 1 {
		stmt.TotalConditionsCount, _ = strconv.Atoi(match[1])
	}
	if match := indigCondRegex.FindStringSubmatch(text); len(match) > 1 {
		stmt.IndigenousConsultationConditions, _ = strconv.Atoi(match[1])
	}
	if match := aquaticCondRegex.FindStringSubmatch(text); len(match) > 1 {
		stmt.AquaticHabitatConditions, _ = strconv.Atoi(match[1])
	}
	if match := wildlifeCondRegex.FindStringSubmatch(text); len(match) > 1 {
		stmt.WildlifeCaribouConditions, _ = strconv.Atoi(match[1])
	}

	if match := bondRegex.FindStringSubmatch(text); len(match) > 2 {
		val, _ := strconv.ParseFloat(match[1], 64)
		mult := 1_000_000.0
		if strings.EqualFold(match[2], "b") {
			mult = 1_000_000_000.0
		}
		stmt.FinancialAssuranceBondCAD = int64(math.Round(val * mult))
	}

	// Assess Duty to Consult depth based on Haida Nation criteria
	lower := strings.ToLower(text)
	hasTreatyRights := strings.Contains(lower, "aboriginal title") || strings.Contains(lower, "treaty rights") || strings.Contains(lower, "unceded")
	hasOpposition := strings.Contains(lower, "unresolved objection") || strings.Contains(lower, "without consent") || strings.Contains(lower, "dissenting nation")
	hasIBA := strings.Contains(lower, "impact benefit agreement") || strings.Contains(lower, "iba signed") || strings.Contains(lower, "consent agreement")

	var depth DutyToConsultLevel
	var vulnScore float64

	switch {
	case hasTreatyRights || stmt.IndigenousConsultationConditions >= 10:
		depth = ConsultDeepAccommodation
		vulnScore = 55.0
	case stmt.IndigenousConsultationConditions >= 4:
		depth = ConsultSeriousConsultation
		vulnScore = 35.0
	default:
		depth = ConsultNoticeOnly
		vulnScore = 15.0
	}

	if hasOpposition {
		vulnScore += 35.0
	}
	if hasIBA {
		vulnScore -= 30.0
	}
	if vulnScore > 100.0 {
		vulnScore = 100.0
	} else if vulnScore < 0.0 {
		vulnScore = 0.0
	}

	var exposure LitigationRiskRating
	switch {
	case vulnScore >= 75.0:
		exposure = LitigationCritical
	case vulnScore >= 50.0:
		exposure = LitigationHigh
	case vulnScore >= 25.0:
		exposure = LitigationModerate
	default:
		exposure = LitigationLow
	}

	stmt.DutyToConsultDepth = depth
	stmt.InjunctionVulnerabilityScore = math.Round(vulnScore*10) / 10
	stmt.LitigationExposure = exposure
	stmt.LegalRiskSummary = fmt.Sprintf("Governed by %s obligations. Injunction exposure rated %s (Score: %.1f/100) across %d enforceable statutory conditions.",
		depth, exposure, stmt.InjunctionVulnerabilityScore, stmt.TotalConditionsCount)

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d|%.1f", stmt.ProjectName, stmt.Decision, stmt.TotalConditionsCount, stmt.InjunctionVulnerabilityScore)))
	stmt.AuditHash = hex.EncodeToString(h[:])

	return stmt, nil
}
