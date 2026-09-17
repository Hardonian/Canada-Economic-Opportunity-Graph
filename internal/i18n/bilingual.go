package i18n

import (
	"strings"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// Language represents an official language of Canada.
type Language string

const (
	LangEnglish Language = "en"
	LangFrench  Language = "fr"
)

// SectorTranslation maps Canadian infrastructure sectors to English and French.
var SectorTranslation = map[domain.Sector]map[Language]string{
	domain.SectorCriticalMinerals: {
		LangEnglish: "Critical Minerals",
		LangFrench:  "Minéraux critiques",
	},
	domain.SectorNuclearEnergy: {
		LangEnglish: "Nuclear & Clean Power",
		LangFrench:  "Énergie nucléaire et propre",
	},
	domain.SectorCleanEnergy: {
		LangEnglish: "Clean Energy & Grid",
		LangFrench:  "Énergie propre et réseau électrique",
	},
	domain.SectorAICompute: {
		LangEnglish: "AI Compute & Data Centres",
		LangFrench:  "Calcul IA et centres de données",
	},
	domain.SectorDefenceArctic: {
		LangEnglish: "Defence & Arctic",
		LangFrench:  "Défense et Arctique",
	},
	domain.SectorTransportation: {
		LangEnglish: "Transportation & Ports",
		LangFrench:  "Transports et ports",
	},
	domain.SectorIndustrialMfg: {
		LangEnglish: "Industrial & Manufacturing",
		LangFrench:  "Industrie et fabrication",
	},
	domain.SectorHousingEnabling: {
		LangEnglish: "Housing-Enabling Infrastructure",
		LangFrench:  "Infrastructures habilitantes au logement",
	},
	domain.SectorMiningMetals: {
		LangEnglish: "Mining & Metals",
		LangFrench:  "Mines et métaux",
	},
	domain.SectorEnergyFuels: {
		LangEnglish: "Energy & Fuels",
		LangFrench:  "Énergie et carburants",
	},
	domain.SectorForestryBioeconomy: {
		LangEnglish: "Forestry & Bioeconomy",
		LangFrench:  "Foresterie et bioéconomie",
	},
}

// StageTranslation maps project lifecycle stages.
var StageTranslation = map[domain.LifecycleStage]map[Language]string{
	domain.StageAnnounced: {
		LangEnglish: "Announced",
		LangFrench:  "Annoncé",
	},
	domain.StageEnvironmentalReview: {
		LangEnglish: "Under Environmental Review",
		LangFrench:  "En cours d'évaluation environnementale",
	},
	domain.StagePermitting: {
		LangEnglish: "Permitting",
		LangFrench:  "Délivrance de permis",
	},
	domain.StageFID: {
		LangEnglish: "Final Investment Decision (FID)",
		LangFrench:  "Décision finale d'investissement (DFI)",
	},
	domain.StageConstruction: {
		LangEnglish: "Construction",
		LangFrench:  "En construction",
	},
	domain.StageOperating: {
		LangEnglish: "Operating",
		LangFrench:  "Opérationnel",
	},
	domain.StageDelayed: {
		LangEnglish: "Delayed",
		LangFrench:  "Retardé",
	},
}

// TranslateSector returns the localized sector name conforming to Bill C-13.
func TranslateSector(sector domain.Sector, lang Language) string {
	if m, ok := SectorTranslation[sector]; ok {
		if val, exists := m[lang]; exists {
			return val
		}
		return m[LangEnglish]
	}
	return string(sector)
}

// TranslateStage returns the localized project stage.
func TranslateStage(stage domain.LifecycleStage, lang Language) string {
	if m, ok := StageTranslation[stage]; ok {
		if val, exists := m[lang]; exists {
			return val
		}
		return m[LangEnglish]
	}
	return string(stage)
}

// NormalizeLanguage extracts the preferred language from an Accept-Language HTTP header.
func NormalizeLanguage(acceptHeader string) Language {
	lower := strings.ToLower(acceptHeader)
	if strings.Contains(lower, "fr") {
		return LangFrench
	}
	return LangEnglish
}
