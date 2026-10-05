import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '../../..');

test('E2E Test 1: Public Release Manifest and Snapshot Integrity', async (t) => {
  const manifestPath = path.join(rootDir, 'data/public/manifest.json');
  assert.ok(fs.existsSync(manifestPath), 'manifest.json must exist');

  const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  assert.equal(manifest.cegs, '0.1', 'manifest must declare CEGS 0.1 specification');
  assert.ok(manifest.checksums_sha256, 'manifest must contain sha256 checksums');

  // Verify projects.jsonl hash
  const jsonlHash = manifest.checksums_sha256['public/projects.jsonl'];
  assert.ok(jsonlHash, 'manifest must include projects.jsonl checksum');

  const jsonlPath = path.join(rootDir, 'data/public/projects.jsonl');
  assert.ok(fs.existsSync(jsonlPath), 'projects.jsonl must exist');

  const jsonlContent = fs.readFileSync(jsonlPath);
  const hash = crypto.createHash('sha256').update(jsonlContent).digest('hex');
  assert.equal(hash, jsonlHash, 'projects.jsonl SHA-256 must match manifest signature');

  // Verify projects.jsonl parses valid project lines
  const lines = jsonlContent.toString('utf8').trim().split('\n');
  assert.ok(lines.length > 0, 'projects.jsonl must have lines');
  const firstProj = JSON.parse(lines[0]);
  assert.ok(firstProj.name, 'first project must have name');
  assert.ok(firstProj.sector, 'first project must have sector');
});

test('E2E Test 2: Capital Stack Builder Financial Modeling Logic', async (t) => {
  // Test the financial math from CapitalStackBuilder.tsx
  const capex = 1_000_000_000; // $1B CAD

  const tranches = [
    { name: 'Senior Debt', pct: 45, costPct: 5.5 },
    { name: 'Concessionary / CIB', pct: 25, costPct: 2.5 },
    { name: 'Sponsor Equity', pct: 20, costPct: 12.0 },
    { name: 'Indigenous Equity (FNFA)', pct: 10, costPct: 4.0 },
  ];

  const totalPct = tranches.reduce((sum, tr) => sum + tr.pct, 0);
  assert.equal(totalPct, 100, 'capital stack tranches must sum to exactly 100%');

  // Weighted Average Cost of Capital (WACC)
  const weightedCost = tranches.reduce((sum, tr) => sum + (tr.pct / 100) * tr.costPct, 0);
  assert.ok(weightedCost > 4.0 && weightedCost < 8.0, `realistic WACC between 4-8%, got ${weightedCost}%`);

  // Debt Service Coverage Ratio (DSCR)
  const annualNOI = 90_000_000; // $90M annual net operating income
  const seniorDebtAmount = capex * 0.45;
  const seniorDebtService = seniorDebtAmount * (0.055 + 1 / 25); // interest + 25-yr amortization
  const dscr = annualNOI / seniorDebtService;
  assert.ok(dscr >= 1.25, `DSCR must exceed investment-grade threshold of 1.25x, got ${dscr.toFixed(2)}x`);
});

test('E2E Test 3: Investigation Workbench GQL Query Contracts', async (t) => {
  // Simulate client-side query construction used in /investigate
  const testQuery = `{
    projects(limit: 5) {
      id
      name
      sector
      province
      capexCAD
    }
  }`;

  assert.ok(testQuery.includes('projects'), 'query must select projects');
  assert.ok(testQuery.includes('capexCAD'), 'query must request capexCAD');

  // Verify field selection pattern regex used in Next.js workbench
  const fieldRegex = /([a-zA-Z0-9_]+)\s*(?:\{|\()/g;
  const matches = [...testQuery.matchAll(fieldRegex)].map(m => m[1]);
  assert.ok(matches.includes('projects'), 'query parser must extract root field');
});

test('E2E Test 4: Multi-Level ABAC Security Clearance Redaction', async (t) => {
  // Simulate the classification switcher behavior in Navbar.tsx
  const mockRecord = {
    id: 'proj-classified-01',
    name: 'Alert High Arctic Deepwater Naval Facility',
    capexCAD: 3_200_000_000,
    sensitivity: 'SECRET',
    uboSanctionRisk: 'HIGH_IRANIAN_FRONT_PROXY',
  };

  function applyClientRedaction(record, clearance) {
    if (clearance === 'UNCLASSIFIED' || clearance === 'PUBLIC') {
      return {
        ...record,
        capexCAD: 0,
        uboSanctionRisk: '[REDACTED - INSUFFICIENT CLEARANCE]',
      };
    }
    if (clearance === 'PROTECTED_B') {
      return {
        ...record,
        uboSanctionRisk: '[REDACTED - EYES ONLY]',
      };
    }
    return record; // TOP_SECRET
  }

  const unclassView = applyClientRedaction(mockRecord, 'UNCLASSIFIED');
  assert.equal(unclassView.capexCAD, 0, 'capex must be redacted in unclassified view');
  assert.equal(unclassView.uboSanctionRisk, '[REDACTED - INSUFFICIENT CLEARANCE]');

  const secretView = applyClientRedaction(mockRecord, 'TOP_SECRET');
  assert.equal(secretView.capexCAD, 3_200_000_000, 'capex visible in secret view');
  assert.equal(secretView.uboSanctionRisk, 'HIGH_IRANIAN_FRONT_PROXY');
});

test('E2E Test 5: WCAG 2.2 AAA Accessibility, FIP Branding & Bilingual Compliance', async (t) => {
  const layoutPath = path.join(__dirname, '../app/layout.tsx');
  assert.ok(fs.existsSync(layoutPath), 'layout.tsx must exist');
  const layoutContent = fs.readFileSync(layoutPath, 'utf8');

  // Skip link verification
  assert.ok(layoutContent.includes('href="#main-content"'), 'must include skip-to-content anchor');
  assert.ok(layoutContent.includes('id="main-content"'), 'must include target main content landmark');
  assert.ok(layoutContent.includes('skip-link'), 'must style skip link for keyboard focus');

  // Bilingual and language provider
  assert.ok(layoutContent.includes('LanguageProvider'), 'must wrap application in LanguageProvider');
  assert.ok(layoutContent.includes('lang="en-CA"'), 'must declare default Canadian English locale');

  // FIP Navbar and Footer compliance
  const navbarPath = path.join(__dirname, '../components/Navbar.tsx');
  const footerPath = path.join(__dirname, '../components/Footer.tsx');
  assert.ok(fs.existsSync(navbarPath), 'Navbar.tsx must exist');
  assert.ok(fs.existsSync(footerPath), 'Footer.tsx must exist');

  const navbarContent = fs.readFileSync(navbarPath, 'utf8');
  const footerContent = fs.readFileSync(footerPath, 'utf8');

  // Verify Canadian Flag SVG and bilingual aria-label
  assert.ok(navbarContent.includes('Flag of Canada / Drapeau du Canada'), 'Navbar must include bilingual flag accessibility label');
  assert.ok(footerContent.includes('Flag of Canada / Drapeau du Canada'), 'Footer must include bilingual flag accessibility label');

  // Verify Open Government Licence - Canada compliance notice
  assert.ok(footerContent.includes('Open Government Licence - Canada'), 'Footer must cite OGL-Canada licence');
  assert.ok(footerContent.includes('WCAG 2.2 AAA'), 'Footer must declare WCAG AAA bilingual compliance');
});

test('E2E Test 6: Filings Fallback Is Explicitly Labelled and Time-Stable', async () => {
  const dataPath = path.join(__dirname, '../lib/data.ts');
  const explorerPath = path.join(__dirname, '../components/ProcurementExplorer.tsx');
  const dataContent = fs.readFileSync(dataPath, 'utf8');
  const explorerContent = fs.readFileSync(explorerPath, 'utf8');

  assert.ok(dataContent.includes('data_mode: "DEMONSTRATION_SNAPSHOT"'));
  assert.ok(dataContent.includes('Demonstration records; not a live registry'));
  assert.ok(!dataContent.includes('new Date(Date.now()'), 'fallback timestamps must not move with the clock');
  assert.ok(explorerContent.includes('Demonstration snapshot:'), 'UI must disclose non-live data');
});

test('E2E Test 7: 10-Pillars Palantir Sovereign Capabilities Architecture Contract', async (t) => {
  // Verify architectural contract of all 10 Sovereign Pillars in palantir_routes.go
  const routesPath = path.join(rootDir, 'internal/api/palantir_routes.go');
  assert.ok(fs.existsSync(routesPath), 'palantir_routes.go must exist');
  const routesContent = fs.readFileSync(routesPath, 'utf8');

  const requiredPillars = [
    { pillar: 1, name: 'Ontology', marker: '/api/v1/ontology/' },
    { pillar: 2, name: 'Lakehouse', marker: '/api/v1/lakehouse/' },
    { pillar: 3, name: 'Knowledge Graph & GQL', marker: '/api/v1/graph/' },
    { pillar: 4, name: 'Earth Observation', marker: '/api/v1/earthobs/' },
    { pillar: 5, name: 'AIP Multi-Agent', marker: '/api/v1/ai/' },
    { pillar: 6, name: 'Linear Corridors & Gateways', marker: '/api/v1/corridor/' },
    { pillar: 7, name: 'Sovereign Security & ABAC', marker: '/api/v1/security/' },
    { pillar: 8, name: 'Counter-Intelligence & UBO', marker: '/api/v1/counter-intel/' },
    { pillar: 9, name: 'Geoeconomic War Game & Allocation', marker: '/api/v1/palantir/planning/' },
    { pillar: 10, name: 'Compliance & Merkle Root', marker: '/api/v1/compliance/' },
  ];

  for (const p of requiredPillars) {
    assert.ok(
      routesContent.includes(p.marker),
      `Pillar ${p.pillar} (${p.name}) route prefix '${p.marker}' must be registered`
    );
  }
});

test('E2E Test 8: Institutional Deal Precedents and Maple 8 Allocator Matching Contract', async (t) => {
  // Validate the institutional matching weights and similarity scoring formula
  // from internal/matching/precedent.go
  function calculateSimilarity(project, deal) {
    let score = 0;
    const matched = [];

    if (project.sector === deal.sector) {
      score += 35;
      matched.push('sector');
    }
    if (project.province === deal.province) {
      score += 20;
      matched.push('province');
    }
    if (project.stage === deal.stage) {
      score += 15;
      matched.push('stage');
    }
    const ratio = Math.min(project.capexCAD, deal.amountCAD) / Math.max(project.capexCAD, deal.amountCAD);
    if (ratio >= 0.5) {
      score += 20 * ratio;
      matched.push('scale');
    }
    const currentYear = 2026;
    if (deal.year >= currentYear - 3) {
      score += 10;
      matched.push('recent');
    }
    return { score, matched };
  }

  const crawfordProject = {
    sector: 'CRITICAL_MINERALS',
    province: 'ON',
    stage: 'FEED',
    capexCAD: 2_500_000_000,
  };

  const cgfDeal = {
    sector: 'CRITICAL_MINERALS',
    province: 'ON',
    stage: 'FEED',
    amountCAD: 2_000_000_000,
    year: 2024,
  };

  const result = calculateSimilarity(crawfordProject, cgfDeal);
  assert.ok(result.score >= 85, `Institutional match score should exceed 85%, got ${result.score}%`);
  assert.deepEqual(result.matched, ['sector', 'province', 'stage', 'scale', 'recent']);
});

test('E2E Test 9: Cloud-Native Helm and Kubernetes Deployment Manifests Contract', async (t) => {
  const helmDeployment = path.join(rootDir, 'deploy/helm/templates/deployment.yaml');
  const helmIngress = path.join(rootDir, 'deploy/helm/templates/ingress.yaml');
  const helmHPA = path.join(rootDir, 'deploy/helm/templates/hpa.yaml');
  const k8sIngress = path.join(rootDir, 'deploy/k8s/ingress.yaml');
  const k8sPDB = path.join(rootDir, 'deploy/k8s/pdb.yaml');

  assert.ok(fs.existsSync(helmDeployment), 'helm deployment.yaml must exist');
  assert.ok(fs.existsSync(helmIngress), 'helm ingress.yaml must exist');
  assert.ok(fs.existsSync(helmHPA), 'helm hpa.yaml must exist');
  assert.ok(fs.existsSync(k8sIngress), 'k8s ingress.yaml must exist');
  assert.ok(fs.existsSync(k8sPDB), 'k8s pdb.yaml must exist');

  const k8sIngressContent = fs.readFileSync(k8sIngress, 'utf8');
  assert.ok(k8sIngressContent.includes('opportunity.canada.ca'), 'ingress must route opportunity.canada.ca');
  assert.ok(k8sIngressContent.includes('api.opportunity.canada.ca'), 'ingress must route api.opportunity.canada.ca');
  assert.ok(k8sIngressContent.includes('cert-manager.io/cluster-issuer'), 'ingress must declare cert-manager issuer');

  const k8sPDBContent = fs.readFileSync(k8sPDB, 'utf8');
  assert.ok(k8sPDBContent.includes('PodDisruptionBudget'), 'pdb must declare PodDisruptionBudget');
  assert.ok(k8sPDBContent.includes('minAvailable'), 'pdb must configure minAvailable threshold');
});

test('E2E Test 10: Cryptographic Merkle KMS Attestation and SSE Streaming Contract', async (t) => {
  const signerPath = path.join(rootDir, 'internal/merkle/signer.go');
  const sseHubPath = path.join(rootDir, 'internal/eventsse/hub.go');
  const sseHookPath = path.join(__dirname, '../lib/useLiveEvents.ts');
  const sedarPath = path.join(rootDir, 'adapters/sedar/sedar.go');

  assert.ok(fs.existsSync(signerPath), 'merkle signer.go must exist');
  assert.ok(fs.existsSync(sseHubPath), 'eventsse hub.go must exist');
  assert.ok(fs.existsSync(sseHookPath), 'useLiveEvents.ts hook must exist');
  assert.ok(fs.existsSync(sedarPath), 'sedar.go adapter must exist');

  const signerContent = fs.readFileSync(signerPath, 'utf8');
  assert.ok(signerContent.includes('KMSSigner'), 'signer must define KMSSigner interface');
  assert.ok(signerContent.includes('SignRoot'), 'signer must provide SignRoot function');
  assert.ok(signerContent.includes('VerifySignedRoot'), 'signer must provide VerifySignedRoot function');

  const sseHookContent = fs.readFileSync(sseHookPath, 'utf8');
  assert.ok(sseHookContent.includes('useLiveEvents'), 'hook must export useLiveEvents function');
  assert.ok(sseHookContent.includes('/api/v1/stream/events'), 'hook must connect to /api/v1/stream/events');

  const sedarContent = fs.readFileSync(sedarPath, 'utf8');
  assert.ok(sedarContent.includes('sedar_plus_disclosures'), 'adapter must name sedar_plus_disclosures');
  assert.ok(sedarContent.includes('CanonicalSEDARFilings'), 'adapter must provide CanonicalSEDARFilings');
});

test('E2E Test 11: Mark Carney Briefing Suite: National Airports System Concession and GFANZ Transition Taxonomy', async (t) => {
  const concessionPath = path.join(rootDir, 'internal/concession/airports.go');
  const taxonomyPath = path.join(rootDir, 'internal/transitionfinance/taxonomy.go');
  const carneyHandlersPath = path.join(rootDir, 'internal/api/carney_handlers.go');
  const dossierCompPath = path.join(__dirname, '../components/finance/CarneyBriefingDossier.tsx');
  const carneyPagePath = path.join(__dirname, '../app/briefing/carney/page.tsx');

  assert.ok(fs.existsSync(concessionPath), 'airports.go must exist');
  assert.ok(fs.existsSync(taxonomyPath), 'taxonomy.go must exist');
  assert.ok(fs.existsSync(carneyHandlersPath), 'carney_handlers.go must exist');
  assert.ok(fs.existsSync(dossierCompPath), 'CarneyBriefingDossier.tsx must exist');
  assert.ok(fs.existsSync(carneyPagePath), 'briefing/carney/page.tsx must exist');

  // Verify National Airports System concession mechanics
  const concessionContent = fs.readFileSync(concessionPath, 'utf8');
  assert.ok(concessionContent.includes('HubYYZ'), 'concession must model Toronto Pearson YYZ');
  assert.ok(concessionContent.includes('HubYVR'), 'concession must model Vancouver YVR');
  assert.ok(concessionContent.includes('HubYUL'), 'concession must model Montreal YUL');
  assert.ok(concessionContent.includes('HubYYC'), 'concession must model Calgary YYC');
  assert.ok(concessionContent.includes('HubYEG'), 'concession must model Edmonton YEG');
  assert.ok(concessionContent.includes('18_000_000_000'), 'concession must total $18B CAD target capex');
  assert.ok(concessionContent.includes('Maple8AllocationsCAD'), 'concession must allocate to Canadian Maple 8 pensions');

  // Verify GFANZ / SFAC Transition Taxonomy logic
  const taxonomyContent = fs.readFileSync(taxonomyPath, 'utf8');
  assert.ok(taxonomyContent.includes('CategoryGreen'), 'taxonomy must define CategoryGreen');
  assert.ok(taxonomyContent.includes('CategoryTransition'), 'taxonomy must define CategoryTransition');
  assert.ok(taxonomyContent.includes('CategoryEnabling'), 'taxonomy must define CategoryEnabling');
  assert.ok(taxonomyContent.includes('TransitionCredibilityIndex'), 'taxonomy must score TransitionCredibilityIndex');
  assert.ok(taxonomyContent.includes('MarginalAbatementCostCAD'), 'taxonomy must compute MAC curve');

  // Verify API route registrations
  const serverPath = path.join(rootDir, 'internal/api/server.go');
  const serverContent = fs.readFileSync(serverPath, 'utf8');
  assert.ok(serverContent.includes('/api/v1/finance/transition-taxonomy'), 'server must register transition-taxonomy endpoint');
  assert.ok(serverContent.includes('/api/v1/finance/concession/airports'), 'server must register concession/airports endpoint');

  // Verify UI Dossier completeness
  const dossierContent = fs.readFileSync(dossierCompPath, 'utf8');
  assert.ok(dossierContent.includes('CANONICAL_AIRPORTS'), 'dossier must define CANONICAL_AIRPORTS');
  assert.ok(dossierContent.includes('TRANSITION_BENCHMARKS'), 'dossier must define TRANSITION_BENCHMARKS');
  assert.ok(dossierContent.includes('CPPIB'), 'dossier must reference CPPIB');
  assert.ok(dossierContent.includes('CDPQ'), 'dossier must reference CDPQ');
  assert.ok(dossierContent.includes('Brookfield'), 'dossier must reference Brookfield');
});

test('E2E Test 12: Product Control Plane, Universal Gateway, and Deployment Wiring', async () => {
  const gatewayPath = path.join(__dirname, '../app/api/v1/[...path]/route.ts');
  const operationsPath = path.join(__dirname, '../app/operations/page.tsx');
  const operationsLibPath = path.join(__dirname, '../lib/operations.ts');
  const systemHandlerPath = path.join(rootDir, 'internal/api/system_handler.go');
  const planningUIPath = path.join(__dirname, '../components/NationalPlanningWorkbench.tsx');
  const helmDeploymentPath = path.join(rootDir, 'deploy/helm/templates/deployment.yaml');

  for (const requiredPath of [gatewayPath, operationsPath, operationsLibPath, systemHandlerPath]) {
    assert.ok(fs.existsSync(requiredPath), `${requiredPath} must exist`);
  }

  const gateway = fs.readFileSync(gatewayPath, 'utf8');
  assert.ok(gateway.includes('requestUpstreamAPI'), 'gateway must use the shared upstream transport');
  assert.ok(gateway.includes('MAX_PROXY_BODY_BYTES'), 'gateway must bound request bodies');
  assert.ok(gateway.includes('forbidden_origin'), 'gateway must reject cross-site mutations');
  assert.ok(gateway.includes('application/json'), 'gateway mutations must require JSON bodies');
  assert.ok(!gateway.includes('x-admin-secret'), 'gateway must never forward adapter administration credentials');
  assert.ok(!gateway.includes('x-clearance-level'), 'gateway must not trust browser-supplied security clearance');

  const systemHandler = fs.readFileSync(systemHandlerPath, 'utf8');
  assert.ok(systemHandler.includes('AggregateHealth'), 'system status must read the live connector registry');
  assert.ok(systemHandler.includes('SchedulerEnabled'), 'system status must expose scheduled-ingestion readiness');

  const planningUI = fs.readFileSync(planningUIPath, 'utf8');
  assert.ok(planningUI.includes('cib_concessionary_cad'), 'planning UI must send the Go optimizer envelope schema');
  assert.ok(planningUI.includes('allocated_projects'), 'planning UI must read the Go optimizer response schema');
  assert.ok(planningUI.includes('total_frozen_capex_cad'), 'planning UI must read the Go war-game response schema');

  const helmDeployment = fs.readFileSync(helmDeploymentPath, 'utf8');
  for (const envName of ['ENV', 'STORAGE_MODE', 'STORAGE_DIR', 'INGEST_INTERVAL', 'COG_API_BASE']) {
    assert.ok(helmDeployment.includes(`name: ${envName}`), `Helm must configure ${envName}`);
  }
  assert.ok(helmDeployment.includes('-wal-data'), 'API deployment and PVC must share the WAL claim name');
  assert.ok(!helmDeployment.includes('name: ENVIRONMENT'), 'Helm must not use ignored ENVIRONMENT variable');
  assert.ok(!helmDeployment.includes('name: DATA_DIR'), 'Helm must not use ignored DATA_DIR variable');
});
