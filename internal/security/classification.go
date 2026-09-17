package security

import (
	"fmt"
	"strings"
)

// GenerateSecurityBanner creates the canonical Canadian Government security header.
func GenerateSecurityBanner(level ClassificationLevel, caveats []Caveat) string {
	parts := []string{string(level)}
	for _, c := range caveats {
		parts = append(parts, string(c))
	}
	return fmt.Sprintf("// CLASSIFICATION: %s //", strings.Join(parts, " // "))
}

// Redactor filters structured records based on the subject's clearance level.
type Redactor struct{}

// NewRedactor creates a security redactor.
func NewRedactor() *Redactor {
	return &Redactor{}
}

// RedactRecord redacts fields whose classification exceeds the subject's clearance.
func (r *Redactor) RedactRecord(data map[string]interface{}, sub SecuritySubject, fieldLevels map[string]ClassificationLevel) map[string]interface{} {
	result := make(map[string]interface{})
	subRank := ClearanceRank[sub.Clearance]

	for k, v := range data {
		fieldReq, hasLevel := fieldLevels[k]
		if hasLevel {
			reqRank := ClearanceRank[fieldReq]
			if subRank < reqRank {
				result[k] = fmt.Sprintf("[REDACTED - %s]", fieldReq)
				continue
			}
		}
		result[k] = v
	}
	return result
}
