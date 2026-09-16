package gridphysics

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// SystemOperator identifies provincial balancing authorities and ISOs.
type SystemOperator string

const (
	OperatorIESO        SystemOperator = "IESO (Ontario)"
	OperatorAESO        SystemOperator = "AESO (Alberta)"
	OperatorHydroQuebec SystemOperator = "Hydro-Québec (TransÉnergie)"
	OperatorBCHydro     SystemOperator = "BC Hydro"
	OperatorManitoba    SystemOperator = "Manitoba Hydro"
	OperatorSaskPower   SystemOperator = "SaskPower"
	OperatorAtlantic    SystemOperator = "Atlantic Utilities"
	OperatorNorthern    SystemOperator = "Northern Off-Grid"
)

// InterconnectAssessment provides engineering and queue feasibility for major grid ties.
type InterconnectAssessment struct {
	ProjectID                 string         `json:"project_id"`
	ProjectName               string         `json:"project_name"`
	Province                  string         `json:"province"`
	Operator                  SystemOperator `json:"system_operator"`
	EstimatedLoadOrGenMW      float64        `json:"estimated_load_or_gen_mw"`
	InterconnectVoltageKV     int            `json:"interconnect_voltage_kv"` // 115, 230, 500 kV
	QueueEstimatedMonths      int            `json:"queue_estimated_months"`
	SubstationHeadroomMW      float64        `json:"substation_headroom_mw"`
	DedicatedSubstationNeeded bool           `json:"dedicated_substation_needed"`
	ReinforcementCostCAD      int64          `json:"reinforcement_cost_cad"`
	GridFeasibilityScore      float64        `json:"grid_feasibility_score"` // 0.0 - 100.0
	CleanPowerPurityPct       float64        `json:"clean_power_purity_pct"` // % Hydro/Nuclear on local balancing zone
	InterconnectNotes         []string       `json:"interconnect_notes"`
	AuditHash                 string         `json:"audit_hash"`
	CalculatedAt              time.Time      `json:"calculated_at"`
}

// Engine evaluates grid feasibility and queue delays across Canadian provinces.
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// AssessProject performs deterministic grid interconnection analysis.
func (e *Engine) AssessProject(project *domain.Project) *InterconnectAssessment {
	prov := strings.ToUpper(strings.TrimSpace(project.Province))
	var op SystemOperator
	var cleanPurity float64
	var baseQueueMonths int

	switch prov {
	case "ON", "ONTARIO":
		op = OperatorIESO
		cleanPurity = 92.0 // Nuclear + Hydro + Wind
		baseQueueMonths = 28
	case "AB", "ALBERTA":
		op = OperatorAESO
		cleanPurity = 35.0 // Gas heavy with expanding solar/wind
		baseQueueMonths = 20
	case "QC", "QUEBEC":
		op = OperatorHydroQuebec
		cleanPurity = 99.5 // Pure Hydro
		baseQueueMonths = 36 // Heavy queue backlog for large industrial loads
	case "BC", "BRITISH COLUMBIA":
		op = OperatorBCHydro
		cleanPurity = 98.0 // Hydro
		baseQueueMonths = 24
	case "MB", "MANITOBA":
		op = OperatorManitoba
		cleanPurity = 97.0
		baseQueueMonths = 18
	case "SK", "SASKATCHEWAN":
		op = OperatorSaskPower
		cleanPurity = 40.0
		baseQueueMonths = 16
	case "NU", "NT", "YT", "NUNAVUT", "NORTHWEST TERRITORIES", "YUKON":
		op = OperatorNorthern
		cleanPurity = 20.0 // Diesel microgrids unless dedicated hydro/SMR
		baseQueueMonths = 12
	default:
		op = OperatorAtlantic
		cleanPurity = 65.0
		baseQueueMonths = 22
	}

	// Estimate MW load/gen from CAPEX and Sector
	mw := 50.0
	switch project.Sector {
	case domain.SectorNuclearEnergy:
		mw = 300.0 // SMR scale
	case domain.SectorAICompute:
		mw = 150.0 // Hyperscale data centre cluster
	case domain.SectorCleanEnergy:
		mw = 100.0
	case domain.SectorCriticalMinerals:
		mw = 80.0 // Mining & mill processing
	}

	voltage := 230
	if mw > 200.0 {
		voltage = 500
	} else if mw < 60.0 {
		voltage = 115
	}

	substationNeeded := mw >= 50.0
	var reinforcementCAD int64
	if substationNeeded {
		reinforcementCAD = int64(math.Round(15_000_000 + (mw * 350_000))) // $15M base + $350k/MW
	}

	headroom := 200.0 - mw
	if headroom < 0 {
		headroom = 0
	}

	score := 75.0
	if op == OperatorHydroQuebec && mw > 100.0 {
		score -= 15.0 // Quebec Bill 2 / Hydro-Québec capacity rationing
	}
	if substationNeeded {
		score -= 10.0
	}
	if cleanPurity > 90.0 {
		score += 15.0
	}
	if score > 100.0 {
		score = 100.0
	} else if score < 0.0 {
		score = 0.0
	}

	notes := []string{
		fmt.Sprintf("Governing Balancing Authority: %s", op),
		fmt.Sprintf("Estimated Interconnect Voltage: %d kV at %.0f MW", voltage, mw),
		fmt.Sprintf("Grid Purity Profile: %.1f%% non-emitting generation", cleanPurity),
	}
	if substationNeeded {
		notes = append(notes, fmt.Sprintf("Dedicated on-site substation step-up/down required (~$%.1fM CAD).", float64(reinforcementCAD)/1e6))
	}

	assessment := &InterconnectAssessment{
		ProjectID:                 project.ID,
		ProjectName:               project.Name,
		Province:                  project.Province,
		Operator:                  op,
		EstimatedLoadOrGenMW:      mw,
		InterconnectVoltageKV:     voltage,
		QueueEstimatedMonths:      baseQueueMonths,
		SubstationHeadroomMW:      headroom,
		DedicatedSubstationNeeded: substationNeeded,
		ReinforcementCostCAD:      reinforcementCAD,
		GridFeasibilityScore:      score,
		CleanPowerPurityPct:       cleanPurity,
		InterconnectNotes:         notes,
		CalculatedAt:              time.Now().UTC(),
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%.1f|%d|%.1f", assessment.ProjectID, assessment.Operator, assessment.EstimatedLoadOrGenMW, assessment.ReinforcementCostCAD, assessment.GridFeasibilityScore)))
	assessment.AuditHash = hex.EncodeToString(h[:])

	return assessment
}
