package security

import (
	"regexp"
)

// DLPScanner detects and redacts Personally Identifiable Information (PII), defense coordinates,
// sacred Indigenous territories, and confidential bidding data.
type DLPScanner struct {
	sinRegex             *regexp.Regexp
	cardRegex            *regexp.Regexp
	defenseRegex         *regexp.Regexp
	transitRegex         *regexp.Regexp
	sacredSiteRegex      *regexp.Regexp
	confidentialBidRegex *regexp.Regexp
}

// NewDLPScanner creates a real-time data loss prevention scanner.
func NewDLPScanner() *DLPScanner {
	return &DLPScanner{
		sinRegex:             regexp.MustCompile(`\b\d{3}[ -]?\d{3}[ -]?\d{3}\b`),
		cardRegex:            regexp.MustCompile(`\b(?:\d{4}[ -]?){3}\d{4}\b`),
		defenseRegex:         regexp.MustCompile(`(?i)\b(NORAD_SITE_[A-Z0-9_]+|CFS_ALERT_COORDS_[A-Z0-9_]+|DND_CLASSIFIED_RADAR_[A-Z0-9_]+)\b`),
		transitRegex:         regexp.MustCompile(`\b(?:TRANSIT|INSTITUTION)[- :]+(?:\d{5}[- ]?\d{3})\b`),
		sacredSiteRegex:      regexp.MustCompile(`(?i)\b(SACRED_SITE_[A-Z0-9_]+|BURIAL_GROUND_[A-Z0-9_]+|RESTRICTED_INDIGENOUS_CEREMONY_[A-Z0-9_]+)\b`),
		confidentialBidRegex: regexp.MustCompile(`(?i)\b(CONFIDENTIAL_BID_PRICE_CAD_[0-9]+|PROPRIETARY_DISCOUNT_MARGIN_[0-9]+)\b`),
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

	out = dlp.transitRegex.ReplaceAllStringFunc(out, func(match string) string {
		findings = append(findings, "BANK_TRANSIT")
		return "[REDACTED-BANK-TRANSIT]"
	})

	out = dlp.sacredSiteRegex.ReplaceAllStringFunc(out, func(match string) string {
		findings = append(findings, "SACRED_INDIGENOUS_SITE")
		return "[REDACTED-SACRED-SITE-OCAP]"
	})

	out = dlp.confidentialBidRegex.ReplaceAllStringFunc(out, func(match string) string {
		findings = append(findings, "CONFIDENTIAL_BID_PRICING")
		return "[REDACTED-CONFIDENTIAL-BID]"
	})

	return out, findings
}
