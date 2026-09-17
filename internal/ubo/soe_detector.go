package ubo

import (
	"strings"
)

var knownSOEKeywords = []struct {
	Keyword string
	State   string
	Risk    float64
}{
	{"SASAC", "CHINA", 95.0},
	{"NORINCO", "CHINA", 100.0},
	{"SINOCHEM", "CHINA", 90.0},
	{"CNNC", "CHINA", 95.0},
	{"ZHAOJIN", "CHINA", 85.0},
	{"ROSATOM", "RUSSIA", 100.0},
	{"GAZPROM", "RUSSIA", 100.0},
	{"ROSTEC", "RUSSIA", 100.0},
	{"PDVSA", "VENEZUELA", 90.0},
	{"STATE-OWNED", "FOREIGN_STATE", 80.0},
	{"SOVEREIGN WEALTH", "FOREIGN_STATE", 70.0},
}

// SOEDetector identifies foreign state-owned and military-civil fusion corporations.
type SOEDetector struct{}

// NewSOEDetector creates an SOE detection engine.
func NewSOEDetector() *SOEDetector {
	return &SOEDetector{}
}

// DetectSOE analyzes an entity name and jurisdiction against intelligence watchlists.
func (sd *SOEDetector) DetectSOE(name, jurisdiction string) (bool, string, float64) {
	upperName := strings.ToUpper(name)
	upperJuris := strings.ToUpper(jurisdiction)

	for _, k := range knownSOEKeywords {
		if strings.Contains(upperName, k.Keyword) {
			return true, k.State, k.Risk
		}
	}

	// Geopolitical non-FTA auto-elevate
	if upperJuris == "CN" || upperJuris == "CHINA" || upperJuris == "RU" || upperJuris == "RUSSIA" || upperJuris == "IR" || upperJuris == "IRAN" {
		return true, upperJuris, 85.0
	}

	return false, "", 0.0
}
