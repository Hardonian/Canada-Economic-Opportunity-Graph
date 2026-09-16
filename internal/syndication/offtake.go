package syndication

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// OfftakeCommodity classifies the contracted output.
type OfftakeCommodity string

const (
	OfftakeCleanPowerPPA     OfftakeCommodity = "CLEAN_POWER_PPA"
	OfftakeBatteryNickel     OfftakeCommodity = "BATTERY_GRADE_NICKEL_SULPHATE"
	OfftakeLithiumHydroxide  OfftakeCommodity = "LITHIUM_HYDROXIDE_MONOHYDRATE"
	OfftakeGreenHydrogen     OfftakeCommodity = "GREEN_HYDROGEN_AMMONIA"
	OfftakeLNGLiquefaction   OfftakeCommodity = "LNG_EXPORT_VOLUME"
	OfftakeAirportLeasing    OfftakeCommodity = "AEROSPACE_LOGISTICS_CONCESSION"
)

// PricingStructure categorizes contractual pricing mechanisms.
type PricingStructure string

const (
	PricingFixedFee         PricingStructure = "FIXED_ESCALATING_PRICE"
	PricingMarketIndexedCap PricingStructure = "MARKET_INDEXED_WITH_FLOOR"
	PricingTollingTariff    PricingStructure = "TOLLING_TARIFF_STRUCTURE"
	PricingCostPlusMargin   PricingStructure = "COST_PLUS_TARGET_MARGIN"
)

// OfftakeAgreement models a commercial off-take or power purchase agreement.
type OfftakeAgreement struct {
	ID                     string           `json:"id"`
	ProjectID              string           `json:"project_id"`
	BuyerEntityID          string           `json:"buyer_entity_id"`
	BuyerName              string           `json:"buyer_name"`
	BuyerCreditRating      string           `json:"buyer_credit_rating"` // e.g. "AA-", "A+", "BBB+"
	Commodity              OfftakeCommodity `json:"commodity"`
	VolumeAnnualMetric     string           `json:"volume_annual_metric"` // e.g., "500 GWh/year", "25,000 tonnes/year"
	TermYears              int              `json:"term_years"`
	PricingStructure       PricingStructure `json:"pricing_structure"`
	TakeOrPayObligation    bool             `json:"take_or_pay_obligation"`
	AnnualContractValueCAD int64            `json:"annual_contract_value_cad"`
	ContractExecutionDate  time.Time        `json:"contract_execution_date"`
	AuditHash              string           `json:"audit_hash"`
}

// CanonicalOfftakeAgreements returns representative bilateral offtake precedents.
func CanonicalOfftakeAgreements() []OfftakeAgreement {
	now := time.Now().UTC()
	return []OfftakeAgreement{
		{
			ID:                     "offtake-ppa-amazon-bwrx300",
			ProjectID:              "darlington-small-modular-reactor-deployment-project",
			BuyerEntityID:          "entity-amazon-web-services",
			BuyerName:              "Amazon Web Services (AWS) Data Centers",
			BuyerCreditRating:      "AA",
			Commodity:              OfftakeCleanPowerPPA,
			VolumeAnnualMetric:     "2,400 GWh/year",
			TermYears:              20,
			PricingStructure:       PricingFixedFee,
			TakeOrPayObligation:    true,
			AnnualContractValueCAD: 192_000_000,
			ContractExecutionDate:  now.Add(-60 * 24 * time.Hour),
		},
		{
			ID:                     "offtake-nickel-vw-crawford",
			ProjectID:              "crawford-nickel-project",
			BuyerEntityID:          "entity-powerco-volkswagen",
			BuyerName:              "PowerCo SE (Volkswagen Battery Subsidiary)",
			BuyerCreditRating:      "A-",
			Commodity:              OfftakeBatteryNickel,
			VolumeAnnualMetric:     "30,000 tonnes/year",
			TermYears:              15,
			PricingStructure:       PricingMarketIndexedCap,
			TakeOrPayObligation:    true,
			AnnualContractValueCAD: 450_000_000,
			ContractExecutionDate:  now.Add(-120 * 24 * time.Hour),
		},
		{
			ID:                     "offtake-airport-air-canada-toronto",
			ProjectID:              "canadian-international-airports-global-investor-leasing-infrastructure-hubs",
			BuyerEntityID:          "entity-air-canada-cargo",
			BuyerName:              "Air Canada Cargo & Global Logistics Consortia",
			BuyerCreditRating:      "BBB-",
			Commodity:              OfftakeAirportLeasing,
			VolumeAnnualMetric:     "500,000 m² Apron & Bonded Warehouse Leasing",
			TermYears:              35,
			PricingStructure:       PricingTollingTariff,
			TakeOrPayObligation:    true,
			AnnualContractValueCAD: 320_000_000,
			ContractExecutionDate:  now.Add(-30 * 24 * time.Hour),
		},
	}
}

// GetProjectOfftakes retrieves registered offtake agreements for an asset.
func GetProjectOfftakes(projectID string) []OfftakeAgreement {
	agreements := CanonicalOfftakeAgreements()
	matched := make([]OfftakeAgreement, 0)
	for _, a := range agreements {
		if a.ProjectID == projectID {
			auditData := fmt.Sprintf("%s|%s|%s|%s|%d|%d", a.ID, a.ProjectID, a.BuyerName, a.Commodity, a.TermYears, a.AnnualContractValueCAD)
			h := sha256.Sum256([]byte(auditData))
			a.AuditHash = hex.EncodeToString(h[:])
			matched = append(matched, a)
		}
	}
	return matched
}
