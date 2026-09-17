package security

import (
	"regexp"
)

// DLPScanner detects and redacts Personally Identifiable Information (PII) and defense secrets.
type DLPScanner struct {
	sinRegex     *regexp.Regexp
	cardRegex    *regexp.Regexp
	defenseRegex *regexp.Regexp
}

// NewDLPScanner creates a real-time data loss prevention scanner.
func NewDLPScanner() *DLPScanner {
	return &DLPScanner{
		sinRegex:     regexp.MustCompile(`\b\d{3}[ -]?\d{3}[ -]?\d{3}\b`),
		cardRegex:    regexp.MustCompile(`\b(?:\d{4}[ -]?){3}\d{4}\b`),
		defenseRegex: regexp.MustCompile(`(?i)\b(NORAD_SITE_[A-Z0-9]+|CFS_ALERT_COORDS_[A-Z0-9]+)\b`),
	}
}

// ScrubText replaces detected sensitive tokens with masked placeholders.
func (dlp *DLPScanner) ScrubText(input string) (string, []string) {
	var findings []string

	out := dlp.sinRegex.ReplaceAllStringFunc(input, func(match string) string {
		findings = append(findings, "CANADIAN_SIN")
		return "[REDACTED-SIN]"
	})

	out = dlp.cardRegex.ReplaceAllStringFunc(out, func(match string) string {
		findings = append(findings, "PAYMENT_CARD")
		return "[REDACTED-CARD]"
	})

	out = dlp.defenseRegex.ReplaceAllStringFunc(out, func(match string) string {
		findings = append(findings, "DEFENSE_SECRET")
		return "[REDACTED-DEFENSE-COORDINATES]"
	})

	return out, findings
}
