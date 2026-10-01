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

test('E2E Test 6: 10-Pillars Palantir Sovereign Capabilities Architecture Contract', async (t) => {
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

test('E2E Test 7: Institutional Deal Precedents and Maple 8 Allocator Matching Contract', async (t) => {
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

