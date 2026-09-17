package security

import (
	"fmt"
	"time"
)

// ABACEvaluator evaluates Attribute-Based Access Control policies.
type ABACEvaluator struct{}

// NewABACEvaluator creates a new ABAC authorization engine.
func NewABACEvaluator() *ABACEvaluator {
	return &ABACEvaluator{}
}

// Authorize checks whether a security subject is permitted to access a protected object.
func (e *ABACEvaluator) Authorize(sub SecuritySubject, obj SecurityObject, action string) ABACDecision {
	now := time.Now().UTC()

	// 1. Dominance check: subject clearance must dominate object classification
	subRank := ClearanceRank[sub.Clearance]
	objRank := ClearanceRank[obj.Classification]

	if subRank < objRank {
		return ABACDecision{
			SubjectID:    sub.UserID,
			ResourceID:   obj.ResourceID,
			Permitted:    false,
			DeniedReason: fmt.Sprintf("insufficient clearance: subject holds %s (rank %d), resource requires %s (rank %d)", sub.Clearance, subRank, obj.Classification, objRank),
			EvaluatedAt:  now,
		}
	}

	// 2. Compartment check
	if obj.RequiredCompartment != "" {
		hasCompartment := false
		for _, c := range sub.Compartments {
			if c == obj.RequiredCompartment {
				hasCompartment = true
				break
			}
		}
		if !hasCompartment {
			return ABACDecision{
				SubjectID:    sub.UserID,
				ResourceID:   obj.ResourceID,
				Permitted:    false,
				DeniedReason: fmt.Sprintf("missing required compartment %s", obj.RequiredCompartment),
				EvaluatedAt:  now,
			}
		}
	}

	// 3. Caveat checks
	for _, cav := range obj.Caveats {
		switch cav {
		case CaveatCanadianEyesOnly:
			if sub.Citizenship != "CA" {
				return ABACDecision{
					SubjectID:    sub.UserID,
					ResourceID:   obj.ResourceID,
					Permitted:    false,
					DeniedReason: fmt.Sprintf("resource is CANADIAN_EYES_ONLY; subject citizenship is %s", sub.Citizenship),
					EvaluatedAt:  now,
				}
			}
		case CaveatCommercialConf:
			if sub.Organization != obj.OwnerOrg && sub.Clearance != ClassTopSecret {
				return ABACDecision{
					SubjectID:    sub.UserID,
					ResourceID:   obj.ResourceID,
					Permitted:    false,
					DeniedReason: fmt.Sprintf("resource is COMMERCIAL_IN_CONFIDENCE to %s; subject belongs to %s", obj.OwnerOrg, sub.Organization),
					EvaluatedAt:  now,
				}
			}
		}
	}

	return ABACDecision{
		SubjectID:   sub.UserID,
		ResourceID:  obj.ResourceID,
		Permitted:   true,
		EvaluatedAt: now,
	}
}
