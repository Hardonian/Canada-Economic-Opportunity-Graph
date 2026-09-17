package security

import (
	"time"
)

// ClassificationLevel defines Canadian Government security marking hierarchy.
type ClassificationLevel string

const (
	ClassUnclassified ClassificationLevel = "UNCLASSIFIED"
	ClassProtectedA   ClassificationLevel = "PROTECTED_A"
	ClassProtectedB   ClassificationLevel = "PROTECTED_B"
	ClassProtectedC   ClassificationLevel = "PROTECTED_C"
	ClassConfidential ClassificationLevel = "CONFIDENTIAL"
	ClassSecret       ClassificationLevel = "SECRET"
	ClassTopSecret    ClassificationLevel = "TOP_SECRET"
)

// ClearanceRank maps classification levels to integer dominance levels.
var ClearanceRank = map[ClassificationLevel]int{
	ClassUnclassified: 0,
	ClassProtectedA:   1,
	ClassProtectedB:   2,
	ClassProtectedC:   3,
	ClassConfidential: 4,
	ClassSecret:       5,
	ClassTopSecret:    6,
}

// Caveat defines distribution caveats (e.g. Canadian Eyes Only).
type Caveat string

const (
	CaveatCanadianEyesOnly Caveat = "CANADIAN_EYES_ONLY"
	CaveatFiveEyes         Caveat = "FIVE_EYES"
	CaveatCommercialConf   Caveat = "COMMERCIAL_IN_CONFIDENCE"
)

// SecuritySubject represents the requesting entity/user in an ABAC decision.
type SecuritySubject struct {
	UserID        string              `json:"user_id"`
	Clearance     ClassificationLevel `json:"clearance"`
	Citizenship   string              `json:"citizenship"` // "CA", "US", "GB", etc.
	Organization  string              `json:"organization"`
	Compartments  []string            `json:"compartments"` // e.g. "NUCLEAR", "CRITICAL_MINERALS"
	GrantedCaveats []Caveat           `json:"granted_caveats"`
}

// SecurityObject represents a protected graph node, edge, or field.
type SecurityObject struct {
	ResourceID          string              `json:"resource_id"`
	ResourceType        string              `json:"resource_type"`
	Classification      ClassificationLevel `json:"classification"`
	RequiredCompartment string              `json:"required_compartment,omitempty"`
	Caveats             []Caveat            `json:"caveats,omitempty"`
	OwnerOrg            string              `json:"owner_org,omitempty"`
}

// ABACDecision holds the verdict of an authorization check.
type ABACDecision struct {
	SubjectID    string    `json:"subject_id"`
	ResourceID   string    `json:"resource_id"`
	Permitted    bool      `json:"permitted"`
	DeniedReason string    `json:"denied_reason,omitempty"`
	EvaluatedAt  time.Time `json:"evaluated_at"`
}
