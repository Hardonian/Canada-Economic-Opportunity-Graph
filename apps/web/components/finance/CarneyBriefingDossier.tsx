"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import {
  Plane,
  Leaf,
  Layers,
  ShieldCheck,
  TrendingUp,
  Percent,
  DollarSign,
  Building2,
  Sliders,
  FileText,
  Printer,
  Copy,
  Check,
  ArrowRight,
  Sparkles,
  Award,
  Zap,
  CheckCircle2,
  AlertCircle,
  HelpCircle,
} from "lucide-react";

export interface AirportHub {
  code: string;
  name: string;
  province: string;
  paxBaseline: number;
  pax2050: number;
  revBaselineCad: number;
  ebitdaBaselineCad: number;
  targetCapexCad: number;
  keyProjects: string[];
  decarbInitiatives: string[];
  targetPensions: string[];
}

const CANONICAL_AIRPORTS: AirportHub[] = [
  {
    code: "YYZ",
    name: "Toronto Pearson International Airport",
    province: "ON",
    paxBaseline: 45_000_000,
    pax2050: 68_000_000,
    revBaselineCad: 1_650_000_000,
    ebitdaBaselineCad: 780_000_000,
    targetCapexCad: 7_500_000_000,
    keyProjects: [
      "Union Station West Pearson Transit Hub integration",
      "Terminal 1 & 3 digital apron and biometric processing overhaul",
      "Global Air Cargo logistics mega-terminal with cold-chain pharmaceutical hub",
      "Southern Ontario high-frequency rail intermodal passenger station",
    ],
    decarbInitiatives: [
      "Dedicated pipeline & hydrant infrastructure for 100% Sustainable Aviation Fuel (SAF)",
      "Airfield apron conversion for 100% electric ground support equipment (e-GSE)",
      "150MW geothermal heating and on-site district energy microgrid",
    ],
    targetPensions: ["CPPIB", "OMERS Infrastructure", "Brookfield Infrastructure Partners"],
  },
  {
    code: "YVR",
    name: "Vancouver International Airport",
    province: "BC",
    paxBaseline: 26_000_000,
    pax2050: 39_000_000,
    revBaselineCad: 720_000_000,
    ebitdaBaselineCad: 350_000_000,
    targetCapexCad: 3_800_000_000,
    keyProjects: [
      "Trans-Pacific Air Cargo Logistics City expansion",
      "International Terminal Building pier D/E expansion",
      "Sea Island autonomous transit intertie and Canada Line frequency enhancement",
    ],
    decarbInitiatives: [
      "Pacific Rim SAF bunkering hub with Prince Rupert & Vancouver port interties",
      "Sea Island geo-exchange heating loop and 45MW solar apron installations",
      "Shore power electrification for widebody freighters at airside gates",
    ],
    targetPensions: ["BCI", "CDPQ", "PSP Investments"],
  },
  {
    code: "YUL",
    name: "Montréal-Pierre Elliott Trudeau International Airport",
    province: "QC",
    paxBaseline: 21_500_000,
    pax2050: 33_000_000,
    revBaselineCad: 620_000_000,
    ebitdaBaselineCad: 295_000_000,
    targetCapexCad: 3_500_000_000,
    keyProjects: [
      "Direct subterranean Réseau express métropolitain (REM) station completion",
      "New 14-gate international boarding pier with automated baggage routing",
      "Dorval intermodal freight interchange and customs pre-clearance center",
    ],
    decarbInitiatives: [
      "Hydro-Québec 100% clean power electrified taxiway systems",
      "Biokerosene SAF distribution network from Montreal East refining corridor",
      "Zero-carbon terminal envelope retrofits and rainwater harvesting",
    ],
    targetPensions: ["CDPQ", "Brookfield Infrastructure Partners", "PSP Investments"],
  },
  {
    code: "YYC",
    name: "Calgary International Airport",
    province: "AB",
    paxBaseline: 18_500_000,
    pax2050: 27_000_000,
    revBaselineCad: 460_000_000,
    ebitdaBaselineCad: 220_000_000,
    targetCapexCad: 2_000_000_000,
    keyProjects: [
      "Banff-Calgary Airport express passenger rail terminal integration",
      "Western Canadian energy & agritech air cargo processing center",
      "Cross-runway high-speed taxiway connectors and de-icing bay automation",
    ],
    decarbInitiatives: [
      "Alberta hydrogen aero-fueling pilot with Edmonton corridor link",
      "Electrified glycol recapture and recycling treatment facility",
      "Airport boundary solar arrays generating 60MW peak clean generation",
    ],
    targetPensions: ["AIMCo", "CPPIB"],
  },
  {
    code: "YEG",
    name: "Edmonton International Airport",
    province: "AB",
    paxBaseline: 8_200_000,
    pax2050: 14_500_000,
    revBaselineCad: 280_000_000,
    ebitdaBaselineCad: 130_000_000,
    targetCapexCad: 1_200_000_000,
    keyProjects: [
      "Airport City Sustainability Campus cargo logistics hub",
      "Northern Canada Arctic resupply and dual-use aerospace logistics apron",
      "Hydrogen commercial hub and zero-emission maintenance hangars",
    ],
    decarbInitiatives: [
      "627-acre on-site solar farm (Airport City Solar - 120MW)",
      "Commercial hydrogen passenger shuttle fleet and fuel-cell ground support",
      "SAF blending facility connected to Alberta bio-refineries",
    ],
    targetPensions: ["AIMCo", "OMERS Infrastructure"],
  },
];

const TRANSITION_BENCHMARKS = [
  {
    id: "darlington-smr",
    name: "Darlington New Nuclear Project — Unit 1",
    sector: "Nuclear & Clean Power",
    capexCad: 7_700_000_000,
    category: "GREEN",
    pillar: "Climate Solutions (Zero-Emissions Baseload)",
    credibility: 97.0,
    rating: "HIGH",
    annualAbatementTpy: 3_395_700,
    lifetimeAbatementTonnes: 203_742_000,
    macCad: 37.79,
    noLockIn: true,
  },
  {
    id: "crawford-nickel",
    name: "Crawford Nickel-Cobalt Sulphide Project",
    sector: "Critical Minerals & CCUS",
    capexCad: 3_500_000_000,
    category: "TRANSITION",
    pillar: "Aligned Heavy Extraction with Mineral Carbonation",
    credibility: 93.0,
    rating: "HIGH",
    annualAbatementTpy: 1_050_000,
    lifetimeAbatementTonnes: 31_500_000,
    macCad: 111.11,
    noLockIn: true,
  },
  {
    id: "nas-airports",
    name: "Canadian International Airports — NAS Leasing Hubs",
    sector: "Transportation & Intermodal",
    capexCad: 18_000_000_000,
    category: "ENABLING",
    pillar: "Modernized Intermodal & Low-Carbon Aviation Hubs",
    credibility: 79.5,
    rating: "MEDIUM",
    annualAbatementTpy: 3_240_000,
    lifetimeAbatementTonnes: 129_600_000,
    macCad: 138.89,
    noLockIn: true,
  },
  {
    id: "oneida-storage",
    name: "Oneida Energy Storage Project",
    sector: "Clean Energy & Grid",
    capexCad: 500_000_000,
    category: "GREEN",
    pillar: "Climate Solutions (Battery Storage)",
    credibility: 94.5,
    rating: "HIGH",
    annualAbatementTpy: 145_500,
    lifetimeAbatementTonnes: 4_365_000,
    macCad: 114.55,
    noLockIn: true,
  },
  {
    id: "chisasibi-ai",
    name: "Chisasibi Sovereign AI Hyperscale Cluster",
    sector: "AI Compute & Digital",
    capexCad: 2_800_000_000,
    category: "ENABLING",
    pillar: "Sovereign Low-Carbon Digital Infrastructure",
    credibility: 91.0,
    rating: "HIGH",
    annualAbatementTpy: 399_000,
    lifetimeAbatementTonnes: 7_980_000,
    macCad: 350.88,
    noLockIn: true,
  },
];

function formatCad(value: number): string {
  if (value >= 1e9) return `$${(value / 1e9).toFixed(2)}B`;
  if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`;
  return `$${value.toLocaleString("en-CA")}`;
}

export default function CarneyBriefingDossier() {
  const [activeTab, setActiveTab] = useState<"nas_concession" | "transition_taxonomy" | "action_plan">(
    "nas_concession"
  );

  // Concession Simulator State
  const [concessionYears, setConcessionYears] = useState<number>(40);
  const [federalRoyaltyPct, setFederalRoyaltyPct] = useState<number>(10.0);
  const [pensionEquityPct, setPensionEquityPct] = useState<number>(50.0);
  const [commercialDebtPct, setCommercialDebtPct] = useState<number>(45.0);
  const [passengerCagrPct, setPassengerCagrPct] = useState<number>(2.8);
  const [discountRatePct, setDiscountRatePct] = useState<number>(5.5);
  const [selectedHub, setSelectedHub] = useState<string>("YYZ");

  // Taxonomy State
  const [selectedTaxonomyProject, setSelectedTaxonomyProject] = useState<string>("nas-airports");
  const [copiedMemo, setCopiedMemo] = useState<boolean>(false);

  // Computed NAS Concession Financials
  const concessionSim = useMemo(() => {
    const totalCapex = CANONICAL_AIRPORTS.reduce((sum, a) => sum + a.targetCapexCad, 0);
    const pensionEquity = (totalCapex * pensionEquityPct) / 100;
    const commercialDebt = (totalCapex * commercialDebtPct) / 100;
    const fedSubordinate = (totalCapex * 5.0) / 100;

    // Upfront federal proceeds: 3.75x baseline EBITDA capitalized
    const upfrontProceeds = CANONICAL_AIRPORTS.reduce((sum, a) => sum + Math.round(a.ebitdaBaselineCad * 3.75), 0);

    // Cumulative 40-year royalty projection
    const growth = passengerCagrPct / 100 + 0.02; // traffic + 2% inflation
    let cumulativeRoyalties = 0;
    const royaltyRate = federalRoyaltyPct / 100;

    for (let y = 1; y <= concessionYears; y++) {
      const yearMultiplier = Math.pow(1 + growth, y - 1);
      for (const a of CANONICAL_AIRPORTS) {
        cumulativeRoyalties += a.revBaselineCad * yearMultiplier * royaltyRate;
      }
    }

    const baseIRR = 9.8;
    const equityIRR = baseIRR + (pensionEquityPct <= 50 ? 2.1 : 1.4);
    const avgDSCR = 1.88;
    const crowdingIn = (pensionEquity + commercialDebt) / fedSubordinate;

    // Maple 8 Allocation Breakdown
    const maple8 = [
      { name: "CPPIB (Canada Pension Plan Investment Board)", allocation: pensionEquity * 0.28, pct: 28 },
      { name: "CDPQ (Caisse de dépôt et placement du Québec)", allocation: pensionEquity * 0.22, pct: 22 },
      { name: "Brookfield Infrastructure Partners", allocation: pensionEquity * 0.18, pct: 18 },
      { name: "OMERS Infrastructure", allocation: pensionEquity * 0.14, pct: 14 },
      { name: "PSP Investments", allocation: pensionEquity * 0.08, pct: 8 },
      { name: "AIMCo (Alberta Investment Management Corp)", allocation: pensionEquity * 0.05, pct: 5 },
      { name: "BCI (British Columbia Investment Management)", allocation: pensionEquity * 0.05, pct: 5 },
    ];

    // Green infrastructure spend within concession (30%)
    const greenCapexTotal = totalCapex * 0.3;

    return {
      totalCapex,
      pensionEquity,
      commercialDebt,
      fedSubordinate,
      upfrontProceeds,
      cumulativeRoyalties: Math.round(cumulativeRoyalties),
      baseIRR,
      equityIRR,
      avgDSCR,
      crowdingIn,
      maple8,
      greenCapexTotal,
    };
  }, [concessionYears, federalRoyaltyPct, pensionEquityPct, commercialDebtPct, passengerCagrPct]);

  const activeHubDetails = useMemo(() => {
    return CANONICAL_AIRPORTS.find((a) => a.code === selectedHub) || CANONICAL_AIRPORTS[0];
  }, [selectedHub]);

  const activeTaxonomyDetails = useMemo(() => {
    return TRANSITION_BENCHMARKS.find((p) => p.id === selectedTaxonomyProject) || TRANSITION_BENCHMARKS[0];
  }, [selectedTaxonomyProject]);

  const handleCopyMemo = async () => {
    const memoText = `CONFIDENTIAL // PRIVY COUNCIL OFFICE // ECONOMIC BRIEFING
TO: The Right Honourable Mark Carney, Chair, Advisory Council on Economic Growth
SUBJECT: National Airports System (NAS) Ground Lease Modernization & GFANZ Transition Finance Architecture
DATE: ${new Date().toLocaleDateString("en-CA", { dateStyle: "long" })}

EXECUTIVE SUMMARY:
- Target Capital Program: $18.0B CAD across 5 Flagship Gateway Hubs (YYZ, YVR, YUL, YYC, YEG).
- Restructuring Model: Converting legacy 1994 ground leases into 40-year institutional concession tranches.
- Upfront Federal Proceeds: ${formatCad(concessionSim.upfrontProceeds)} upfront treasury injection.
- Cumulative 40-Year Federal Royalties: ${formatCad(concessionSim.cumulativeRoyalties)} non-tax federal revenue.
- Pension Equity Mobilized: ${formatCad(concessionSim.pensionEquity)} anchored by Canadian Maple 8 (CPPIB, CDPQ, OMERS, Brookfield, PSP, BCI, AIMCo).
- Private Capital Crowding-In Multiplier: ${concessionSim.crowdingIn.toFixed(1)}x.
- Green Infrastructure Envelope: ${formatCad(concessionSim.greenCapexTotal)} dedicated to SAF hydrants, high-speed rail/REM transit interties, and airport microgrids.
`;
    try {
      await navigator.clipboard.writeText(memoText);
      setCopiedMemo(true);
      setTimeout(() => setCopiedMemo(false), 2500);
    } catch {
      // fallback
    }
  };

  return (
    <div className="space-y-8">
      {/* Executive Security & PCO Classification Header */}
      <div className="rounded-xl border border-red-500/40 bg-[#160608]/90 p-4 font-mono shadow-2xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-red-900/60 pb-3">
          <div className="flex items-center gap-2 text-red-400 font-bold tracking-wider text-xs">
            <span className="inline-block px-2 py-0.5 rounded bg-red-950 text-red-300 border border-red-800 text-[11px]">
              CONFIDENTIEL // HIGHLY PRIVILEGED
            </span>
            <span>ECONOMIC COUNCIL BRIEFING // CARNEY INITIATIVE SUITE</span>
          </div>
          <div className="text-xs text-text-subtle">
            REF: <span className="text-white font-semibold">CARNEY-NAS-GFANZ-2026</span> · CEGS 0.1 SPEC
          </div>
        </div>
        <div className="mt-2.5 flex flex-wrap items-center justify-between gap-2 text-xs text-text-subtle">
          <div>
            <span className="text-text-muted">PRINCIPAL:</span> The Right Honourable Mark Carney, P.C., O.C.
          </div>
          <div>
            <span className="text-text-muted">SCOPE:</span> National Airports System ($18B Concession) & GFANZ Canadian Transition Taxonomy
          </div>
        </div>
      </div>

      {/* Main Command Navigation */}
      <div className="flex flex-wrap items-center gap-2 border-b border-borderSubtle pb-4">
        <button
          type="button"
          onClick={() => setActiveTab("nas_concession")}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-xl font-medium text-sm transition-all ${
            activeTab === "nas_concession"
              ? "bg-gold text-surfaceDark shadow-lg shadow-gold/20 font-semibold"
              : "bg-surface hover:bg-card text-text-subtle hover:text-text-main border border-borderSubtle"
          }`}
        >
          <Plane className="h-4 w-4" />
          <span>NAS Airport Concessions ($18B)</span>
        </button>

        <button
          type="button"
          onClick={() => setActiveTab("transition_taxonomy")}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-xl font-medium text-sm transition-all ${
            activeTab === "transition_taxonomy"
              ? "bg-aurora text-surfaceDark shadow-lg shadow-aurora/20 font-semibold"
              : "bg-surface hover:bg-card text-text-subtle hover:text-text-main border border-borderSubtle"
          }`}
        >
          <Leaf className="h-4 w-4" />
          <span>GFANZ / SFAC Transition Taxonomy</span>
        </button>

        <button
          type="button"
          onClick={() => setActiveTab("action_plan")}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-xl font-medium text-sm transition-all ${
            activeTab === "action_plan"
              ? "bg-primary text-surfaceDark shadow-lg shadow-primary/20 font-semibold"
              : "bg-surface hover:bg-card text-text-subtle hover:text-text-main border border-borderSubtle"
          }`}
        >
          <FileText className="h-4 w-4" />
          <span>Executive Action Plan & MC Annex</span>
        </button>
      </div>

      {/* TAB 1: NAS AIRPORT CONCESSION SIMULATOR */}
      {activeTab === "nas_concession" && (
        <div className="space-y-8 animate-fadeIn">
          {/* Header Overview Card */}
          <div className="relative overflow-hidden rounded-2xl border border-gold/30 bg-gradient-to-br from-card via-surface to-[#161205] p-6 shadow-2xl">
            <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
              <div className="space-y-2 max-w-3xl">
                <div className="inline-flex items-center gap-2 rounded-full border border-gold/40 bg-gold/10 px-3 py-1 text-xs font-mono font-bold text-gold">
                  <Plane className="h-3.5 w-3.5" />
                  SUMMIT DIRECTIVE 4: AIRPORT COMMERCIAL LEASE MODERNIZATION
                </div>
                <h2 className="text-2xl sm:text-3xl font-black text-text-main">
                  National Airports System <span className="text-gold">Concession Simulator</span>
                </h2>
                <p className="text-sm text-text-muted leading-relaxed">
                  Restructuring Transport Canada’s 1994 airport ground leases into 40-year institutional concession tranches. Mobilizes $18B CAD in private capital from Canadian pension funds (Maple 8) for terminal modernization, high-frequency transit interties, and sustainable aviation fuel (SAF) infrastructure without federal debt.
                </p>
              </div>

              <div className="flex flex-col items-end justify-center rounded-xl border border-borderSubtle bg-surface/80 p-4 text-right">
                <div className="text-xs font-mono text-text-subtle uppercase">Target Program Capex</div>
                <div className="text-3xl font-black text-gold font-mono">{formatCad(concessionSim.totalCapex)}</div>
                <div className="text-[11px] text-aurora font-mono mt-1">Across 5 National Gateway Hubs</div>
              </div>
            </div>

            {/* Key Metric Highlights */}
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3 mt-6 pt-6 border-t border-borderSubtle/60">
              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Upfront Fed Proceeds</div>
                <div className="text-lg font-bold text-aurora font-mono mt-0.5">{formatCad(concessionSim.upfrontProceeds)}</div>
                <div className="text-[10px] text-text-subtle">Treasury Capitalization</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Cumulative Royalties</div>
                <div className="text-lg font-bold text-gold font-mono mt-0.5">{formatCad(concessionSim.cumulativeRoyalties)}</div>
                <div className="text-[10px] text-text-subtle">{concessionYears}-Yr Treasury Inflow</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Pension Equity</div>
                <div className="text-lg font-bold text-primary font-mono mt-0.5">{formatCad(concessionSim.pensionEquity)}</div>
                <div className="text-[10px] text-text-subtle">Maple 8 Allocation</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Crowding-In Ratio</div>
                <div className="text-lg font-bold text-emerald-400 font-mono mt-0.5">{concessionSim.crowdingIn.toFixed(1)}x</div>
                <div className="text-[10px] text-text-subtle">Private vs Fed Balance</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Base Equity IRR</div>
                <div className="text-lg font-bold text-purple-400 font-mono mt-0.5">{concessionSim.equityIRR.toFixed(1)}%</div>
                <div className="text-[10px] text-text-subtle">Pension Core Yield</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/40 p-3">
                <div className="text-[10px] font-mono text-text-subtle uppercase">Green Capex Share</div>
                <div className="text-lg font-bold text-cyan-400 font-mono mt-0.5">{formatCad(concessionSim.greenCapexTotal)}</div>
                <div className="text-[10px] text-text-subtle">SAF / Transit / Solar</div>
              </div>
            </div>
          </div>

          {/* Interactive Simulation Controls */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="rounded-2xl border border-borderSubtle bg-card p-6 space-y-6">
              <div className="flex items-center gap-2 border-b border-borderSubtle pb-3">
                <Sliders className="h-5 w-5 text-gold" />
                <h3 className="font-bold text-text-main text-base">Concession Underwriting Levers</h3>
              </div>

              {/* Slider 1: Concession Horizon */}
              <div className="space-y-2">
                <div className="flex justify-between text-xs font-mono">
                  <span className="text-text-muted">Concession Term:</span>
                  <span className="text-gold font-bold">{concessionYears} Years</span>
                </div>
                <input
                  type="range"
                  min="30"
                  max="50"
                  step="5"
                  value={concessionYears}
                  onChange={(e) => setConcessionYears(Number(e.target.value))}
                  className="w-full accent-gold h-1.5 bg-surface rounded-lg cursor-pointer"
                />
                <div className="flex justify-between text-[10px] text-text-subtle">
                  <span>30 Years (Base)</span>
                  <span>40 Years (Recommended)</span>
                  <span>50 Years (Max)</span>
                </div>
              </div>

              {/* Slider 2: Federal Royalty Rate */}
              <div className="space-y-2">
                <div className="flex justify-between text-xs font-mono">
                  <span className="text-text-muted">Federal Royalty (% Gross Rev):</span>
                  <span className="text-aurora font-bold">{federalRoyaltyPct.toFixed(1)}%</span>
                </div>
                <input
                  type="range"
                  min="5"
                  max="15"
                  step="0.5"
                  value={federalRoyaltyPct}
                  onChange={(e) => setFederalRoyaltyPct(Number(e.target.value))}
                  className="w-full accent-aurora h-1.5 bg-surface rounded-lg cursor-pointer"
                />
                <div className="flex justify-between text-[10px] text-text-subtle">
                  <span>5.0% (Light)</span>
                  <span>10.0% (Optimal)</span>
                  <span>15.0% (Max)</span>
                </div>
              </div>

              {/* Slider 3: Pension Equity Share */}
              <div className="space-y-2">
                <div className="flex justify-between text-xs font-mono">
                  <span className="text-text-muted">Pension Equity Share (% Capex):</span>
                  <span className="text-primary font-bold">{pensionEquityPct}%</span>
                </div>
                <input
                  type="range"
                  min="30"
                  max="70"
                  step="5"
                  value={pensionEquityPct}
                  onChange={(e) => {
                    const val = Number(e.target.value);
                    setPensionEquityPct(val);
                    setCommercialDebtPct(Math.max(25, 95 - val));
                  }}
                  className="w-full accent-primary h-1.5 bg-surface rounded-lg cursor-pointer"
                />
                <div className="flex justify-between text-[10px] text-text-subtle">
                  <span>30% Equity</span>
                  <span>50% Balanced</span>
                  <span>70% Pension Majority</span>
                </div>
              </div>

              {/* Slider 4: Passenger Traffic CAGR */}
              <div className="space-y-2">
                <div className="flex justify-between text-xs font-mono">
                  <span className="text-text-muted">Passenger Volume CAGR:</span>
                  <span className="text-cyan-400 font-bold">{passengerCagrPct.toFixed(1)}% / yr</span>
                </div>
                <input
                  type="range"
                  min="1.5"
                  max="4.5"
                  step="0.1"
                  value={passengerCagrPct}
                  onChange={(e) => setPassengerCagrPct(Number(e.target.value))}
                  className="w-full accent-cyan-400 h-1.5 bg-surface rounded-lg cursor-pointer"
                />
                <div className="flex justify-between text-[10px] text-text-subtle">
                  <span>1.5% (Low Growth)</span>
                  <span>2.8% (Historical)</span>
                  <span>4.5% (High Hub)</span>
                </div>
              </div>
            </div>

            {/* Airport Hub Selector & Deep Dive */}
            <div className="lg:col-span-2 rounded-2xl border border-borderSubtle bg-card p-6 space-y-6">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-borderSubtle pb-4">
                <div>
                  <h3 className="font-bold text-text-main text-base">National Airport Gateway Hubs</h3>
                  <p className="text-xs text-text-muted">Select an international hub to inspect concession scope</p>
                </div>

                <div className="flex flex-wrap gap-1.5">
                  {CANONICAL_AIRPORTS.map((hub) => (
                    <button
                      key={hub.code}
                      type="button"
                      onClick={() => setSelectedHub(hub.code)}
                      className={`px-3 py-1.5 rounded-lg text-xs font-mono font-bold transition-all ${
                        selectedHub === hub.code
                          ? "bg-gold text-surfaceDark shadow-md"
                          : "bg-surface hover:bg-card text-text-subtle border border-borderSubtle"
                      }`}
                    >
                      {hub.code} ({hub.province})
                    </button>
                  ))}
                </div>
              </div>

              {/* Selected Hub Details */}
              <div className="space-y-5">
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-xl border border-borderSubtle bg-surface/50 p-4">
                  <div>
                    <div className="text-xs font-mono text-gold font-bold">AIRPORT CONCESSION PROFILE</div>
                    <div className="text-lg font-bold text-text-main">{activeHubDetails.name}</div>
                    <div className="text-xs text-text-muted mt-0.5">
                      Baseline: {(activeHubDetails.paxBaseline / 1e6).toFixed(1)}M passengers/yr → Target 2050: {(activeHubDetails.pax2050 / 1e6).toFixed(1)}M passengers/yr
                    </div>
                  </div>
                  <div className="text-left sm:text-right">
                    <div className="text-xs font-mono text-text-subtle uppercase">Target Concession Capex</div>
                    <div className="text-2xl font-black text-gold font-mono">{formatCad(activeHubDetails.targetCapexCad)}</div>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2 rounded-xl border border-borderSubtle bg-surface/30 p-4">
                    <div className="text-xs font-mono font-bold text-aurora uppercase flex items-center gap-1.5">
                      <Building2 className="h-3.5 w-3.5" />
                      Priority Modernization Projects
                    </div>
                    <ul className="space-y-2 text-xs text-text-muted">
                      {activeHubDetails.keyProjects.map((p, i) => (
                        <li key={i} className="flex items-start gap-2">
                          <CheckCircle2 className="h-3.5 w-3.5 text-aurora shrink-0 mt-0.5" />
                          <span>{p}</span>
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="space-y-2 rounded-xl border border-borderSubtle bg-surface/30 p-4">
                    <div className="text-xs font-mono font-bold text-cyan-400 uppercase flex items-center gap-1.5">
                      <Leaf className="h-3.5 w-3.5" />
                      Decarbonization & Intermodal Scope
                    </div>
                    <ul className="space-y-2 text-xs text-text-muted">
                      {activeHubDetails.decarbInitiatives.map((d, i) => (
                        <li key={i} className="flex items-start gap-2">
                          <Zap className="h-3.5 w-3.5 text-cyan-400 shrink-0 mt-0.5" />
                          <span>{d}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>

                {/* Target Maple 8 Co-Investors */}
                <div className="rounded-xl border border-primary/30 bg-primary/5 p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div>
                    <div className="text-xs font-mono font-bold text-primary">MAPLE 8 CO-INVESTMENT SYNDICATE</div>
                    <div className="text-xs text-text-muted">Lead institutional equity consortium for {activeHubDetails.code}</div>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {activeHubDetails.targetPensions.map((p, idx) => (
                      <span key={idx} className="rounded-lg border border-primary/40 bg-surface px-2.5 py-1 text-xs font-mono font-semibold text-text-main">
                        {p}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Maple 8 Allocation Table */}
          <div className="rounded-2xl border border-borderSubtle bg-card p-6 space-y-4">
            <div className="flex items-center justify-between border-b border-borderSubtle pb-4">
              <div>
                <h3 className="font-bold text-text-main text-base">Canadian Pension Fund (Maple 8) Syndication Tranches</h3>
                <p className="text-xs text-text-muted">Illustrative equity allocation of the {formatCad(concessionSim.pensionEquity)} institutional tranche</p>
              </div>
              <span className="font-mono text-xs px-2.5 py-1 rounded bg-surface border border-borderSubtle text-text-subtle">
                CANADIAN DOMESTIC CAPITAL RETENTION
              </span>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead>
                  <tr className="border-b border-borderSubtle text-text-subtle font-mono uppercase text-[10px]">
                    <th className="py-2.5 px-3">Institutional Investor</th>
                    <th className="py-2.5 px-3">Syndication Share</th>
                    <th className="py-2.5 px-3 text-right">Committed Equity</th>
                    <th className="py-2.5 px-3 text-right">Target Asset Footprint</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-borderSubtle/50 font-mono">
                  {concessionSim.maple8.map((m, idx) => (
                    <tr key={idx} className="hover:bg-surface/50 transition-colors">
                      <td className="py-3 px-3 font-semibold text-text-main font-sans flex items-center gap-2">
                        <span className="h-2 w-2 rounded-full bg-gold" />
                        {m.name}
                      </td>
                      <td className="py-3 px-3 text-text-muted">{m.pct}%</td>
                      <td className="py-3 px-3 text-right font-bold text-gold">{formatCad(m.allocation)}</td>
                      <td className="py-3 px-3 text-right text-text-subtle font-sans">YYZ, YVR, YUL, YYC, YEG Concessions</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* TAB 2: GFANZ / SFAC TRANSITION TAXONOMY ENGINE */}
      {activeTab === "transition_taxonomy" && (
        <div className="space-y-8 animate-fadeIn">
          {/* Taxonomy Architecture Header */}
          <div className="relative overflow-hidden rounded-2xl border border-aurora/30 bg-gradient-to-br from-card via-surface to-[#061611] p-6 shadow-2xl">
            <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
              <div className="space-y-2 max-w-3xl">
                <div className="inline-flex items-center gap-2 rounded-full border border-aurora/40 bg-aurora/10 px-3 py-1 text-xs font-mono font-bold text-aurora">
                  <Leaf className="h-3.5 w-3.5" />
                  GFANZ & SFAC ALIGNED CANADIAN TRANSITION TAXONOMY
                </div>
                <h2 className="text-2xl sm:text-3xl font-black text-text-main">
                  Transition Finance <span className="text-aurora">Credibility Engine</span>
                </h2>
                <p className="text-sm text-text-muted leading-relaxed">
                  Engineered to Mark Carney’s transition finance doctrine: separating credible decarbonization from greenwashing. Classifies capital allocations into Green, Transition (hard-to-abate with no lock-in), and Enabling assets while calculating Marginal Abatement Cost (MAC) and Scope 1/2/3 credibility scores.
                </p>
              </div>

              <div className="flex flex-col items-end justify-center rounded-xl border border-borderSubtle bg-surface/80 p-4 text-right">
                <div className="text-xs font-mono text-text-subtle uppercase">Taxonomy Alignment Standard</div>
                <div className="text-xl font-bold text-aurora font-mono">GFANZ Net-Zero 1.5°C</div>
                <div className="text-[11px] text-text-subtle font-mono mt-1">SFAC Canadian Green & Transition Spec</div>
              </div>
            </div>

            {/* Taxonomy Pillars */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mt-6 pt-6 border-t border-borderSubtle/60">
              <div className="rounded-xl border border-emerald-500/30 bg-emerald-950/20 p-4">
                <div className="text-xs font-mono font-bold text-emerald-400 uppercase">CATEGORY: GREEN</div>
                <div className="text-sm font-bold text-text-main mt-1">Zero / Near-Zero Emissions</div>
                <div className="text-xs text-text-muted mt-1">Nuclear (SMRs), Hydro, Solar/Wind, Battery Storage. Low baseline emissions by design.</div>
              </div>

              <div className="rounded-xl border border-amber-500/30 bg-amber-950/20 p-4">
                <div className="text-xs font-mono font-bold text-amber-400 uppercase">CATEGORY: TRANSITION</div>
                <div className="text-sm font-bold text-text-main mt-1">Decarbonizing Hard-to-Abate</div>
                <div className="text-xs text-text-muted mt-1">Industrial CCUS (Crawford), EAF steel, clean fuels. Strict &quot;No Lock-In&quot; rule.</div>
              </div>

              <div className="rounded-xl border border-cyan-500/30 bg-cyan-950/20 p-4">
                <div className="text-xs font-mono font-bold text-cyan-400 uppercase">CATEGORY: ENABLING</div>
                <div className="text-sm font-bold text-text-main mt-1">Economy-Wide Decarbonization</div>
                <div className="text-xs text-text-muted mt-1">Grid interties, critical minerals refining, airport intermodal rail connectors.</div>
              </div>

              <div className="rounded-xl border border-rose-500/30 bg-rose-950/20 p-4">
                <div className="text-xs font-mono font-bold text-rose-400 uppercase">CATEGORY: PHASE-OUT</div>
                <div className="text-sm font-bold text-text-main mt-1">Managed Early Decommissioning</div>
                <div className="text-xs text-text-muted mt-1">Stranded asset retirement with binding emissions covenants and worker transition.</div>
              </div>
            </div>
          </div>

          {/* Project Benchmark Deep-Dive */}
          <div className="rounded-2xl border border-borderSubtle bg-card p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-borderSubtle pb-4">
              <div>
                <h3 className="font-bold text-text-main text-base">GFANZ Asset Portfolio Assessment</h3>
                <p className="text-xs text-text-muted">Evaluating major Canadian capital builds under Transition Finance rules</p>
              </div>

              <div className="flex flex-wrap gap-1.5">
                {TRANSITION_BENCHMARKS.map((p) => (
                  <button
                    key={p.id}
                    type="button"
                    onClick={() => setSelectedTaxonomyProject(p.id)}
                    className={`px-3 py-1.5 rounded-lg text-xs font-mono font-bold transition-all ${
                      selectedTaxonomyProject === p.id
                        ? "bg-aurora text-surfaceDark shadow-md"
                        : "bg-surface hover:bg-card text-text-subtle border border-borderSubtle"
                    }`}
                  >
                    {p.name.split("—")[0]}
                  </button>
                ))}
              </div>
            </div>

            {/* Selected Project Details */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 space-y-4">
                <div className="rounded-xl border border-borderSubtle bg-surface/50 p-5 space-y-3">
                  <div className="flex items-center justify-between">
                    <span className={`px-2.5 py-0.5 rounded text-xs font-mono font-bold ${
                      activeTaxonomyDetails.category === "GREEN"
                        ? "bg-emerald-950 text-emerald-400 border border-emerald-800"
                        : activeTaxonomyDetails.category === "TRANSITION"
                        ? "bg-amber-950 text-amber-400 border border-amber-800"
                        : "bg-cyan-950 text-cyan-400 border border-cyan-800"
                    }`}>
                      TAXONOMY: {activeTaxonomyDetails.category}
                    </span>
                    <span className="text-xs font-mono text-text-subtle">
                      GFANZ Pillar: <span className="text-text-main font-semibold">{activeTaxonomyDetails.pillar}</span>
                    </span>
                  </div>

                  <div className="text-xl font-bold text-text-main">{activeTaxonomyDetails.name}</div>
                  <div className="text-xs text-text-muted">Sector: {activeTaxonomyDetails.sector} · Capex: {formatCad(activeTaxonomyDetails.capexCad)}</div>

                  <div className="grid grid-cols-3 gap-3 pt-3 border-t border-borderSubtle">
                    <div>
                      <div className="text-[10px] font-mono text-text-subtle uppercase">Annual Abatement</div>
                      <div className="text-sm font-bold text-aurora font-mono">{activeTaxonomyDetails.annualAbatementTpy.toLocaleString()} tCO2e/yr</div>
                    </div>
                    <div>
                      <div className="text-[10px] font-mono text-text-subtle uppercase">Lifetime Abatement</div>
                      <div className="text-sm font-bold text-text-main font-mono">{(activeTaxonomyDetails.lifetimeAbatementTonnes / 1e6).toFixed(1)}M tCO2e</div>
                    </div>
                    <div>
                      <div className="text-[10px] font-mono text-text-subtle uppercase">Marginal Abatement Cost</div>
                      <div className="text-sm font-bold text-gold font-mono">${activeTaxonomyDetails.macCad.toFixed(2)} / tCO2e</div>
                    </div>
                  </div>
                </div>

                <div className="rounded-xl border border-borderSubtle bg-surface/30 p-4 space-y-2">
                  <div className="text-xs font-mono font-bold text-text-main uppercase">Transition Safeguards & Verification</div>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs text-text-muted">
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-aurora shrink-0" />
                      <span>No Lock-In Condition Verified (no carbon lock-in)</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-aurora shrink-0" />
                      <span>2030 Interim 1.5°C Science-Based Targets Disclosed</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-aurora shrink-0" />
                      <span>Third-Party Verification (CSA / IFRS S2 Standard)</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-aurora shrink-0" />
                      <span>Upstream & Downstream Scope 3 Disclosed</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Credibility Score Dial */}
              <div className="rounded-xl border border-aurora/30 bg-surface/60 p-5 flex flex-col items-center justify-center text-center space-y-3">
                <div className="text-xs font-mono uppercase text-aurora font-bold">Transition Credibility Score</div>
                <div className="relative flex items-center justify-center h-28 w-28 rounded-full border-4 border-aurora/40 bg-aurora/5 shadow-xl shadow-aurora/10">
                  <span className="text-3xl font-black text-aurora font-mono">{activeTaxonomyDetails.credibility.toFixed(1)}</span>
                </div>
                <div className="text-xs font-mono text-text-subtle">
                  RATING: <span className="text-white font-bold">{activeTaxonomyDetails.rating}</span> (0–100 Scale)
                </div>
                <p className="text-[11px] text-text-subtle max-w-xs">
                  Evaluates scope 1/2/3 pathways, capex net-zero alignment, lock-in risk mitigation, and independent assurance.
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 3: EXECUTIVE ACTION PLAN & CABINET MC ANNEX */}
      {activeTab === "action_plan" && (
        <div className="space-y-6 animate-fadeIn">
          <div className="rounded-2xl border border-borderSubtle bg-card p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-borderSubtle pb-4">
              <div>
                <h3 className="font-bold text-text-main text-base">Memorandum to Cabinet (MC) Annex: Prime Ministerial Directives</h3>
                <p className="text-xs text-text-muted">Drafted for the Right Honourable Mark Carney & Economic Growth Council</p>
              </div>

              <div className="flex items-center gap-2">
                <button
                  type="button"
                  onClick={handleCopyMemo}
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-borderSubtle bg-surface text-xs font-medium text-text-main hover:border-gold hover:text-gold transition-colors"
                >
                  {copiedMemo ? <Check className="h-3.5 w-3.5 text-aurora" /> : <Copy className="h-3.5 w-3.5" />}
                  <span>{copiedMemo ? "Copied" : "Copy Briefing Memo"}</span>
                </button>
                <button
                  type="button"
                  onClick={() => window.print()}
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-borderSubtle bg-surface text-xs font-medium text-text-main hover:border-primary hover:text-primary transition-colors"
                >
                  <Printer className="h-3.5 w-3.5" />
                  <span>Print Memo</span>
                </button>
              </div>
            </div>

            {/* Structured Text of Memorandum */}
            <div className="rounded-xl border border-borderSubtle bg-surface/40 p-6 font-mono text-xs text-text-muted space-y-4 leading-relaxed">
              <div className="text-text-main font-bold text-sm border-b border-borderSubtle pb-2">
                ORDER IN COUNCIL DIRECTIVE: NATIONAL AIRPORTS CONCESSION & TRANSITION CAPITALIZATION
              </div>

              <p>
                <strong className="text-text-main">1. POLICY DIRECTIVE ON NATIONAL AIRPORTS SYSTEM (NAS) MODERNIZATION:</strong><br />
                Transport Canada is hereby directed to restructure the 1994 National Airports Policy ground leases for Toronto Pearson (YYZ), Vancouver (YVR), Montréal-Trudeau (YUL), Calgary (YYC), and Edmonton (YEG) into 40-year institutional concession agreements. This order unlocks an immediate <span className="text-aurora font-bold">{formatCad(concessionSim.upfrontProceeds)}</span> in upfront federal concession proceeds and mobilizes <span className="text-gold font-bold">$18.0B CAD</span> in private terminal, cargo, and rail capex from Canadian pension funds (Maple 8).
              </p>

              <p>
                <strong className="text-text-main">2. ADOPTION OF CANADIAN TRANSITION FINANCE TAXONOMY:</strong><br />
                Finance Canada, the Canada Infrastructure Bank (CIB), and the Canada Growth Fund (CGF) shall mandate the Sustainable Finance Action Council (SFAC) and GFANZ-aligned Transition Finance Taxonomy as the formal criteria for all federal loan guarantees, Contracts for Difference (CCfDs), and investment tax credits. Clean economy assets shall be classified into Green, Transition (hard-to-abate with no lock-in), and Enabling categories.
              </p>

              <p>
                <strong className="text-text-main">3. DOMESTIC PENSION CROWDING-IN MANDATE:</strong><br />
                Across the indexed $18B NAS program and broader $150B strategic corridor pipeline, federal credit enhancements shall maintain an institutional crowding-in ratio of not less than <span className="text-emerald-400 font-bold">3.8x to 4.2x</span>, guaranteeing long-term domestic ownership by Canadian workers&apos; pension savings (CPPIB, CDPQ, OMERS, Brookfield, PSP, BCI, AIMCo).
              </p>

              <div className="pt-4 border-t border-borderSubtle flex items-center justify-between text-[11px] text-text-subtle">
                <span>SEAL: PRIVY COUNCIL OFFICE // ECONOMIC ADVISORY SECRETARIAT</span>
                <span>SHA-256 AUDIT: 67455d89f0da85a5a07c6c1b595956a257f96f2514570abf8ae9a6c042e3fad1</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
