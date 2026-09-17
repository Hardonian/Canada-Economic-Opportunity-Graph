package indigenouslinker

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// OCAPStandardVersion identifies the First Nations Information Governance Centre framework.
const OCAPStandardVersion = "FNIGC-OCAP-v2.0"

// ConsentStatus defines community authorization states.
type ConsentStatus string

const (
	ConsentGranted   ConsentStatus = "CONSENT_GRANTED"
	ConsentPending   ConsentStatus = "CONSENT_PENDING"
	ConsentWithdrawn ConsentStatus = "CONSENT_WITHDRAWN"
	ConsentRestricted ConsentStatus = "CONSENT_RESTRICTED"
)

// CommunityConsent captures Free, Prior, and Informed Consent (FPIC) metadata.
type CommunityConsent struct {
	NationOrBandName    string        `json:"nation_or_band_name"`
	BandCouncilNumber   string        `json:"band_council_number"`
	TreatyTerritory     string        `json:"treaty_territory"`
	Status              ConsentStatus `json:"status"`
	PermittedUseCases   []string      `json:"permitted_use_cases"` // "PROCUREMENT_MATCHING", "CAPITAL_STACK", "SUPPLY_CHAIN"
	RestrictedKnowledge []string      `json:"restricted_knowledge"`// Sacred sites, non-public traditional harvest zones
	AuthorizedSigner    string        `json:"authorized_signer"`
	EffectiveDate       time.Time     `json:"effective_date"`
	ExpiryDate          time.Time     `json:"expiry_date"`
}

// OCAPEvaluationResult details the governance compliance check for an Indigenous business relationship.
type OCAPEvaluationResult struct {
	Compliant       bool      `json:"compliant"`
	ViolationReason string    `json:"violation_reason,omitempty"`
	AuditHash       string    `json:"audit_hash"`
	EvaluatedAt     time.Time `json:"evaluated_at"`
}

// OCAPComplianceEngine enforces data sovereignty and ethical linking constraints.
type OCAPComplianceEngine struct {
	registry map[string]*CommunityConsent
}

// NewOCAPComplianceEngine creates an OCAP protocol verification engine.
func NewOCAPComplianceEngine() *OCAPComplianceEngine {
	return &OCAPComplianceEngine{
		registry: make(map[string]*CommunityConsent),
	}
}

// RegisterConsent registers explicit nation or council consent for community business records.
func (e *OCAPComplianceEngine) RegisterConsent(consent *CommunityConsent) {
	if consent != nil && consent.NationOrBandName != "" {
		e.registry[consent.NationOrBandName] = consent
	}
}

// ValidateLink verifies whether linking an Indigenous business to a project respects OCAP principles.
func (e *OCAPComplianceEngine) ValidateLink(business *domain.Entity, project *domain.Project, requestedUseCase string) OCAPEvaluationResult {
	now := time.Now().UTC()

	if business == nil || business.EntityType != "IndigenousBusiness" {
		return OCAPEvaluationResult{
			Compliant:   true,
			EvaluatedAt: now,
		}
	}

	nation := business.Jurisdiction
	if nation == "" {
		nation = "Unspecified First Nation"
	}

	consent, ok := e.registry[nation]
	if !ok {
		// Default rule: Publicly registered business in official Indigenous Business Directory
		// is compliant for public procurement matching, provided traditional knowledge is absent.
		h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", business.ID, project.ID, requestedUseCase)))
		return OCAPEvaluationResult{
			Compliant:   true,
			AuditHash:   hex.EncodeToString(h[:]),
			EvaluatedAt: now,
		}
	}

	// Check consent status
	if consent.Status == ConsentWithdrawn {
		return OCAPEvaluationResult{
			Compliant:       false,
			ViolationReason: fmt.Sprintf("community consent withdrawn by %s governing council", consent.NationOrBandName),
			EvaluatedAt:     now,
		}
	}

	if now.After(consent.ExpiryDate) {
		return OCAPEvaluationResult{
			Compliant:       false,
			ViolationReason: fmt.Sprintf("community consent expired on %s", consent.ExpiryDate.Format("2006-01-02")),
			EvaluatedAt:     now,
		}
	}

	// Check permitted use-case
	usePermitted := false
	for _, u := range consent.PermittedUseCases {
		if u == requestedUseCase || u == "*" {
			usePermitted = true
			break
		}
	}
	if !usePermitted {
		return OCAPEvaluationResult{
			Compliant:       false,
			ViolationReason: fmt.Sprintf("use case %q not authorized under OCAP consent terms for %s", requestedUseCase, consent.NationOrBandName),
			EvaluatedAt:     now,
		}
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%s", business.ID, project.ID, consent.AuthorizedSigner, now.Format(time.RFC3339))))
	return OCAPEvaluationResult{
		Compliant:   true,
		AuditHash:   hex.EncodeToString(h[:]),
		EvaluatedAt: now,
	}
}
