package i18n

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestBilingualTranslation_OfficialLanguagesAct(t *testing.T) {
	// English sector
	enSector := TranslateSector(domain.SectorCriticalMinerals, LangEnglish)
	if enSector != "Critical Minerals" {
		t.Errorf("expected Critical Minerals, got %s", enSector)
	}

	// French sector (Bill C-13)
	frSector := TranslateSector(domain.SectorCriticalMinerals, LangFrench)
	if frSector != "Minéraux critiques" {
		t.Errorf("expected Minéraux critiques, got %s", frSector)
	}

	// Stage translation
	enStage := TranslateStage(domain.StageConstruction, LangEnglish)
	if enStage != "Construction" {
		t.Errorf("expected Construction, got %s", enStage)
	}

	frStage := TranslateStage(domain.StageConstruction, LangFrench)
	if frStage != "En construction" {
		t.Errorf("expected En construction, got %s", frStage)
	}

	// Header parsing
	if NormalizeLanguage("fr-CA,fr;q=0.9,en;q=0.8") != LangFrench {
		t.Error("expected French language from fr-CA header")
	}
	if NormalizeLanguage("en-US,en;q=0.5") != LangEnglish {
		t.Error("expected English language from en-US header")
	}
}
