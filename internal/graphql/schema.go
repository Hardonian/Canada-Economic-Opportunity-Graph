// Package graphql provides a zero-dependency GraphQL-over-HTTP handler for the
// CanadaOpportunityGraph API. It implements a minimal hand-rolled executor that
// covers introspection, simple selections, and typed arguments — sufficient for
// the project/organization/event/signal/aiSovereignty/reconciliation queries
// required by the enterprise API surface.
//
// No external GraphQL libraries are used; the implementation parses the query
// string directly and routes to the appropriate resolver, keeping the module
// dependency-free beyond github.com/google/uuid.
package graphql

// Schema is the normative CEGS GraphQL SDL. It is exposed at GET /api/v1/graphql
// and is returned verbatim by the __schema introspection query.
const Schema = `
"""
CanadaOpportunityGraph Enterprise GraphQL API — CEGS 1.0
"""
schema {
  query:        Query
  subscription: Subscription
}

"""
Root query type.
"""
type Query {
  """List projects with optional filters."""
  projects(
    sector:   String
    province: String
    stage:    String
    limit:    Int
    offset:   Int
  ): [Project!]!

  """Retrieve a single project by its canonical ID or slug."""
  project(id: String!): Project

  """List organizations."""
  organizations(limit: Int): [Organization!]!

  """List recent events, optionally filtered by project."""
  events(projectId: String, limit: Int): [Event!]!

  """List momentum signals."""
  signals(limit: Int): [Signal!]!

  """Multi-jurisdiction reconciliation summary."""
  reconciliation: ReconciliationReport!

  """AI Sovereignty benchmark scores."""
  aiSovereignty: [AISovereigntyScore!]!

  """StatCan input-output macro multipliers for a project."""
  mrio(projectId: String!): MRIOImpact

  """Bayesian reference class cost and schedule overrun hazard curve."""
  flyvbjerg(projectId: String!): FlyvbjergRiskForecast

  """Ultimate Beneficial Ownership & Investment Canada Act screening."""
  ubo(projectId: String!): UBOScreening

  """Electrical grid feasibility and interconnect queue headroom."""
  grid(projectId: String!): GridAssessment

  """Satellite SAR and optical ground-truth corroboration dossier."""
  earthobs(projectId: String!): GroundTruthDossier

  """Sovereign capital allocation optimizer."""
  planningOptimize(objective: String): OptimizationResult!

  """Geopolitical macro-shock stress testing."""
  planningWarGame(scenario: String): WarGameResult!

  """Red Seal craft labor constraint and collision report."""
  planningLabor(province: String): RegionalLaborReport!

  """Machine-learned topological link prediction for partners and supply chain."""
  predictedLinks(projectId: String!): [PredictedLink!]!
}

"""
Subscription webhooks for moving-project notifications.
"""
type Subscription {
  """Fires when a project transitions stage or scores change."""
  projectUpdated(id: String!): ProjectUpdate!

  """Fires when the aggregate portfolio composition changes."""
  portfolioUpdated: PortfolioUpdate!
}

# ── Domain object types ──────────────────────────────────────────────────────

type Project {
  id:             String!
  name:           String!
  slug:           String!
  sector:         String!
  province:       String!
  stage:          String!
  capexCAD:       Float!
  buildability:   Float
  investability:  Float
  updatedAt:      String!
}

type Organization {
  id:         String!
  slug:       String!
  commonName: String!
  legalName:  String
  entityType: String!
  updatedAt:  String!
}

type Event {
  id:          String!
  projectId:   String!
  eventType:   String!
  eventDate:   String!
  title:       String!
  description: String
}

type Signal {
  id:          String!
  projectId:   String!
  signalType:  String!
  strength:    Float!
  detectedAt:  String!
  description: String
}

type ReconciliationReport {
  totalRecords: Int!
  merged:       Int!
  linked:       Int!
  conflicts:    Int!
}

type AISovereigntyScore {
  entityId:    String!
  entityName:  String!
  overallScore: Float!
  tier:         String!
}

type ProjectUpdate {
  projectId:  String!
  changeType: String!
  occurredAt: String!
}

type PortfolioUpdate {
  totalProjects: Int!
  updatedAt:     String!
}

type MRIOImpact {
  projectId:            String!
  projectName:          String!
  capexCAD:             Float!
  directGDPCAD:         Float!
  indirectGDPCAD:       Float!
  inducedGDPCAD:        Float!
  totalGDPCAD:          Float!
  totalMultiplier:       Float!
  personYearsJobs:      Int!
  federalTaxCAD:        Float!
  provincialTaxCAD:     Float!
  municipalTaxCAD:      Float!
  totalFiscalReturnCAD: Float!
  modelVersion:         String!
  auditHash:            String!
}

type FlyvbjergRiskForecast {
  projectId:              String!
  projectName:            String!
  sector:                 String!
  baseCapexCAD:           Float!
  referenceClass:         String!
  historicalSampleSize:   Int!
  expectedCostOverrunPct: Float!
  expectedDelayMonths:    Int!
  auditHash:              String!
}

type UBOScreening {
  projectId:            String!
  proponentName:        String!
  icaRisk:              String!
  domesticControlShare: Float!
  ftaPartnerShare:      Float!
  nonFTAShare:          Float!
  soeExposurePercent:   Float!
  criticalMineralFlag:  Boolean!
  dualUseSovereignty:   Boolean!
  auditHash:            String!
}

type GridAssessment {
  projectId:                 String!
  projectName:               String!
  province:                  String!
  systemOperator:            String!
  estimatedLoadOrGenMW:      Float!
  interconnectVoltageKV:     Int!
  queueEstimatedMonths:      Int!
  substationHeadroomMW:      Float!
  dedicatedSubstationNeeded: Boolean!
  reinforcementCostCAD:      Float!
  gridFeasibilityScore:      Float!
  cleanPowerPurityPct:       Float!
  auditHash:                 String!
}

type GroundTruthDossier {
  projectId:             String!
  claimedStage:          String!
  corroborationStatus:   String!
  physicalProgressScore: Float!
  earthworksConfirmed:   Boolean!
  structuresConfirmed:   Boolean!
  telemetrySummary:      String!
  auditHash:             String!
}

type OptimizationResult {
  requestId:               String!
  objective:               String!
  totalPublicInvestedCAD:  Float!
  totalPrivateMobilizedCAD: Float!
  crowdingInMultiplier:    Float!
  totalGHGAbatedMtYr:      Float!
  auditHash:               String!
}

type WarGameResult {
  simulationId:                String!
  scenario:                    String!
  scenarioTitle:               String!
  scenarioDescription:         String!
  totalAssetsStalledCount:     Int!
  totalFrozenCapexCAD:         Float!
  estimatedNationalGDPLossCAD: Float!
  auditHash:                   String!
}

type RegionalLaborReport {
  province:            String!
  totalActiveCapexCAD: Float!
  concurrentProjects:  Int!
  totalLaborDemandFTE: Int!
  collisionDetected:   Boolean!
  strategicAdvice:     String!
  auditHash:           String!
}

type PredictedLink {
  entityId:        String!
  entityName:      String!
  projectId:       String!
  projectName:     String!
  predictedRole:   String!
  confidenceScore: Float!
  adamicAdarScore: Float!
  rationale:       String!
  auditHash:       String!
}
`
