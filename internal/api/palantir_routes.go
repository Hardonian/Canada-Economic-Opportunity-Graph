package api

import (
	"context"
	"encoding/json"
	"net/http"
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

func (s *Server) registerPalantirGradeRoutes() {
	// Pillar 1: Dynamic Sovereign Ontology
	s.mux.HandleFunc("POST /api/v1/ontology/actions/execute", s.handleOntologyExecuteAction)
	s.mux.HandleFunc("GET /api/v1/ontology/branches", s.handleOntologyBranches)
	s.mux.HandleFunc("POST /api/v1/ontology/branches", s.handleOntologyCreateBranch)

	// Pillar 2: Real-Time Data Fabric & Lakehouse
	s.mux.HandleFunc("GET /api/v1/lakehouse/query", s.handleLakehouseQuery)
	s.mux.HandleFunc("POST /api/v1/lakehouse/query", s.handleLakehouseQuery)
	s.mux.HandleFunc("GET /api/v1/lakehouse/quality", s.handleLakehouseQuality)

	// Pillar 3: Knowledge Graph & GQL
	s.mux.HandleFunc("POST /api/v1/graph/gql", s.handleGraphGQL)
	s.mux.HandleFunc("GET /api/v1/graph/gnn-predictions", s.handleGraphGNNPredictions)
	s.mux.HandleFunc("GET /api/v1/graph/communities", s.handleGraphCommunities)

	// Pillar 4: Earth Observation & Sensor Fusion
	s.mux.HandleFunc("GET /api/v1/earthobs/insar/{id}", s.handleEarthObsInSAR)
	s.mux.HandleFunc("GET /api/v1/earthobs/ais/congestion", s.handleEarthObsAISCongestion)
	s.mux.HandleFunc("GET /api/v1/earthobs/discrepancies", s.handleEarthObsDiscrepancies)

	// Pillar 5: Sovereign AIP Multi-Agent Mesh
	s.mux.HandleFunc("POST /api/v1/ai/copilot/run", s.handleAICopilotRun)
	s.mux.HandleFunc("POST /api/v1/ai/negotiate/ppa", s.handleAINegotiatePPA)

	// Pillar 7: Sovereign Security & ABAC
	s.mux.HandleFunc("POST /api/v1/security/authorize", s.handleSecurityAuthorize)
	s.mux.HandleFunc("POST /api/v1/security/dlp/scan", s.handleSecurityDLPScan)

	// Pillar 8: Counter-Intelligence & UBO
	s.mux.HandleFunc("GET /api/v1/counter-intel/ubo/trace", s.handleCounterIntelUBOTrace)
	s.mux.HandleFunc("GET /api/v1/counter-intel/dark-fleet", s.handleCounterIntelDarkFleet)
	s.mux.HandleFunc("POST /api/v1/counter-intel/ica/screen", s.handleCounterIntelICAScreen)

	// Pillar 10: Sovereign Compliance & Merkle Root
	s.mux.HandleFunc("GET /api/v1/compliance/c59/audit/{id}", s.handleComplianceC59Audit)
	s.mux.HandleFunc("GET /api/v1/compliance/c69/clock/{id}", s.handleComplianceC69Clock)
	s.mux.HandleFunc("GET /api/v1/compliance/trc92/scorecard/{id}", s.handleComplianceTRC92Scorecard)
	s.mux.HandleFunc("GET /api/v1/compliance/merkle/root", s.handleComplianceMerkleRoot)
}

// ----------------------------------------------------------------------
// Pillar 1: Ontology Handlers
// ----------------------------------------------------------------------

func (s *Server) handleOntologyExecuteAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID     string                 `json:"action_id"`
		TargetNodeID string                 `json:"target_node_id"`
		ActionType   string                 `json:"action_type"`
		UserID       string                 `json:"user_id"`
		Parameters   map[string]interface{} `json:"parameters"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.ActionID == "" {
		req.ActionID = "act-" + time.Now().Format("20060102150405")
	}
	if req.TargetNodeID == "" {
		req.TargetNodeID = "crawford-nickel"
	}
	if req.ActionType == "" {
		req.ActionType = "UPDATE_PROJECT_STAGE"
	}
	if req.UserID == "" {
		req.UserID = "USER-SOVEREIGN"
	}

	invalidator := ontology.NewCascadeInvalidator()
	engine := ontology.NewWriteBackEngine(invalidator)

	act := &ontology.Action{
		ID:             req.ActionID,
		Type:           ontology.ActionType(req.ActionType),
		TargetNodeID:   req.TargetNodeID,
		TargetNodeType: "PROJECT",
		BranchID:       "main",
		Initiator:      req.UserID,
		Parameters:     req.Parameters,
		ValidTime:      time.Now().UTC(),
		SystemTime:     time.Now().UTC(),
	}

	res, err := engine.ExecuteAction(context.Background(), act)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "execution_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleOntologyBranches(w http.ResponseWriter, r *http.Request) {
	invalidator := ontology.NewCascadeInvalidator()
	engine := ontology.NewWriteBackEngine(invalidator)
	mgr := ontology.NewBranchManager(engine)

	branches := mgr.ListBranches()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"branches": branches,
		"count":    len(branches),
	})
}

func (s *Server) handleOntologyCreateBranch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		BaseBranch  string `json:"base_branch"`
		CreatedBy   string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid_json", "Invalid branch specification")
		return
	}
	if req.BaseBranch == "" {
		req.BaseBranch = "main"
	}
	if req.CreatedBy == "" {
		req.CreatedBy = "USER"
	}

	invalidator := ontology.NewCascadeInvalidator()
	engine := ontology.NewWriteBackEngine(invalidator)
	mgr := ontology.NewBranchManager(engine)

	b, err := mgr.CreateBranch(req.ID, req.Name, req.Description, req.BaseBranch, req.CreatedBy)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "branch_exists", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, b)
}

// ----------------------------------------------------------------------
// Pillar 2: Lakehouse Handlers
// ----------------------------------------------------------------------

func (s *Server) handleLakehouseQuery(w http.ResponseWriter, r *http.Request) {
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
		},
	}

	res, err := engine.Execute(rb, queryReq)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "lakehouse_error", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleLakehouseQuality(w http.ResponseWriter, r *http.Request) {
	qe := lakehouse.NewQuarantineEngine()
	valid, rec := qe.Validate("projects", map[string]interface{}{
		"id":        "crawford-nickel",
		"name":      "Crawford Nickel Sulphide",
		"capex_cad": 3500000000.0,
		"latitude":  48.7,
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"record_valid":      valid,
		"quarantine_record": rec,
		"quarantine_count":  len(qe.GetQuarantinedRecords()),
	})
}

// ----------------------------------------------------------------------
// Pillar 3: Knowledge Graph Handlers
// ----------------------------------------------------------------------

func (s *Server) handleGraphGQL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GQLQuery string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GQLQuery == "" {
		req.GQLQuery = "MATCH (p:Project) WHERE p.sector = 'Critical Minerals' RETURN p"
	}

	ast, err := graphanalytics.ParseGQL(req.GQLQuery)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "gql_parse_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"query": req.GQLQuery,
		"ast":   ast,
	})
}

func (s *Server) handleGraphGNNPredictions(w http.ResponseWriter, r *http.Request) {
	predictor := graphanalytics.NewGNNLinkPredictor()
	projects, _, _ := s.store.ListProjects(r.Context(), database.ProjectFilter{})
	predictions := predictor.PredictPartnerships(projects, 0.5)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"predictions": predictions,
		"count":       len(predictions),
	})
}

func (s *Server) handleGraphCommunities(w http.ResponseWriter, r *http.Request) {
	clusterer := graphanalytics.NewLouvainClusterer()
	nodes := []string{"crawford-nickel", "timmins-substation", "ontario-hydro-one", "noront-eagles-nest", "james-bay-road"}
	edges := map[string][]string{
		"crawford-nickel":     {"timmins-substation"},
		"timmins-substation":  {"crawford-nickel", "ontario-hydro-one"},
		"ontario-hydro-one":   {"timmins-substation"},
		"noront-eagles-nest":  {"james-bay-road"},
		"james-bay-road":      {"noront-eagles-nest"},
	}

	communities := clusterer.DetectCommunities(nodes, edges)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"communities": communities,
		"algorithm":   "Louvain Modular Graph Partitioning",
	})
}

// ----------------------------------------------------------------------
// Pillar 4: Earth Observation Handlers
// ----------------------------------------------------------------------

func (s *Server) handleEarthObsInSAR(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	pipeline := earthobs.NewSARPipeline()
	now := time.Now().UTC()
	prev := now.AddDate(0, 0, -12)

	analysis := pipeline.ComputeInSAR(prev, now, 0.45, -12.4, -18.2)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"project_id": projectID,
		"insar":      analysis,
	})
}

func (s *Server) handleEarthObsAISCongestion(w http.ResponseWriter, r *http.Request) {
	ate := earthobs.NewAISTrackingEngine()
	now := time.Now().UTC()
	vesselsPRR := []earthobs.AISVessel{
		{MMSI: "316001234", VesselName: "PACIFIC SOVEREIGN", VesselType: "BULK_CARRIER", Status: "MOORED", ArrivalDate: now.Add(-36 * time.Hour)},
		{MMSI: "316009876", VesselName: "MAPLE ARCTIC", VesselType: "CONTAINER", Status: "ANCHORED", ArrivalDate: now.Add(-48 * time.Hour)},
	}

	telemetry := ate.EvaluatePort("PRR", "Port of Prince Rupert", vesselsPRR)
	writeJSON(w, http.StatusOK, telemetry)
}

func (s *Server) handleEarthObsDiscrepancies(w http.ResponseWriter, r *http.Request) {
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

	report := engine.EvaluateDiscrepancy("crawford-nickel", 50.0, mach, sar)
	writeJSON(w, http.StatusOK, report)
}

// ----------------------------------------------------------------------
// Pillar 5: Sovereign AIP Multi-Agent Handlers
// ----------------------------------------------------------------------

func (s *Server) handleAICopilotRun(w http.ResponseWriter, r *http.Request) {
	guardrails := ai.NewAIPGuardrails()
	orchestrator := ai.NewAgentMeshOrchestrator(guardrails)

	project, err := s.store.GetProject(r.Context(), "crawford-nickel")
	if err != nil || project == nil {
		project = &domain.Project{
			ID:          "crawford-nickel",
			Name:        "Crawford Nickel Sulphide",
			Sector:      domain.SectorCriticalMinerals,
			Subsector:   "Nickel-Cobalt",
			CapexCAD:    3500000000,
			Province:    "ON",
		}
	}

	report := orchestrator.EvaluateProject(project)
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleAINegotiatePPA(w http.ResponseWriter, r *http.Request) {
	negotiator := ai.NewMultiAgentNegotiator()
	proposal := negotiator.OptimizeContractTerms(3500000000, 450.0)
	writeJSON(w, http.StatusOK, proposal)
}

// ----------------------------------------------------------------------
// Pillar 7: Sovereign Security & ABAC Handlers
// ----------------------------------------------------------------------

func (s *Server) handleSecurityAuthorize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubjectClearance     string `json:"subject_clearance"`
		ObjectClassification string `json:"object_classification"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.SubjectClearance = "SECRET"
		req.ObjectClassification = "PROTECTED_B"
	}

	sub := security.SecuritySubject{
		UserID:         "OFFICER-001",
		Clearance:      security.ClassificationLevel(req.SubjectClearance),
		Citizenship:    "CA",
		Organization:   "Treasury Board Secretariat",
		GrantedCaveats: []security.Caveat{security.CaveatCanadianEyesOnly},
	}
	obj := security.SecurityObject{
		ResourceID:     "crawford-nickel-geotechnical",
		ResourceType:   "GEOTECHNICAL_REPORT",
		Classification: security.ClassificationLevel(req.ObjectClassification),
		Caveats:        []security.Caveat{security.CaveatCanadianEyesOnly},
	}

	evaluator := security.NewABACEvaluator()
	decision := evaluator.Authorize(sub, obj, "READ")
	writeJSON(w, http.StatusOK, decision)
}

func (s *Server) handleSecurityDLPScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		req.Content = "Proponent contact SIN 123-456-789 deployed at NORAD_SITE_NORTH_BAY."
	}

	dlp := security.NewDLPScanner()
	cleanContent, findings := dlp.ScrubText(req.Content)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sanitized_content": cleanContent,
		"findings_count":    len(findings),
		"findings":          findings,
	})
}

// ----------------------------------------------------------------------
// Pillar 8: Counter-Intelligence & UBO Handlers
// ----------------------------------------------------------------------

func (s *Server) handleCounterIntelUBOTrace(w http.ResponseWriter, r *http.Request) {
	unraveler := ubo.NewUBOUnraveler()
	holdings := []ubo.HoldingNode{
		{EntityID: "E1", Name: "Crawford Mining Ltd", Jurisdiction: "CA", DirectShare: 1.0, ParentID: "E2"},
		{EntityID: "E2", Name: "Northern Horizon Holdings", Jurisdiction: "KY", DirectShare: 0.70, ParentID: "E3"},
		{EntityID: "E3", Name: "State Strategic Assets Inc", Jurisdiction: "CN", DirectShare: 1.0},
	}

	owners := unraveler.UnravelChain("E1", holdings)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"root_entity_id": "E1",
		"owners":         owners,
		"count":          len(owners),
	})
}

func (s *Server) handleCounterIntelDarkFleet(w http.ResponseWriter, r *http.Request) {
	detector := corridor.NewDarkFleetDetector()
	alerts := detector.EvaluateVessel("316001234", "MAPLE SHADOW", "PA", 6.5, 120.0, 3)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"dark_fleet_anomalies": alerts,
		"count":                len(alerts),
	})
}

func (s *Server) handleCounterIntelICAScreen(w http.ResponseWriter, r *http.Request) {
	evaluator := ubo.NewEvaluator()
	project, err := s.store.GetProject(r.Context(), "crawford-nickel")
	if err != nil || project == nil {
		project = &domain.Project{
			ID:          "crawford-nickel",
			Name:        "Crawford Nickel Sulphide",
			Sector:      domain.SectorCriticalMinerals,
			Subsector:   "Nickel",
			ProponentID: "canada-nickel-corp",
		}
	}

	owners := []ubo.BeneficialOwner{
		{EntityID: "O1", Name: "Public Institutional Float", Jurisdiction: "CA", OwnershipPercent: 75.0, IsSOE: false, IsFTA: true},
		{EntityID: "O2", Name: "Overseas Strategic Minerals", Jurisdiction: "CN", OwnershipPercent: 25.0, IsSOE: true, IsFTA: false},
	}

	result := evaluator.ScreenProject(project, owners)
	writeJSON(w, http.StatusOK, result)
}

// ----------------------------------------------------------------------
// Pillar 10: Compliance & Merkle Root Handlers
// ----------------------------------------------------------------------

func (s *Server) handleComplianceC59Audit(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	auditor := cegs.NewC59Auditor()
	report := auditor.AuditProject(projectID, 3500000000.0, 100000.0, 12500.0, true, 30.0)

	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleComplianceC69Clock(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	tracker := cegs.NewC69ClockAuditor()
	clock := tracker.AuditClock(projectID, cegs.PhaseImpactStatement, 145, false)

	writeJSON(w, http.StatusOK, clock)
}

func (s *Server) handleComplianceTRC92Scorecard(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	evaluator := cegs.NewTRC92Evaluator()
	scorecard := evaluator.Evaluate(projectID, 20.0, 8.5, 100.0)

	writeJSON(w, http.StatusOK, scorecard)
}

func (s *Server) handleComplianceMerkleRoot(w http.ResponseWriter, r *http.Request) {
	root := merkle.BuildRoot(r.Context(), s.store)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"merkle_root": root,
		"algorithm":   "SHA-256 Binary Tree (RFC 6962)",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}
