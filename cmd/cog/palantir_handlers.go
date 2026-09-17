package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ai"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/cegs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/corridor"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/earthobs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/graphanalytics"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/lakehouse"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/merkle"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ontology"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ubo"
)

// ----------------------------------------------------------------------
// Pillar 1: Dynamic Sovereign Ontology
// ----------------------------------------------------------------------

func handleOntology(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog ontology <branch|action> [options]")
		return
	}

	invalidator := ontology.NewCascadeInvalidator()
	engine := ontology.NewWriteBackEngine(invalidator)
	mgr := ontology.NewBranchManager(engine)

	switch args[0] {
	case "branch":
		if len(args) > 1 && args[1] == "create" {
			if len(args) < 4 {
				fmt.Println("Usage: cog ontology branch create <id> <name>")
				return
			}
			b, err := mgr.CreateBranch(args[2], args[3], "Workspace branch", "main", "USER")
			if err != nil {
				fmt.Printf("Branch creation error: %v\n", err)
				return
			}
			fmt.Printf("Branch created: %s (%s) [Base: %s]\n", b.ID, b.Name, b.BaseBranch)
			return
		}
		branches := mgr.ListBranches()
		fmt.Println("=== Sovereign Ontology Branches ===")
		for _, b := range branches {
			fmt.Printf("• [%s] %s — Author: %s, Created: %s\n", b.ID, b.Name, b.CreatedBy, b.CreatedAt.Format(time.RFC3339))
		}

	case "action":
		target := "crawford-nickel"
		actionType := "UPDATE_PROJECT_STAGE"
		if len(args) > 1 {
			target = args[1]
		}
		if len(args) > 2 {
			actionType = args[2]
		}

		act := &ontology.Action{
			ID:             fmt.Sprintf("ACT-%d", time.Now().Unix()),
			Type:           ontology.ActionType(actionType),
			TargetNodeID:   target,
			TargetNodeType: "PROJECT",
			BranchID:       "main",
			Initiator:      "CLI-OPERATOR",
			Parameters:     map[string]interface{}{"status": "APPROVED", "stage": "CONSTRUCTION"},
			ValidTime:      time.Now().UTC(),
			SystemTime:     time.Now().UTC(),
		}

		res, err := engine.ExecuteAction(context.Background(), act)
		if err != nil {
			fmt.Printf("Action execution error: %v\n", err)
			return
		}
		fmt.Printf("Executed Action %s [%s] on Node %s -> Status: %s\n", res.ID, res.Type, res.TargetNodeID, res.Status)
	default:
		fmt.Println("Usage: cog ontology <branch|action>")
	}
}

// ----------------------------------------------------------------------
// Pillar 2: Real-Time Data Fabric & Lakehouse
// ----------------------------------------------------------------------

func handleLakehouse(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog lakehouse <query|quality>")
		return
	}

	switch args[0] {
	case "query":
		schema := lakehouse.TableSchema{
			SchemaID: 1,
			Fields: []lakehouse.FieldSchema{
				{ID: 1, Name: "project_id", Type: lakehouse.TypeString, Required: true},
				{ID: 2, Name: "sector", Type: lakehouse.TypeString, Required: true},
				{ID: 3, Name: "capex_cad", Type: lakehouse.TypeFloat64, Required: true},
			},
		}
		rb := lakehouse.NewRecordBatch(schema)
		_ = rb.AppendRow(map[string]interface{}{"project_id": "crawford-nickel", "sector": "Critical Minerals", "capex_cad": 3500.0})
		_ = rb.AppendRow(map[string]interface{}{"project_id": "eagles-nest", "sector": "Critical Minerals", "capex_cad": 2200.0})
		_ = rb.AppendRow(map[string]interface{}{"project_id": "bécancour-hub", "sector": "Clean Tech", "capex_cad": 1400.0})

		engine := lakehouse.NewVectorQueryEngine()
		queryReq := lakehouse.QueryRequest{
			GroupBy: []string{"sector"},
			Aggregations: []lakehouse.AggregationSpec{
				{Column: "capex_cad", Type: lakehouse.AggSum, Alias: "total_capex_cad"},
				{Column: "capex_cad", Type: lakehouse.AggCount, Alias: "project_count"},
			},
		}
		res, err := engine.Execute(rb, queryReq)
		if err != nil {
			fmt.Printf("Query error: %v\n", err)
			return
		}

		out, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(out))

	case "quality":
		qe := lakehouse.NewQuarantineEngine()
		valid, rec := qe.Validate("projects", map[string]interface{}{
			"id":        "crawford-nickel",
			"name":      "Crawford Nickel Sulphide",
			"capex_cad": 3500000000.0,
			"latitude":  48.7,
		})
		fmt.Printf("Quality Firewall Invariant Evaluation: Valid=%v, Quarantined Count=%d\n", valid, len(qe.GetQuarantinedRecords()))
		if rec != nil {
			fmt.Printf("Violations: %v\n", rec.Violations)
		}
	}
}

// ----------------------------------------------------------------------
// Pillar 3: Knowledge Graph & GQL
// ----------------------------------------------------------------------

func handleGQL(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog gql <query|gnn|communities>")
		return
	}

	switch args[0] {
	case "query":
		query := "MATCH (p:Project) WHERE p.sector = 'Critical Minerals' RETURN p"
		if len(args) > 1 {
			query = strings.Join(args[1:], " ")
		}
		ast, err := graphanalytics.ParseGQL(query)
		if err != nil {
			fmt.Printf("GQL Parse Error: %v\n", err)
			return
		}
		fmt.Printf("ISO/IEC 39075 GQL Query Parsed Successfully: SourceType=%s, RelType=%s, TargetType=%s, Filter=%s %s '%s'\n",
			ast.SourceType, ast.RelType, ast.TargetType, ast.FilterField, ast.FilterOp, ast.FilterVal)

	case "gnn":
		store := getSeededStore()
		predictor := graphanalytics.NewGNNLinkPredictor()
		projects, _, _ := store.ListProjects(context.Background(), database.ProjectFilter{Limit: 20})
		preds := predictor.PredictPartnerships(projects, 0.5)

		fmt.Println("=== GNN Relational Link Predictions (Latent Synergies) ===")
		for i, p := range preds {
			fmt.Printf("%d. %s <==[%s (%.1f%%, Synergy: %.1f)]==> %s\n   Rationale: %s\n",
				i+1, p.SourceName, p.PredictedRelType, p.Probability*100, p.SynergyScore, p.TargetName, p.Rationale)
		}

	case "communities":
		clusterer := graphanalytics.NewLouvainClusterer()
		nodes := []string{"crawford-nickel", "timmins-substation", "ontario-hydro-one", "noront-eagles-nest", "james-bay-road"}
		edges := map[string][]string{
			"crawford-nickel":    {"timmins-substation"},
			"timmins-substation": {"crawford-nickel", "ontario-hydro-one"},
			"ontario-hydro-one":  {"timmins-substation"},
			"noront-eagles-nest": {"james-bay-road"},
			"james-bay-road":     {"noront-eagles-nest"},
		}
		clusters := clusterer.DetectCommunities(nodes, edges)
		fmt.Println("=== Louvain Modularity Industrial Clusters ===")
		for _, c := range clusters {
			fmt.Printf("Cluster #%d (%s): Members=%v (Density: %.2f)\n", c.ClusterID, c.Name, c.NodeIDs, c.InternalDensity)
		}
	}
}

// ----------------------------------------------------------------------
// Pillar 4: Earth Observation & Sensor Fusion
// ----------------------------------------------------------------------

func handleEarthObsCLI(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog earthobs <insar|ais|audit> [id]")
		return
	}

	switch args[0] {
	case "insar":
		id := "crawford-nickel"
		if len(args) > 1 {
			id = args[1]
		}
		now := time.Now().UTC()
		pipeline := earthobs.NewSARPipeline()
		analysis := pipeline.ComputeInSAR(now.AddDate(0, 0, -12), now, 0.35, -11.5, -17.8)
		fmt.Printf("=== InSAR Sentinel-1 Telemetry for %s ===\n", id)
		fmt.Printf("Coherence:          %.2f\n", analysis.InterferometricCoherence)
		fmt.Printf("LOS Displacement:   %.2f mm\n", analysis.LineOfSightDisplacementMM)
		fmt.Printf("Tailings Dam Status: %s\n", analysis.TailingsStabilityRating)

	case "ais":
		ate := earthobs.NewAISTrackingEngine()
		now := time.Now().UTC()
		vesselsPRR := []earthobs.AISVessel{
			{MMSI: "316001234", VesselName: "PACIFIC SOVEREIGN", VesselType: "BULK_CARRIER", Status: "MOORED", ArrivalDate: now.Add(-36 * time.Hour)},
			{MMSI: "316009876", VesselName: "MAPLE ARCTIC", VesselType: "CONTAINER", Status: "ANCHORED", ArrivalDate: now.Add(-48 * time.Hour)},
		}
		t := ate.EvaluatePort("PRR", "Port of Prince Rupert", vesselsPRR)
		fmt.Printf("=== Maritime Port Velocity (%s) ===\n", t.PortName)
		fmt.Printf("Anchored: %d, Moored: %d\n", t.AnchoredVesselCount, t.MooredVesselCount)
		fmt.Printf("Average Dwell: %.1f days, Berth Utilization: %.1f%%\n", t.AverageDwellDays, t.BerthUtilizationPct)
		fmt.Printf("Friction Rating: %s (Score: %.1f)\n", t.FrictionRating, t.SupplyChainFrictionScore)

	case "audit":
		engine := earthobs.NewGroundTruthDiscrepancyEngine()
		mach := &earthobs.MachineryDetectionResult{
			TotalMachineryUnits:  8,
			ExcavatorCount:       4,
			HaulTruckCount:       3,
			CraneCount:           1,
			EstimatedSiteWorkers: 16,
		}
		now := time.Now().UTC()
		sar := earthobs.NewSARPipeline().ComputeInSAR(now.AddDate(0, 0, -12), now, 0.2, -10.0, -15.0)
		rep := engine.EvaluateDiscrepancy("crawford-nickel", 45.0, mach, sar)
		fmt.Println("=== Filing vs Satellite Ground Truth Reconciliation ===")
		fmt.Printf("Project: %s | Claimed Progress: %.1f%% | Observed: %.1f%%\n", rep.ProjectID, rep.ClaimedProgressPct, rep.ObservedProgressPct)
		fmt.Printf("Verdict: %s (Flagged: %v)\n", rep.RiskVerdict, rep.IsDiscrepancyFlagged)
		fmt.Printf("Explanation: %s\n", rep.Explanation)
	}
}

// ----------------------------------------------------------------------
// Pillar 5: Sovereign AIP Multi-Agent Mesh
// ----------------------------------------------------------------------

func handleAICLI(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog ai <copilot|ppa> [options]")
		return
	}

	switch args[0] {
	case "copilot":
		store := getSeededStore()
		project, _ := store.GetProject(context.Background(), "crawford-nickel")
		if project == nil {
			project = &domain.Project{
				ID:        "crawford-nickel",
				Name:      "Crawford Nickel Sulphide",
				Sector:    domain.SectorCriticalMinerals,
				CapexCAD:  3500000000,
				Province:  "ON",
			}
		}

		guardrails := ai.NewAIPGuardrails()
		orchestrator := ai.NewAgentMeshOrchestrator(guardrails)
		report := orchestrator.EvaluateProject(project)

		fmt.Printf("=== Sovereign AIP Multi-Agent Consensus for %s ===\n", project.Name)
		fmt.Printf("Consensus Score: %.1f/100 (Consensus Verdict: %s)\n", report.ConsensusScore, report.ConsensusVerdict)
		for role, e := range report.EvaluationsByRole {
			finding := "Evaluation completed"
			if len(e.KeyFindings) > 0 {
				finding = e.KeyFindings[0]
			}
			fmt.Printf("• [%s] Score: %.1f | Recommendation: %s | Finding: %s\n", role, e.ConfidenceScore, e.Recommendation, finding)
		}

	case "ppa":
		capex := int64(3500000000)
		mw := 450.0
		if len(args) > 1 {
			if v, err := strconv.ParseInt(args[1], 10, 64); err == nil {
				capex = v
			}
		}
		negotiator := ai.NewMultiAgentNegotiator()
		proposal := negotiator.OptimizeContractTerms(capex, mw)
		fmt.Println("=== Bilateral Clean Power PPA & Debt Structuring ===")
		fmt.Printf("PPA Strike Price: $%.2f CAD/MWh\n", proposal.PPAStrikePriceCADPerMWh)
		fmt.Printf("Debt Tenor:       %d years\n", proposal.DebtTenorYears)
		fmt.Printf("Target IRR:       %.1f%%\n", proposal.TargetEquityIRRPct)
		fmt.Printf("Indigenous Yield: $%.1fM CAD/yr\n", proposal.IndigenousDividendYieldCADPerYr/1e6)
		fmt.Printf("Efficiency Score: %.1f/100\n", proposal.ParetoEfficiencyScore)
		fmt.Printf("Summary:          %s\n", proposal.TermsSummary)
	}
}

// ----------------------------------------------------------------------
// Pillar 7: Sovereign Security & ABAC
// ----------------------------------------------------------------------

func handleSecurityCLI(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog security <auth|dlp> [options]")
		return
	}

	switch args[0] {
	case "auth":
		subClearance := "SECRET"
		objClassification := "PROTECTED_B"
		if len(args) > 1 {
			subClearance = args[1]
		}
		if len(args) > 2 {
			objClassification = args[2]
		}

		sub := security.SecuritySubject{
			UserID:         "OFFICER-001",
			Clearance:      security.ClassificationLevel(subClearance),
			Citizenship:    "CA",
			Organization:   "Treasury Board",
			GrantedCaveats: []security.Caveat{security.CaveatCanadianEyesOnly},
		}
		obj := security.SecurityObject{
			ResourceID:     "res-101",
			ResourceType:   "GEOTECHNICAL_ASSESSMENT",
			Classification: security.ClassificationLevel(objClassification),
			Caveats:        []security.Caveat{security.CaveatCanadianEyesOnly},
		}

		evaluator := security.NewABACEvaluator()
		decision := evaluator.Authorize(sub, obj, "READ")
		fmt.Printf("=== ABAC Clearance Authorization Check ===\n")
		fmt.Printf("Subject Clearance: %s | Object Classification: %s\n", sub.Clearance, obj.Classification)
		fmt.Printf("Verdict: Permitted=%v (Reason: %s)\n", decision.Permitted, decision.DeniedReason)

	case "dlp":
		text := "Proponent contact SIN 123-456-789 deployed at NORAD_SITE_ALERT."
		if len(args) > 1 {
			text = strings.Join(args[1:], " ")
		}
		dlp := security.NewDLPScanner()
		clean, findings := dlp.ScrubText(text)
		fmt.Printf("Original:  %s\n", text)
		fmt.Printf("Sanitized: %s\n", clean)
		fmt.Printf("Redacted:  %v\n", findings)
	}
}

// ----------------------------------------------------------------------
// Pillar 8: Counter-Intelligence & UBO
// ----------------------------------------------------------------------

func handleUBOCLI(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog ubo <trace|darkfleet|ica> [options]")
		return
	}

	switch args[0] {
	case "trace":
		unraveler := ubo.NewUBOUnraveler()
		holdings := []ubo.HoldingNode{
			{EntityID: "E1", Name: "Crawford Mining Ltd", Jurisdiction: "CA", DirectShare: 1.0, ParentID: "E2"},
			{EntityID: "E2", Name: "Northern Horizon Holdings", Jurisdiction: "KY", DirectShare: 0.70, ParentID: "E3"},
			{EntityID: "E3", Name: "State Strategic Assets Inc", Jurisdiction: "CN", DirectShare: 1.0},
		}
		owners := unraveler.UnravelChain("E1", holdings)
		fmt.Println("=== Beneficial Ownership Multi-Tier Unraveler ===")
		for _, o := range owners {
			fmt.Printf("• Owner: %s (%s) — Effective Share: %.1f%%, Hops: %d (Offshore Obfuscated: %v)\n",
				o.Name, o.Jurisdiction, o.EffectiveShare*100, o.ChainHops, o.IsOffshoreObfuscated)
		}

	case "darkfleet":
		detector := corridor.NewDarkFleetDetector()
		alerts := detector.EvaluateVessel("316001234", "MAPLE SHADOW", "PA", 6.5, 120.0, 3)
		fmt.Println("=== Dark Fleet Sanctions Evasion Detections ===")
		for _, a := range alerts {
			fmt.Printf("• [%s] %s (%s): %s (Risk Score: %.1f)\n", a.Anomaly, a.VesselName, a.Flag, a.Description, a.RiskScore)
		}

	case "ica":
		evaluator := ubo.NewEvaluator()
		store := getSeededStore()
		project, _ := store.GetProject(context.Background(), "crawford-nickel")
		if project == nil {
			project = &domain.Project{
				ID:          "crawford-nickel",
				Name:        "Crawford Nickel Sulphide",
				Sector:      domain.SectorCriticalMinerals,
				Subsector:   "Nickel",
				ProponentID: "canada-nickel-corp",
			}
		}

		owners := []ubo.BeneficialOwner{
			{EntityID: "O1", Name: "Canadian Public Float", Jurisdiction: "CA", OwnershipPercent: 75.0, IsSOE: false, IsFTA: true},
			{EntityID: "O2", Name: "Overseas SASAC SOE", Jurisdiction: "CN", OwnershipPercent: 25.0, IsSOE: true, IsFTA: false},
		}
		res := evaluator.ScreenProject(project, owners)
		fmt.Printf("=== Investment Canada Act Section 25.3 National Security Review ===\n")
		fmt.Printf("Project: %s | ICA Risk Verdict: %s\n", res.ProjectID, res.ICARisk)
		fmt.Printf("Domestic Control: %.1f%% | SOE Exposure: %.1f%%\n", res.DomesticControlShare*100, res.SOEExposurePercent)
		for _, note := range res.NationalSecurityNotes {
			fmt.Printf("  • %s\n", note)
		}
	}
}

// ----------------------------------------------------------------------
// Pillar 10: Sovereign Compliance & Merkle Provenance
// ----------------------------------------------------------------------

func handleAuditCLI(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog audit <c59|c69|trc92|merkle> [id]")
		return
	}

	id := "crawford-nickel"
	if len(args) > 1 {
		id = args[1]
	}

	switch args[0] {
	case "c59":
		auditor := cegs.NewC59Auditor()
		rep := auditor.AuditProject(id, 3500000000.0, 100000.0, 12500.0, true, 30.0)
		fmt.Println("=== Bill C-59 Clean Economy ITC Labor Compliance Audit ===")
		fmt.Printf("Project: %s | Status: %s\n", rep.ProjectID, rep.AuditStatus)
		fmt.Printf("Apprentice Ratio: %.1f%% (Meets 10%% Mandate: %v)\n", rep.ApprenticeRatioPct, rep.Meets10PctMandate)
		fmt.Printf("Wages Certified: %v | Eligible ITC Rate: %.1f%%\n", rep.PrevailingWagesCertified, rep.EligibleITCRatePct)
		fmt.Printf("Estimated Monetized ITC: $%.1fM CAD\n", rep.EstimatedITCMonetizedCAD/1e6)

	case "c69":
		auditor := cegs.NewC69ClockAuditor()
		clock := auditor.AuditClock(id, cegs.PhaseImpactStatement, 145, false)
		fmt.Println("=== Bill C-69 Impact Assessment Act Statutory Clock ===")
		fmt.Printf("Phase: %s | Elapsed: %d days / %d statutory cap days\n", clock.Phase, clock.ElapsedDays, clock.StatutoryCapDays)
		fmt.Printf("Remaining Days: %d | Statutory Breach: %v\n", clock.RemainingDays, clock.IsStatutoryBreach)
		fmt.Printf("Litigation Risk: %s\n", clock.LitigationRisk)

	case "trc92":
		evaluator := cegs.NewTRC92Evaluator()
		sc := evaluator.Evaluate(id, 20.0, 8.5, 100.0)
		fmt.Println("=== TRC Call to Action #92 Corporate Reconciliation Scorecard ===")
		fmt.Printf("Project: %s | Overall Score: %.1f/100 (Rating: %s)\n", sc.ProjectID, sc.OverallReconciliationScore, sc.TRCRating)
		fmt.Printf("Pillar 1 (FPIC & Equity): %.1f/100\n", sc.Pillar1ConsentAndEquity)
		fmt.Printf("Pillar 2 (Jobs & Procurement): %.1f/100\n", sc.Pillar2JobsAndProcurement)
		fmt.Printf("Pillar 3 (UNDRIP Education): %.1f/100\n", sc.Pillar3EducationAndUNDRIP)

	case "merkle":
		store := getSeededStore()
		root := merkle.BuildRoot(context.Background(), store)
		fmt.Println("=== Merkle Tree Cryptographic Provenance Root ===")
		fmt.Printf("Root Hash:     %s\n", root.RootHash)
		fmt.Printf("Leaf Count:    %d\n", root.LeafCount)
		fmt.Printf("Published At:  %s\n", root.PublishedAt.Format(time.RFC3339))
	}
}
