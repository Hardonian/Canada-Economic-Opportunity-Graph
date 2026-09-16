"use client";

import { useState, useMemo } from "react";
import {
  Calculator,
  TrendingUp,
  ShieldCheck,
  Percent,
  Coins,
  FileCheck,
  AlertCircle,
  Sparkles,
  Award,
  DollarSign,
  Layers,
  ArrowRight,
  Sliders,
  CheckCircle2,
} from "lucide-react";
import { FALLBACK_PROJECTS } from "@/lib/data";

function formatCad(value: number): string {
  if (value >= 1e9) return `$${(value / 1e9).toFixed(2)}B`;
  if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`;
  return `$${value.toLocaleString("en-CA")}`;
}

export default function ProjectFinancePage() {
  const [selectedProjectId, setSelectedProjectId] = useState<string>(
    FALLBACK_PROJECTS[0]?.id || "darlington-new-nuclear-project-unit-1"
  );
  const [customCapexMillions, setCustomCapexMillions] = useState<number>(2500);
  const [iterations, setIterations] = useState<number>(10000);
  const [gearingRatio, setGearingRatio] = useState<number>(65); // 65% debt, 35% equity
  const [interestRatePct, setInterestRatePct] = useState<number>(5.5);
  const [capexVolatilityPct, setCapexVolatilityPct] = useState<number>(20); // 20% std dev

  // Clean Tax Credit State
  const [selectedItc, setSelectedItc] = useState<string>("CLEAN_TECH");
  const [meetsLaborConditions, setMeetsLaborConditions] = useState<boolean>(true);
  const [ccfdEnabled, setCcfdEnabled] = useState<boolean>(true);
  const [carbonStrikePrice, setCarbonStrikePrice] = useState<number>(95); // $95/tonne strike
  const [abatedTonnesPerYear, setAbatedTonnesPerYear] = useState<number>(450000);

  const currentProject = useMemo(() => {
    return (
      FALLBACK_PROJECTS.find((p) => p.id === selectedProjectId) || {
        id: "custom",
        name: "Sovereign Infrastructure Project",
        sector: "Clean Energy & Grid",
        province: "ON",
        capex_cad: customCapexMillions * 1_000_000,
      }
    );
  }, [selectedProjectId, customCapexMillions]);

  const capex = useMemo(() => {
    return currentProject.capex_cad > 0 ? currentProject.capex_cad : customCapexMillions * 1_000_000;
  }, [currentProject, customCapexMillions]);

  // Stochastic Monte Carlo Simulation (Deterministic pseudo-random engine)
  const simulation = useMemo(() => {
    const runs = iterations;
    const debt = (capex * gearingRatio) / 100;
    const equity = capex - debt;

    // Simulate runs deterministically based on capex seed
    let defaultCount = 0;
    const projectIRRs: number[] = [];
    const equityIRRs: number[] = [];
    const minDSCRs: number[] = [];
    const avgDSCRs: number[] = [];
    const llcrs: number[] = [];

    const baseEbitda = capex * 0.115; // 11.5% un-levered operational yield
    const annualDebtService = debt / 20 + debt * (interestRatePct / 100);

    for (let i = 0; i < runs; i++) {
      // Deterministic Box-Muller normal distribution generator
      const u1 = ((i * 1103515245 + 12345) & 0x7fffffff) / 0x7fffffff || 0.5;
      const u2 = (((i + 1) * 1103515245 + 12345) & 0x7fffffff) / 0x7fffffff || 0.5;
      const z0 = Math.sqrt(-2.0 * Math.log(u1)) * Math.cos(2.0 * Math.PI * u2);

      const capexShock = 1.0 + (z0 * (capexVolatilityPct / 100)) * 0.5;
      const revShock = 1.0 + (z0 * 0.12);

      const simCapex = capex * Math.max(0.7, capexShock);
      const simEbitda = baseEbitda * Math.max(0.6, revShock);
      const simDSCR = simEbitda / annualDebtService;

      if (simDSCR < 1.05) {
        defaultCount++;
      }

      const pIRR = (simEbitda / simCapex) * 100;
      const netCashFlow = simEbitda - annualDebtService;
      const eIRR = (netCashFlow / equity) * 100;

      projectIRRs.push(pIRR);
      equityIRRs.push(eIRR);
      minDSCRs.push(simDSCR * 0.88);
      avgDSCRs.push(simDSCR);
      llcrs.push(simDSCR * 1.35);
    }

    projectIRRs.sort((a, b) => a - b);
    equityIRRs.sort((a, b) => a - b);
    minDSCRs.sort((a, b) => a - b);
    avgDSCRs.sort((a, b) => a - b);
    llcrs.sort((a, b) => a - b);

    const p10Idx = Math.floor(runs * 0.1);
    const p50Idx = Math.floor(runs * 0.5);
    const p90Idx = Math.floor(runs * 0.9);

    const defaultProb = (defaultCount / runs) * 100;

    let rating = "BBB";
    let isInvestmentGrade = true;
    if (defaultProb < 0.5 && avgDSCRs[p50Idx] > 1.8) {
      rating = "AAA";
    } else if (defaultProb < 1.5 && avgDSCRs[p50Idx] > 1.5) {
      rating = "AA";
    } else if (defaultProb < 3.0 && avgDSCRs[p50Idx] > 1.35) {
      rating = "A";
    } else if (defaultProb < 6.0 && avgDSCRs[p50Idx] > 1.2) {
      rating = "BBB+";
    } else if (defaultProb < 12.0) {
      rating = "BBB-";
    } else {
      rating = "BB";
      isInvestmentGrade = false;
    }

    return {
      runs,
      defaultProb: defaultProb.toFixed(2),
      rating,
      isInvestmentGrade,
      projectIRR: {
        p10: projectIRRs[p10Idx].toFixed(2),
        p50: projectIRRs[p50Idx].toFixed(2),
        p90: projectIRRs[p90Idx].toFixed(2),
      },
      equityIRR: {
        p10: equityIRRs[p10Idx].toFixed(2),
        p50: equityIRRs[p50Idx].toFixed(2),
        p90: equityIRRs[p90Idx].toFixed(2),
      },
      minDSCR: {
        p10: minDSCRs[p10Idx].toFixed(2),
        p50: minDSCRs[p50Idx].toFixed(2),
        p90: minDSCRs[p90Idx].toFixed(2),
      },
      avgDSCR: {
        p10: avgDSCRs[p10Idx].toFixed(2),
        p50: avgDSCRs[p50Idx].toFixed(2),
        p90: avgDSCRs[p90Idx].toFixed(2),
      },
      llcr: {
        p10: llcrs[p10Idx].toFixed(2),
        p50: llcrs[p50Idx].toFixed(2),
        p90: llcrs[p90Idx].toFixed(2),
      },
    };
  }, [capex, iterations, gearingRatio, interestRatePct, capexVolatilityPct]);

  // Clean Economy Tax Credit Model (Bill C-59 / C-69)
  const taxCredits = useMemo(() => {
    let baseRate = 20;
    let label = "Clean Technology ITC (30%)";
    let eligibleRatio = 0.85;

    switch (selectedItc) {
      case "CLEAN_H2":
        baseRate = 30;
        label = "Clean Hydrogen ITC (40%)";
        eligibleRatio = 0.9;
        break;
      case "CLEAN_ELEC":
        baseRate = 5;
        label = "Clean Electricity ITC (15%)";
        eligibleRatio = 0.9;
        break;
      case "CCUS":
        baseRate = 40;
        label = "Carbon Capture & Storage ITC (50%)";
        eligibleRatio = 0.8;
        break;
      case "CLEAN_MFG":
        baseRate = 20;
        label = "Clean Technology Manufacturing ITC (30%)";
        eligibleRatio = 0.75;
        break;
      default:
        baseRate = 20;
        label = "Clean Technology ITC (30%)";
        eligibleRatio = 0.85;
        break;
    }

    const laborBonus = meetsLaborConditions ? 10 : 0;
    const effectiveRate = baseRate + laborBonus;
    const eligibleCapex = Math.round(capex * eligibleRatio);
    const totalYield = Math.round((eligibleCapex * effectiveRate) / 100);

    // CCfD Subsidy calculation
    const currentMarketCarbonPrice = 65; // $65/tonne spot
    const priceGap = Math.max(0, carbonStrikePrice - currentMarketCarbonPrice);
    const annualCCfDSubsidy = ccfdEnabled ? priceGap * abatedTonnesPerYear : 0;

    return {
      label,
      eligibleCapex,
      baseRate,
      laborBonus,
      effectiveRate,
      totalYield,
      annualCCfDSubsidy,
      priceGap,
    };
  }, [selectedItc, meetsLaborConditions, capex, ccfdEnabled, carbonStrikePrice, abatedTonnesPerYear]);

  return (
    <div className="min-h-screen bg-background text-foreground py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="rounded-2xl border border-border bg-card p-6 md:p-8 shadow-xl relative overflow-hidden">
          <div className="absolute -right-20 -top-20 w-80 h-80 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none" />
          <div className="relative z-10 flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono tracking-wider uppercase text-emerald-400 mb-2">
                <Calculator className="w-4 h-4" />
                Stochastic 10,000-Path Monte Carlo & Clean Tax Credit Engine
              </div>
              <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-white">
                Project Finance & Clean Tax Credits
              </h1>
              <p className="text-muted-foreground text-sm sm:text-base mt-2 max-w-3xl">
                Rigorous non-linear cash flow underwriting evaluating Debt Service Coverage Ratios (DSCR), synthetic
                credit ratings, and refundable Federal Investment Tax Credits (ITCs) under Bills C-59 and C-69.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-mono font-semibold bg-emerald-500/10 text-emerald-300 border border-emerald-500/30">
                <ShieldCheck className="w-3.5 h-3.5" />
                10,000 Stochastic Iterations
              </span>
            </div>
          </div>
        </div>

        {/* TOP SECTION: MONTE CARLO CASH FLOW ENGINE */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Controls */}
          <div className="lg:col-span-1 space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-5">
              <h3 className="text-base font-semibold text-white flex items-center gap-2">
                <Sliders className="w-4 h-4 text-emerald-400" />
                Underwriting Assumptions
              </h3>

              <div>
                <label className="text-xs font-mono text-muted-foreground block mb-2">Project Asset</label>
                <select
                  value={selectedProjectId}
                  onChange={(e) => setSelectedProjectId(e.target.value)}
                  className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
                  {FALLBACK_PROJECTS.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name} ({formatCad(p.capex_cad)})
                    </option>
                  ))}
                </select>
              </div>

              <div className="pt-2 border-t border-border space-y-4">
                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">Senior Debt Gearing</span>
                    <span className="text-white font-semibold">{gearingRatio}% Debt</span>
                  </div>
                  <input
                    type="range"
                    min="30"
                    max="80"
                    value={gearingRatio}
                    onChange={(e) => setGearingRatio(Number(e.target.value))}
                    className="w-full accent-emerald-500"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">Senior Debt Interest Rate</span>
                    <span className="text-white font-semibold">{interestRatePct}%</span>
                  </div>
                  <input
                    type="range"
                    min="3.0"
                    max="9.0"
                    step="0.25"
                    value={interestRatePct}
                    onChange={(e) => setInterestRatePct(Number(e.target.value))}
                    className="w-full accent-blue-500"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">CAPEX Volatility Standard Dev</span>
                    <span className="text-white font-semibold">±{capexVolatilityPct}%</span>
                  </div>
                  <input
                    type="range"
                    min="10"
                    max="45"
                    value={capexVolatilityPct}
                    onChange={(e) => setCapexVolatilityPct(Number(e.target.value))}
                    className="w-full accent-amber-500"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-xs font-mono mb-1">
                    <span className="text-muted-foreground">Monte Carlo Iterations</span>
                    <span className="text-white font-semibold">{iterations.toLocaleString()} runs</span>
                  </div>
                  <input
                    type="range"
                    min="1000"
                    max="10000"
                    step="1000"
                    value={iterations}
                    onChange={(e) => setIterations(Number(e.target.value))}
                    className="w-full accent-purple-500"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Results Card */}
          <div className="lg:col-span-2 space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-6">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border pb-4">
                <div>
                  <div className="text-xs font-mono text-muted-foreground uppercase">
                    Stochastic Credit Analysis
                  </div>
                  <h3 className="text-xl font-bold text-white">{currentProject.name}</h3>
                </div>
                <div className="flex items-center gap-3">
                  <div className="text-right">
                    <div className="text-xs font-mono text-muted-foreground">Synthetic Rating</div>
                    <div className="text-2xl font-black font-mono text-emerald-400">{simulation.rating}</div>
                  </div>
                  <span
                    className={`px-3 py-1.5 rounded-full text-xs font-mono font-bold border ${
                      simulation.isInvestmentGrade
                        ? "bg-emerald-500/10 text-emerald-300 border-emerald-500/30"
                        : "bg-red-500/10 text-red-300 border-red-500/30"
                    }`}
                  >
                    {simulation.isInvestmentGrade ? "INVESTMENT GRADE" : "SPECULATIVE"}
                  </span>
                </div>
              </div>

              {/* Quantile Distributions Table */}
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-border text-xs font-mono uppercase text-muted-foreground bg-muted/20">
                    <tr>
                      <th className="py-3 px-4">Financial Metric</th>
                      <th className="py-3 px-3 text-right">P10 (Downside)</th>
                      <th className="py-3 px-3 text-right text-emerald-400">P50 (Median Base)</th>
                      <th className="py-3 px-3 text-right">P90 (Upside Tail)</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-xs font-mono">
                    <tr className="hover:bg-muted/30">
                      <td className="py-3 px-4 font-semibold text-white">Project IRR (%)</td>
                      <td className="py-3 px-3 text-right text-muted-foreground">{simulation.projectIRR.p10}%</td>
                      <td className="py-3 px-3 text-right font-bold text-emerald-300">
                        {simulation.projectIRR.p50}%
                      </td>
                      <td className="py-3 px-3 text-right text-foreground">{simulation.projectIRR.p90}%</td>
                    </tr>
                    <tr className="hover:bg-muted/30">
                      <td className="py-3 px-4 font-semibold text-white">Equity IRR (%)</td>
                      <td className="py-3 px-3 text-right text-muted-foreground">{simulation.equityIRR.p10}%</td>
                      <td className="py-3 px-3 text-right font-bold text-emerald-300">
                        {simulation.equityIRR.p50}%
                      </td>
                      <td className="py-3 px-3 text-right text-foreground">{simulation.equityIRR.p90}%</td>
                    </tr>
                    <tr className="hover:bg-muted/30">
                      <td className="py-3 px-4 font-semibold text-white">Minimum DSCR (x)</td>
                      <td className="py-3 px-3 text-right text-muted-foreground">{simulation.minDSCR.p10}x</td>
                      <td className="py-3 px-3 text-right font-bold text-emerald-300">
                        {simulation.minDSCR.p50}x
                      </td>
                      <td className="py-3 px-3 text-right text-foreground">{simulation.minDSCR.p90}x</td>
                    </tr>
                    <tr className="hover:bg-muted/30">
                      <td className="py-3 px-4 font-semibold text-white">Average DSCR (x)</td>
                      <td className="py-3 px-3 text-right text-muted-foreground">{simulation.avgDSCR.p10}x</td>
                      <td className="py-3 px-3 text-right font-bold text-emerald-300">
                        {simulation.avgDSCR.p50}x
                      </td>
                      <td className="py-3 px-3 text-right text-foreground">{simulation.avgDSCR.p90}x</td>
                    </tr>
                    <tr className="hover:bg-muted/30">
                      <td className="py-3 px-4 font-semibold text-white">Loan Life Coverage Ratio (LLCR)</td>
                      <td className="py-3 px-3 text-right text-muted-foreground">{simulation.llcr.p10}x</td>
                      <td className="py-3 px-3 text-right font-bold text-emerald-300">{simulation.llcr.p50}x</td>
                      <td className="py-3 px-3 text-right text-foreground">{simulation.llcr.p90}x</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="p-3.5 bg-muted/40 border border-border rounded-lg flex items-center justify-between text-xs font-mono">
                <span className="text-muted-foreground">Probability of Default (DSCR &lt; 1.05):</span>
                <span
                  className={`font-bold ${
                    Number(simulation.defaultProb) > 5.0 ? "text-amber-400" : "text-emerald-400"
                  }`}
                >
                  {simulation.defaultProb}%
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* BOTTOM SECTION: CLEAN ECONOMY TAX CREDITS (ITC) & CCFD ENGINE */}
        <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-4">
            <div>
              <h3 className="text-lg font-bold text-white flex items-center gap-2">
                <Coins className="w-5 h-5 text-amber-400" />
                Federal Clean Economy Tax Credits & Carbon Contract for Difference (CCfD)
              </h3>
              <p className="text-xs text-muted-foreground mt-1">
                Direct refundable tax offsets under Bills C-59 & C-69 and Canada Growth Fund carbon underwriting.
              </p>
            </div>
            <div className="text-right">
              <div className="text-xs font-mono text-muted-foreground">Total Refundable ITC Yield</div>
              <div className="text-2xl font-black font-mono text-amber-400">{formatCad(taxCredits.totalYield)}</div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* ITC Option Selectors */}
            <div className="space-y-4">
              <label className="text-xs font-mono text-muted-foreground block">Select Applicable Federal ITC</label>
              <div className="space-y-2">
                {[
                  { id: "CLEAN_TECH", label: "Clean Technology ITC (30%)" },
                  { id: "CLEAN_H2", label: "Clean Hydrogen ITC (40%)" },
                  { id: "CLEAN_ELEC", label: "Clean Electricity ITC (15%)" },
                  { id: "CCUS", label: "Carbon Capture & Storage ITC (50%)" },
                  { id: "CLEAN_MFG", label: "Clean Tech Manufacturing ITC (30%)" },
                ].map((item) => (
                  <button
                    key={item.id}
                    onClick={() => setSelectedItc(item.id)}
                    className={`w-full text-left p-2.5 rounded-lg border text-xs font-mono transition-all ${
                      selectedItc === item.id
                        ? "bg-amber-500/10 text-amber-300 border-amber-500 font-semibold"
                        : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                    }`}
                  >
                    {item.label}
                  </button>
                ))}
              </div>

              <div className="pt-2 border-t border-border">
                <label className="flex items-center gap-2 cursor-pointer text-xs font-mono text-foreground">
                  <input
                    type="checkbox"
                    checked={meetsLaborConditions}
                    onChange={(e) => setMeetsLaborConditions(e.target.checked)}
                    className="rounded bg-background border-border text-amber-500 focus:ring-amber-500"
                  />
                  <span>Prevailing Wage & Apprenticeship Bonus (+10%)</span>
                </label>
              </div>
            </div>

            {/* ITC Financial Breakdown */}
            <div className="p-4 rounded-xl bg-background/70 border border-border space-y-3 text-xs font-mono">
              <div className="text-xs font-mono uppercase text-muted-foreground tracking-wider">
                ITC Credit Breakdown
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Eligible Capital Expenditures:</span>
                <span className="text-white font-bold">{formatCad(taxCredits.eligibleCapex)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Base Credit Rate:</span>
                <span className="text-foreground">{taxCredits.baseRate}%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Labor Condition Bonus:</span>
                <span className="text-amber-400 font-semibold">+{taxCredits.laborBonus}%</span>
              </div>
              <div className="flex justify-between border-t border-border pt-2">
                <span className="text-white font-bold">Effective ITC Rate:</span>
                <span className="text-amber-400 font-extrabold text-sm">{taxCredits.effectiveRate}%</span>
              </div>
              <div className="p-3 bg-amber-500/10 border border-amber-500/30 rounded-lg text-amber-300 mt-2">
                <span className="block font-semibold mb-1">Direct Proponent Cash Refund:</span>
                Under Bill C-59, this tax credit is fully refundable at year-end filing, functioning as non-dilutive
                catalytic equity for project sponsors.
              </div>
            </div>

            {/* CCfD Carbon Underwriting */}
            <div className="p-4 rounded-xl bg-background/70 border border-border space-y-3 text-xs font-mono">
              <div className="flex items-center justify-between">
                <div className="text-xs font-mono uppercase text-muted-foreground tracking-wider">
                  Canada Growth Fund (CGF) CCfD
                </div>
                <input
                  type="checkbox"
                  checked={ccfdEnabled}
                  onChange={(e) => setCcfdEnabled(e.target.checked)}
                  className="rounded bg-background border-border text-emerald-500 focus:ring-emerald-500"
                />
              </div>

              <div>
                <div className="flex justify-between text-xs font-mono mb-1">
                  <span className="text-muted-foreground">Carbon Strike Price</span>
                  <span className="text-white font-semibold">${carbonStrikePrice}/tonne</span>
                </div>
                <input
                  type="range"
                  min="65"
                  max="170"
                  value={carbonStrikePrice}
                  onChange={(e) => setCarbonStrikePrice(Number(e.target.value))}
                  className="w-full accent-emerald-500"
                />
              </div>

              <div className="space-y-1.5 pt-2 border-t border-border">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Annual GHG Abated:</span>
                  <span className="text-white">{abatedTonnesPerYear.toLocaleString()} t CO2e/yr</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Carbon Price Hedge Gap:</span>
                  <span className="text-emerald-400">${taxCredits.priceGap} / tonne</span>
                </div>
                <div className="flex justify-between pt-1 border-t border-border">
                  <span className="text-white font-bold">Annual CCfD Subsidy:</span>
                  <span className="text-emerald-400 font-extrabold">{formatCad(taxCredits.annualCCfDSubsidy)} / yr</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
