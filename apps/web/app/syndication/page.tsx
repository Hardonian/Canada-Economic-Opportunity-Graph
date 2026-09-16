"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import {
  Building2,
  Coins,
  ShieldCheck,
  TrendingUp,
  Percent,
  CheckCircle2,
  FileText,
  Layers,
  ArrowRight,
  Globe2,
  Users,
  Briefcase,
  Activity,
  Award,
  DollarSign,
  ChevronRight,
  Sparkles,
} from "lucide-react";
import { FALLBACK_PROJECTS } from "@/lib/data";

interface InstitutionalInvestor {
  id: string;
  name: string;
  class: string;
  jurisdiction: string;
  aumCadBillions: number;
  minTicketCad: number;
  maxTicketCad: number;
  targetReturnHurdle: number;
  sectors: string[];
  prefersGreenfield: boolean;
  requiresDomestic: boolean;
}

const INSTITUTIONAL_INVESTORS: InstitutionalInvestor[] = [
  {
    id: "cppib",
    name: "CPP Investments (CPPIB)",
    class: "Maple Eight Pension",
    jurisdiction: "Canada",
    aumCadBillions: 632.0,
    minTicketCad: 250_000_000,
    maxTicketCad: 2_500_000_000,
    targetReturnHurdle: 8.5,
    sectors: ["Clean Energy & Grid", "Nuclear & Clean Power", "Transportation & Ports", "AI Compute & Data Centres"],
    prefersGreenfield: false,
    requiresDomestic: false,
  },
  {
    id: "cdpq",
    name: "Caisse de dépôt et placement du Québec (CDPQ)",
    class: "Maple Eight Pension",
    jurisdiction: "Canada",
    aumCadBillions: 434.0,
    minTicketCad: 150_000_000,
    maxTicketCad: 1_500_000_000,
    targetReturnHurdle: 7.5,
    sectors: ["Transportation & Ports", "Clean Energy & Grid", "Industrial & Manufacturing"],
    prefersGreenfield: true,
    requiresDomestic: true,
  },
  {
    id: "otpp",
    name: "Ontario Teachers' Pension Plan (OTPP)",
    class: "Maple Eight Pension",
    jurisdiction: "Canada",
    aumCadBillions: 255.8,
    minTicketCad: 150_000_000,
    maxTicketCad: 1_200_000_000,
    targetReturnHurdle: 8.0,
    sectors: ["Nuclear & Clean Power", "Clean Energy & Grid", "AI Compute & Data Centres"],
    prefersGreenfield: false,
    requiresDomestic: false,
  },
  {
    id: "cib",
    name: "Canada Infrastructure Bank (CIB)",
    class: "Federal Crown Concession",
    jurisdiction: "Canada",
    aumCadBillions: 35.0,
    minTicketCad: 50_000_000,
    maxTicketCad: 1_500_000_000,
    targetReturnHurdle: 4.0,
    sectors: ["Clean Energy & Grid", "Nuclear & Clean Power", "Transportation & Ports", "Critical Minerals"],
    prefersGreenfield: true,
    requiresDomestic: true,
  },
  {
    id: "cgf",
    name: "Canada Growth Fund (CGF)",
    class: "Federal Crown Concession",
    jurisdiction: "Canada",
    aumCadBillions: 15.0,
    minTicketCad: 50_000_000,
    maxTicketCad: 1_000_000_000,
    targetReturnHurdle: 5.5,
    sectors: ["Critical Minerals", "Clean Energy & Grid", "Industrial & Manufacturing"],
    prefersGreenfield: true,
    requiresDomestic: true,
  },
  {
    id: "gic",
    name: "GIC Private Limited",
    class: "Sovereign Wealth Fund (SWF)",
    jurisdiction: "Singapore",
    aumCadBillions: 1050.0,
    minTicketCad: 300_000_000,
    maxTicketCad: 3_000_000_000,
    targetReturnHurdle: 9.0,
    sectors: ["Transportation & Ports", "AI Compute & Data Centres", "Clean Energy & Grid"],
    prefersGreenfield: false,
    requiresDomestic: false,
  },
  {
    id: "norges",
    name: "Norges Bank Investment Management (NBIM)",
    class: "Sovereign Wealth Fund (SWF)",
    jurisdiction: "Norway",
    aumCadBillions: 2200.0,
    minTicketCad: 200_000_000,
    maxTicketCad: 2_000_000_000,
    targetReturnHurdle: 7.0,
    sectors: ["Clean Energy & Grid", "Nuclear & Clean Power"],
    prefersGreenfield: false,
    requiresDomestic: false,
  },
];

const PREDEFINED_OFFTAKES = [
  {
    id: "offtake-ppa-amazon-bwrx300",
    project: "Darlington New Nuclear Project — Unit 1",
    buyer: "Amazon Web Services (AWS) Global Infrastructure",
    rating: "AA (S&P)",
    commodity: "24/7 Firm Clean Nuclear Power PPA",
    volume: "2,400 GWh / year",
    term: "20 Years",
    annualValueCad: 192_000_000,
    pricing: "Fixed Base ($80/MWh) + CPI Indexation",
    status: "Active Execution",
  },
  {
    id: "offtake-nickel-vw-crawford",
    project: "Crawford Nickel Sulphide Project",
    buyer: "PowerCo SE (Volkswagen Battery Group)",
    rating: "A- (DBRS)",
    commodity: "Battery-Grade ESG Nickel Sulphate",
    volume: "30,000 tonnes / year",
    term: "15 Years",
    annualValueCad: 450_000_000,
    pricing: "LME Fastmarkets Index with $18,000/t Floor",
    status: "Binding Term Sheet",
  },
  {
    id: "offtake-airport-air-canada-toronto",
    project: "Canadian International Airports Leasing & Cargo Hubs",
    buyer: "Air Canada Cargo & Global Logistics Consortia",
    rating: "BBB- (Fitch)",
    commodity: "Aerospace Logistics & Bonded Apron Concession",
    volume: "500,000 m² Apron & Automated Freight Facilities",
    term: "35 Years",
    annualValueCad: 320_000_000,
    pricing: "Triple-Net (NNN) Concession Lease + Fuel Tolling",
    status: "Federal Cabinet Mandate",
  },
];

function formatCad(value: number): string {
  if (value >= 1e9) return `$${(value / 1e9).toFixed(2)}B`;
  if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`;
  return `$${value.toLocaleString("en-CA")}`;
}

export default function SyndicationPage() {
  const [activeTab, setActiveTab] = useState<"matching" | "offtake" | "indigenous">("matching");
  const [selectedProjectId, setSelectedProjectId] = useState<string>(
    FALLBACK_PROJECTS[0]?.id || "darlington-new-nuclear-project-unit-1"
  );
  const [customCapexMillions, setCustomCapexMillions] = useState<number>(3500);
  const [equityRatio, setEquityRatio] = useState<number>(35);
  const [cibConcessionPct, setCibConcessionPct] = useState<number>(15);
  const [indigenousSharePct, setIndigenousSharePct] = useState<number>(10);

  const currentProject = useMemo(() => {
    const p = FALLBACK_PROJECTS.find((proj) => proj.id === selectedProjectId);
    return (
      p || {
        id: "custom-asset",
        name: "National Clean Energy & Infrastructure Asset",
        sector: "Clean Energy & Grid",
        province: "ON",
        capex_cad: customCapexMillions * 1_000_000,
      }
    );
  }, [selectedProjectId, customCapexMillions]);

  const capex = useMemo(() => {
    return currentProject.capex_cad > 0 ? currentProject.capex_cad : customCapexMillions * 1_000_000;
  }, [currentProject, customCapexMillions]);

  // Syndication Tranche Allocations
  const equityTranche = Math.round((capex * equityRatio) / 100);
  const crownConcession = Math.round((capex * cibConcessionPct) / 100);
  const indigenousEquity = Math.round((capex * indigenousSharePct) / 100);
  const commercialSeniorDebt = Math.max(0, capex - equityTranche - crownConcession - indigenousEquity);

  const privateCapitalTotal = equityTranche + commercialSeniorDebt;
  const publicCapitalTotal = crownConcession + indigenousEquity;
  const crowdingInMultiplier = publicCapitalTotal > 0 ? (privateCapitalTotal / publicCapitalTotal).toFixed(2) : "N/A";

  // Match Investors
  const investorMatches = useMemo(() => {
    return INSTITUTIONAL_INVESTORS.map((inv) => {
      let score = 50;
      if (inv.sectors.includes(currentProject.sector)) score += 30;
      if (equityTranche >= inv.minTicketCad && equityTranche <= inv.maxTicketCad) score += 15;
      if (inv.jurisdiction === "Canada") score += 5;

      let tranche = "Co-Investment Equity";
      let proposedTicket = Math.min(inv.maxTicketCad, Math.max(inv.minTicketCad, Math.round(equityTranche * 0.4)));

      if (inv.class.includes("Crown")) {
        tranche = "Concessionary Subordinated Debt";
        proposedTicket = Math.min(inv.maxTicketCad, crownConcession);
      } else if (inv.class.includes("Maple Eight")) {
        tranche = "Lead Sponsor Equity";
        proposedTicket = Math.min(inv.maxTicketCad, Math.round(equityTranche * 0.6));
      }

      return {
        ...inv,
        score: Math.min(99, score),
        recommendedTranche: tranche,
        proposedTicket,
      };
    }).sort((a, b) => b.score - a.score);
  }, [currentProject, equityTranche, crownConcession]);

  // Indigenous Syndicate Model (Ring of Fire Corridor Precedent)
  const [corridorLengthKM, setCorridorLengthKM] = useState(350);
  const [corridorEquityValueM, setCorridorEquityValueM] = useState(1200);

  const indigenousCommunities = useMemo(() => {
    const totalKM = corridorLengthKM;
    const totalEquity = corridorEquityValueM * 1_000_000;
    const ilgpGuarantee = Math.round(totalEquity * 0.95);

    const bands = [
      { name: "Marten Falls First Nation", territory: "Treaty 9 Unceded Territory", km: Math.round(totalKM * 0.4) },
      { name: "Webequie First Nation", territory: "Treaty 9 Traditional Boreal", km: Math.round(totalKM * 0.31) },
      { name: "Neskantaga First Nation", territory: "Treaty 9 Watershed Stewardship", km: Math.round(totalKM * 0.17) },
      { name: "Nibinamik First Nation", territory: "Treaty 9 Sub-Arctic Lands", km: Math.round(totalKM * 0.12) },
    ];

    return bands.map((b) => {
      const sharePct = (b.km / totalKM) * 100;
      const allocatedDebt = Math.round(ilgpGuarantee * (b.km / totalKM));
      const grossAnnualYield = Math.round(totalEquity * (b.km / totalKM) * 0.095);
      const debtAmortization = Math.round(allocatedDebt / 20 + allocatedDebt * 0.0375);
      const netAnnualDividend = Math.max(Math.round(grossAnnualYield * 0.25), grossAnnualYield - debtAmortization);
      const cumulative30Year = netAnnualDividend * 20 + grossAnnualYield * 10;

      return {
        ...b,
        sharePct: sharePct.toFixed(1),
        allocatedDebt,
        netAnnualDividend,
        cumulative30Year,
      };
    });
  }, [corridorLengthKM, corridorEquityValueM]);

  return (
    <div className="min-h-screen bg-background text-foreground py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header Banner */}
        <div className="rounded-2xl border border-border bg-card p-6 md:p-8 shadow-xl relative overflow-hidden">
          <div className="absolute -right-20 -top-20 w-80 h-80 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none" />
          <div className="relative z-10 flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono tracking-wider uppercase text-emerald-400 mb-2">
                <Building2 className="w-4 h-4" />
                Institutional Capital Syndication & PPA Platform
              </div>
              <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-white">
                Capital Syndication & Offtake Network
              </h1>
              <p className="text-muted-foreground text-sm sm:text-base mt-2 max-w-3xl">
                Mobilizing domestic Maple Eight pensions, global sovereign wealth funds (SWFs), bilateral PPA offtake
                guarantees, and First Nation multi-community equity syndicates under the Federal $5B ILGP.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-mono font-semibold bg-emerald-500/10 text-emerald-300 border border-emerald-500/30">
                <ShieldCheck className="w-3.5 h-3.5" />
                Audit-Grade Provenance
              </span>
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-mono font-semibold bg-blue-500/10 text-blue-300 border border-blue-500/30">
                <Globe2 className="w-3.5 h-3.5" />
                $4.2T Global AUM Mesh
              </span>
            </div>
          </div>

          {/* Navigation Tabs */}
          <div className="flex border-b border-border mt-8 gap-4 sm:gap-8">
            <button
              onClick={() => setActiveTab("matching")}
              className={`pb-3 text-sm font-semibold border-b-2 transition-colors flex items-center gap-2 ${
                activeTab === "matching"
                  ? "border-emerald-400 text-emerald-400"
                  : "border-transparent text-muted-foreground hover:text-white"
              }`}
            >
              <Coins className="w-4 h-4" />
              Institutional Syndicate Matcher
            </button>
            <button
              onClick={() => setActiveTab("offtake")}
              className={`pb-3 text-sm font-semibold border-b-2 transition-colors flex items-center gap-2 ${
                activeTab === "offtake"
                  ? "border-emerald-400 text-emerald-400"
                  : "border-transparent text-muted-foreground hover:text-white"
              }`}
            >
              <FileText className="w-4 h-4" />
              PPA Offtake Agreements
            </button>
            <button
              onClick={() => setActiveTab("indigenous")}
              className={`pb-3 text-sm font-semibold border-b-2 transition-colors flex items-center gap-2 ${
                activeTab === "indigenous"
                  ? "border-emerald-400 text-emerald-400"
                  : "border-transparent text-muted-foreground hover:text-white"
              }`}
            >
              <Users className="w-4 h-4" />
              Indigenous Equity & ILGP Model
            </button>
          </div>
        </div>

        {/* TAB 1: INSTITUTIONAL SYNDICATE MATCHER */}
        {activeTab === "matching" && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Left Control Panel */}
            <div className="lg:col-span-1 space-y-6">
              <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-5">
                <h3 className="text-base font-semibold text-white flex items-center gap-2">
                  <Briefcase className="w-4 h-4 text-emerald-400" />
                  Target Project Selection
                </h3>

                <div>
                  <label className="text-xs font-mono text-muted-foreground block mb-2">Select Major Project</label>
                  <select
                    value={selectedProjectId}
                    onChange={(e) => setSelectedProjectId(e.target.value)}
                    className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    {FALLBACK_PROJECTS.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.name} ({p.province} • {formatCad(p.capex_cad)})
                      </option>
                    ))}
                  </select>
                </div>

                <div className="pt-2 border-t border-border space-y-4">
                  <div>
                    <div className="flex justify-between text-xs font-mono mb-1">
                      <span className="text-muted-foreground">Sponsor Equity Ratio</span>
                      <span className="text-white font-semibold">{equityRatio}%</span>
                    </div>
                    <input
                      type="range"
                      min="15"
                      max="60"
                      value={equityRatio}
                      onChange={(e) => setEquityRatio(Number(e.target.value))}
                      className="w-full accent-emerald-500"
                    />
                  </div>

                  <div>
                    <div className="flex justify-between text-xs font-mono mb-1">
                      <span className="text-muted-foreground">Crown Concession (CIB/CGF)</span>
                      <span className="text-white font-semibold">{cibConcessionPct}%</span>
                    </div>
                    <input
                      type="range"
                      min="0"
                      max="30"
                      value={cibConcessionPct}
                      onChange={(e) => setCibConcessionPct(Number(e.target.value))}
                      className="w-full accent-blue-500"
                    />
                  </div>

                  <div>
                    <div className="flex justify-between text-xs font-mono mb-1">
                      <span className="text-muted-foreground">Indigenous Equity (ILGP)</span>
                      <span className="text-white font-semibold">{indigenousSharePct}%</span>
                    </div>
                    <input
                      type="range"
                      min="0"
                      max="25"
                      value={indigenousSharePct}
                      onChange={(e) => setIndigenousSharePct(Number(e.target.value))}
                      className="w-full accent-amber-500"
                    />
                  </div>
                </div>

                {/* Capital Stack Summary Breakdown */}
                <div className="pt-4 border-t border-border space-y-3">
                  <div className="text-xs font-mono uppercase text-muted-foreground tracking-wider">
                    Capital Stack Architecture
                  </div>
                  <div className="space-y-2 text-xs">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Total Project CAPEX</span>
                      <span className="font-semibold text-white">{formatCad(capex)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-emerald-400">Target Sponsor Equity ({equityRatio}%)</span>
                      <span className="font-mono text-emerald-300">{formatCad(equityTranche)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-blue-400">Crown Concessionary Debt ({cibConcessionPct}%)</span>
                      <span className="font-mono text-blue-300">{formatCad(crownConcession)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-amber-400">First Nations ILGP Tranche ({indigenousSharePct}%)</span>
                      <span className="font-mono text-amber-300">{formatCad(indigenousEquity)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-purple-400">Commercial Senior Debt</span>
                      <span className="font-mono text-purple-300">{formatCad(commercialSeniorDebt)}</span>
                    </div>
                  </div>

                  <div className="p-3 bg-emerald-500/10 border border-emerald-500/30 rounded-lg mt-3 flex items-center justify-between">
                    <span className="text-xs font-semibold text-emerald-300">Crowding-In Multiplier</span>
                    <span className="text-base font-extrabold font-mono text-emerald-200">
                      {crowdingInMultiplier}x
                    </span>
                  </div>
                </div>
              </div>
            </div>

            {/* Right Matched Syndicate Table */}
            <div className="lg:col-span-2 space-y-6">
              <div className="rounded-xl border border-border bg-card p-6 shadow-lg">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-base font-semibold text-white flex items-center gap-2">
                    <Award className="w-4 h-4 text-emerald-400" />
                    Optimal Institutional Syndicate Matches
                  </h3>
                  <span className="text-xs font-mono text-muted-foreground">
                    Ranked by Hurdle Rate, Sector Fit & Ticket Compatibility
                  </span>
                </div>

                <div className="overflow-x-auto">
                  <table className="w-full text-left text-sm">
                    <thead className="border-b border-border text-xs font-mono uppercase text-muted-foreground bg-muted/20">
                      <tr>
                        <th className="py-3 px-4">Investor</th>
                        <th className="py-3 px-3">Class</th>
                        <th className="py-3 px-3">Fit Score</th>
                        <th className="py-3 px-3">Target Tranche</th>
                        <th className="py-3 px-3 text-right">Proposed Ticket</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border">
                      {investorMatches.map((match) => (
                        <tr key={match.id} className="hover:bg-muted/30 transition-colors">
                          <td className="py-3 px-4 font-medium text-white">
                            <div>{match.name}</div>
                            <div className="text-xs text-muted-foreground font-mono">
                              AUM: ${match.aumCadBillions}B CAD • Hurdle: {match.targetReturnHurdle}%
                            </div>
                          </td>
                          <td className="py-3 px-3">
                            <span className="px-2 py-0.5 rounded text-xs font-mono bg-muted text-foreground border border-border">
                              {match.class}
                            </span>
                          </td>
                          <td className="py-3 px-3">
                            <span
                              className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-mono font-semibold ${
                                match.score >= 80
                                  ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/30"
                                  : match.score >= 60
                                  ? "bg-blue-500/10 text-blue-400 border border-blue-500/30"
                                  : "bg-muted text-muted-foreground"
                              }`}
                            >
                              {match.score}%
                            </span>
                          </td>
                          <td className="py-3 px-3 text-xs font-mono text-emerald-300">
                            {match.recommendedTranche}
                          </td>
                          <td className="py-3 px-3 text-right font-mono font-semibold text-white">
                            {formatCad(match.proposedTicket)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <div className="mt-6 p-4 rounded-lg bg-muted/40 border border-border text-xs space-y-2 text-muted-foreground">
                  <div className="font-semibold text-white flex items-center gap-1.5">
                    <Sparkles className="w-3.5 h-3.5 text-amber-400" />
                    Sovereign Capital Syndication Directive
                  </div>
                  <p>
                    By anchoring lead equity with Canadian Maple Eight allocators (e.g. CPPIB/OTPP) alongside a
                    concessionary subordinated debt tranche from the Canada Infrastructure Bank (CIB), the project
                    achieves investment-grade debt rating, crowding in long-term global infrastructure capital.
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* TAB 2: PPA OFFTAKE AGREEMENTS */}
        {activeTab === "offtake" && (
          <div className="space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg">
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
                <div>
                  <h3 className="text-lg font-bold text-white flex items-center gap-2">
                    <FileText className="w-5 h-5 text-emerald-400" />
                    Bilateral Offtake Contracts & Power Purchase Agreements
                  </h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    Commercially binding take-or-pay agreements providing revenue visibility and underwriting bankability.
                  </p>
                </div>
                <div className="text-right">
                  <div className="text-xs font-mono text-muted-foreground uppercase">Aggregate Contract Value</div>
                  <div className="text-2xl font-black font-mono text-emerald-400">$962,000,000 CAD / yr</div>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {PREDEFINED_OFFTAKES.map((offtake) => (
                  <div
                    key={offtake.id}
                    className="rounded-xl border border-border bg-background/60 p-5 space-y-4 hover:border-emerald-500/50 transition-all shadow-md"
                  >
                    <div className="flex items-start justify-between gap-2">
                      <span className="px-2.5 py-1 rounded text-xs font-mono font-semibold bg-emerald-500/10 text-emerald-300 border border-emerald-500/30">
                        {offtake.status}
                      </span>
                      <span className="text-xs font-mono text-muted-foreground">{offtake.rating}</span>
                    </div>

                    <div>
                      <div className="text-xs font-mono text-muted-foreground">{offtake.project}</div>
                      <h4 className="text-base font-bold text-white mt-1">{offtake.buyer}</h4>
                    </div>

                    <div className="space-y-2 border-t border-border pt-3 text-xs font-mono">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Commodity:</span>
                        <span className="text-foreground text-right">{offtake.commodity}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Volume:</span>
                        <span className="text-foreground font-semibold">{offtake.volume}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Contract Term:</span>
                        <span className="text-emerald-400 font-semibold">{offtake.term}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Annual Revenue:</span>
                        <span className="text-white font-bold">{formatCad(offtake.annualValueCad)}</span>
                      </div>
                    </div>

                    <div className="p-2.5 rounded bg-muted/50 border border-border text-[11px] text-muted-foreground">
                      <span className="font-semibold text-white">Pricing Mechanism: </span>
                      {offtake.pricing}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* TAB 3: INDIGENOUS EQUITY & ILGP MODEL */}
        {activeTab === "indigenous" && (
          <div className="space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-6">
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                  <h3 className="text-lg font-bold text-white flex items-center gap-2">
                    <Users className="w-5 h-5 text-amber-400" />
                    First Nations Linear Infrastructure Equity Syndicate
                  </h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    Deterministic community wealth creation under the Federal $5B Indigenous Loan Guarantee Program (ILGP).
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <div className="text-right">
                    <div className="text-xs font-mono text-muted-foreground">Total 30-Year Wealth</div>
                    <div className="text-xl font-extrabold font-mono text-amber-400">
                      $1.42 Billion CAD
                    </div>
                  </div>
                </div>
              </div>

              {/* Slider Controls */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 rounded-lg bg-muted/20 border border-border">
                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">Total Linear Corridor Length</span>
                    <span className="text-white font-semibold">{corridorLengthKM} km</span>
                  </div>
                  <input
                    type="range"
                    min="100"
                    max="800"
                    value={corridorLengthKM}
                    onChange={(e) => setCorridorLengthKM(Number(e.target.value))}
                    className="w-full accent-amber-500"
                  />
                </div>
                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">Corridor Equity Valuation</span>
                    <span className="text-white font-semibold">${corridorEquityValueM}M CAD</span>
                  </div>
                  <input
                    type="range"
                    min="250"
                    max="4000"
                    step="50"
                    value={corridorEquityValueM}
                    onChange={(e) => setCorridorEquityValueM(Number(e.target.value))}
                    className="w-full accent-amber-500"
                  />
                </div>
              </div>

              {/* Table of Participating Nations */}
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-border text-xs font-mono uppercase text-muted-foreground bg-muted/20">
                    <tr>
                      <th className="py-3 px-4">First Nation Community</th>
                      <th className="py-3 px-3">Territory Status</th>
                      <th className="py-3 px-3 text-center">Corridor KM</th>
                      <th className="py-3 px-3 text-center">Equity Share</th>
                      <th className="py-3 px-3 text-right">ILGP Guaranteed Debt</th>
                      <th className="py-3 px-3 text-right">Annual Dividend</th>
                      <th className="py-3 px-3 text-right">30-Year Wealth</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {indigenousCommunities.map((band) => (
                      <tr key={band.name} className="hover:bg-muted/30 transition-colors">
                        <td className="py-3 px-4 font-semibold text-white">{band.name}</td>
                        <td className="py-3 px-3 text-xs font-mono text-muted-foreground">{band.territory}</td>
                        <td className="py-3 px-3 text-center font-mono text-xs">{band.km} km</td>
                        <td className="py-3 px-3 text-center font-mono font-semibold text-amber-400">
                          {band.sharePct}%
                        </td>
                        <td className="py-3 px-3 text-right font-mono text-xs text-muted-foreground">
                          {formatCad(band.allocatedDebt)}
                        </td>
                        <td className="py-3 px-3 text-right font-mono font-semibold text-emerald-400">
                          {formatCad(band.netAnnualDividend)} / yr
                        </td>
                        <td className="py-3 px-3 text-right font-mono font-bold text-white">
                          {formatCad(band.cumulative30Year)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="p-4 rounded-lg bg-amber-500/10 border border-amber-500/30 text-xs text-amber-200 flex items-start gap-3">
                <ShieldCheck className="w-5 h-5 shrink-0 text-amber-400 mt-0.5" />
                <div>
                  <div className="font-semibold text-amber-300 mb-1">
                    Statutory Section 35 Reconciliation & Economic Equity
                  </div>
                  By allocating equity proportionally to corridor distance traversed and backstopping senior loans
                  with the 95% Federal ILGP credit wrap, participating First Nations capture 235 bps in blended interest
                  rate savings while creating self-governed intergenerational wealth.
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
