"use client";

import { useState, useMemo } from "react";
import {
  Layers,
  Percent,
  TrendingUp,
  ShieldCheck,
  Building2,
  Landmark,
  Coins,
  ArrowRight,
  Sliders,
  DollarSign,
  AlertCircle,
  HelpCircle,
} from "lucide-react";

export interface TrancheConfig {
  id: string;
  name: string;
  category: "senior_debt" | "cib_concessionary" | "ilgp_indigenous" | "tax_credits" | "subordinated_debt" | "equity";
  sharePct: number;
  interestRatePct: number;
  amortizationYears: number;
  sovereignGuaranteed: boolean;
  color: string;
  badge: string;
}

export const PRESET_CAPITAL_PROFILES = [
  {
    id: "critical_minerals_smr",
    title: "Critical Minerals & SMR Smelter (ON/QC)",
    capexMillions: 3200,
    ebitdaMarginPct: 38,
    tranches: [
      { id: "sr_debt", name: "Commercial Senior Debt (Export Development Canada syndication)", category: "senior_debt", sharePct: 35, interestRatePct: 6.2, amortizationYears: 18, sovereignGuaranteed: false, color: "#3B82F6", badge: "EDC/COMM" },
      { id: "cib_tranche", name: "Canada Infrastructure Bank (CIB Concessionary Tranche)", category: "cib_concessionary", sharePct: 20, interestRatePct: 3.8, amortizationYears: 25, sovereignGuaranteed: true, color: "#10B981", badge: "CIB-GROWTH" },
      { id: "ilgp_tranche", name: "Indigenous Loan Guarantee Program (ILGP 95% Backed)", category: "ilgp_indigenous", sharePct: 15, interestRatePct: 4.1, amortizationYears: 20, sovereignGuaranteed: true, color: "#F59E0B", badge: "ILGP-AAA" },
      { id: "c59_itc", name: "Bill C-59 Clean Tech & Critical Mineral ITC (Federal)", category: "tax_credits", sharePct: 10, interestRatePct: 0.0, amortizationYears: 0, sovereignGuaranteed: true, color: "#8B5CF6", badge: "C-59 ITC" },
      { id: "sponsor_equity", name: "Sponsor Equity & Strategic Consortium", category: "equity", sharePct: 20, interestRatePct: 14.0, amortizationYears: 0, sovereignGuaranteed: false, color: "#EC4899", badge: "SPONSOR" },
    ] as TrancheConfig[],
  },
  {
    id: "clean_hydrogen_ammonia",
    title: "Atlantic Clean Hydrogen & Export Port (NL/NS)",
    capexMillions: 4800,
    ebitdaMarginPct: 42,
    tranches: [
      { id: "sr_debt", name: "International Green Syndicate Senior Facilities", category: "senior_debt", sharePct: 40, interestRatePct: 5.9, amortizationYears: 20, sovereignGuaranteed: false, color: "#3B82F6", badge: "GREEN-SYN" },
      { id: "c59_itc", name: "Bill C-59 Clean Hydrogen ITC (40% Tier 1 Carbon Intensity)", category: "tax_credits", sharePct: 25, interestRatePct: 0.0, amortizationYears: 0, sovereignGuaranteed: true, color: "#8B5CF6", badge: "C-59 40%" },
      { id: "cib_tranche", name: "CIB Clean Energy Infrastructure Facility", category: "cib_concessionary", sharePct: 15, interestRatePct: 3.5, amortizationYears: 22, sovereignGuaranteed: true, color: "#10B981", badge: "CIB-LOW" },
      { id: "sponsor_equity", name: "Consortium Equity & German Offtaker Pre-Paid Equity", category: "equity", sharePct: 20, interestRatePct: 13.5, amortizationYears: 0, sovereignGuaranteed: false, color: "#EC4899", badge: "OFFTAKER" },
    ] as TrancheConfig[],
  },
  {
    id: "arctic_defence_corridor",
    title: "Northern Corridors & Dual-Use Deepwater Terminal (NT/NU)",
    capexMillions: 1850,
    ebitdaMarginPct: 28,
    tranches: [
      { id: "cib_tranche", name: "CIB Arctic & Northern Sovereign Infrastructure", category: "cib_concessionary", sharePct: 45, interestRatePct: 3.2, amortizationYears: 30, sovereignGuaranteed: true, color: "#10B981", badge: "ARCTIC-CIB" },
      { id: "ilgp_tranche", name: "Inuit Development Corporation Co-Ownership Tranche", category: "ilgp_indigenous", sharePct: 25, interestRatePct: 3.9, amortizationYears: 25, sovereignGuaranteed: true, color: "#F59E0B", badge: "INUIT-ILGP" },
      { id: "sr_debt", name: "Commercial & Defence Dual-Use Commercial Tranche", category: "senior_debt", sharePct: 15, interestRatePct: 6.5, amortizationYears: 15, sovereignGuaranteed: false, color: "#3B82F6", badge: "DEFENCE" },
      { id: "sponsor_equity", name: "Crown / Territorial Sponsor Equity", category: "equity", sharePct: 15, interestRatePct: 8.0, amortizationYears: 0, sovereignGuaranteed: true, color: "#EC4899", badge: "CROWN-EQ" },
    ] as TrancheConfig[],
  },
];

export default function CapitalStackBuilder() {
  const [selectedPresetId, setSelectedPresetId] = useState<string>("critical_minerals_smr");
  const [capexMillions, setCapexMillions] = useState<number>(3200);
  const [revenueEstMillions, setRevenueEstMillions] = useState<number>(850);
  const [ebitdaMarginPct, setEbitdaMarginPct] = useState<number>(38);
  const [tranches, setTranches] = useState<TrancheConfig[]>(PRESET_CAPITAL_PROFILES[0].tranches);

  const applyPreset = (presetId: string) => {
    const p = PRESET_CAPITAL_PROFILES.find((x) => x.id === presetId);
    if (!p) return;
    setSelectedPresetId(p.id);
    setCapexMillions(p.capexMillions);
    setEbitdaMarginPct(p.ebitdaMarginPct);
    setRevenueEstMillions(Math.round(p.capexMillions * 0.28));
    setTranches(JSON.parse(JSON.stringify(p.tranches)));
  };

  const handleShareChange = (id: string, newShare: number) => {
    setTranches((prev) =>
      prev.map((t) => (t.id === id ? { ...t, sharePct: Math.max(0, Math.min(100, newShare)) } : t))
    );
  };

  const handleRateChange = (id: string, newRate: number) => {
    setTranches((prev) =>
      prev.map((t) => (t.id === id ? { ...t, interestRatePct: Math.max(0, Math.min(30, newRate)) } : t))
    );
  };

  const metrics = useMemo(() => {
    const totalShare = tranches.reduce((sum, t) => sum + t.sharePct, 0);
    const normalizedCapex = capexMillions * 1_000_000;
    const ebitda = (revenueEstMillions * 1_000_000 * ebitdaMarginPct) / 100;

    let weightedCostSum = 0;
    let totalDebtServiceAnnual = 0;
    let totalSovereignBackedCapex = 0;
    let commercialRateBenchmark = 7.5; // Benchmark unsupported CAD project debt rate

    const breakdown = tranches.map((t) => {
      const trancheCapex = (normalizedCapex * t.sharePct) / (totalShare || 1);
      weightedCostSum += (t.sharePct / (totalShare || 1)) * t.interestRatePct;

      let annualDebtService = 0;
      if (t.amortizationYears > 0 && t.interestRatePct > 0) {
        const r = t.interestRatePct / 100;
        const n = t.amortizationYears;
        // Standard annuity formula: P * (r / (1 - (1+r)^-n))
        annualDebtService = trancheCapex * (r / (1 - Math.pow(1 + r, -n)));
        totalDebtServiceAnnual += annualDebtService;
      } else if (t.category === "senior_debt" || t.category === "subordinated_debt") {
        annualDebtService = trancheCapex * (t.interestRatePct / 100);
        totalDebtServiceAnnual += annualDebtService;
      }

      if (t.sovereignGuaranteed) {
        totalSovereignBackedCapex += trancheCapex;
      }

      return {
        ...t,
        dollarValue: trancheCapex,
        annualService: annualDebtService,
      };
    });

    const wacc = totalShare > 0 ? weightedCostSum : 0;
    const dscr = totalDebtServiceAnnual > 0 ? ebitda / totalDebtServiceAnnual : 99.9;
    const sovereignCoveragePct = normalizedCapex > 0 ? (totalSovereignBackedCapex / normalizedCapex) * 100 : 0;
    
    // Annual interest savings vs standard commercial debt without CIB / ILGP / ITC
    const benchmarkAnnualInterest = (normalizedCapex * 0.7 * (commercialRateBenchmark / 100));
    const actualDebtInterest = breakdown
      .filter((b) => b.category !== "equity" && b.category !== "tax_credits")
      .reduce((acc, curr) => acc + curr.dollarValue * (curr.interestRatePct / 100), 0);
    const annualInterestSavings = Math.max(0, benchmarkAnnualInterest - actualDebtInterest);

    return {
      totalShare,
      breakdown,
      wacc,
      dscr,
      ebitda,
      sovereignCoveragePct,
      annualInterestSavings,
      totalDebtServiceAnnual,
    };
  }, [capexMillions, revenueEstMillions, ebitdaMarginPct, tranches]);

  return (
    <div className="space-y-6 rounded-2xl border border-border bg-card p-6 shadow-2xl">
      {/* Header & Preset Selector */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between border-b border-borderSubtle pb-5">
        <div>
          <div className="flex items-center gap-2">
            <Layers className="h-6 w-6 text-aurora" />
            <h2 className="text-xl font-black tracking-tight text-text-main">
              Sovereign Capital Stack & Financing Waterfall
            </h2>
            <span className="rounded-full border border-primary/40 bg-primary/10 px-2.5 py-0.5 text-xs font-mono font-bold text-aurora">
              PALANTIR FOUNDRY PARITY
            </span>
          </div>
          <p className="mt-1 text-xs text-text-muted">
            Multi-tranche debt & equity financial engineering simulator evaluating CIB concessionary facilities, ILGP Indigenous guarantees, and Bill C-59 ITC capital relief.
          </p>
        </div>

        {/* Preset Buttons */}
        <div className="flex flex-wrap items-center gap-1.5 bg-surface p-1 rounded-xl border border-borderSubtle">
          {PRESET_CAPITAL_PROFILES.map((p) => (
            <button
              key={p.id}
              onClick={() => applyPreset(p.id)}
              className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                selectedPresetId === p.id
                  ? "bg-primary text-black font-bold shadow-sm"
                  : "text-text-muted hover:text-text-main hover:bg-card"
              }`}
            >
              {p.title.split(" (")[0]}
            </button>
          ))}
        </div>
      </div>

      {/* Macro Input Sliders */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between text-xs font-semibold text-text-muted mb-2">
            <span className="flex items-center gap-1.5">
              <DollarSign className="h-3.5 w-3.5 text-aurora" /> Total Project CapEx
            </span>
            <span className="font-mono text-sm font-black text-text-main">${capexMillions.toLocaleString()}M CAD</span>
          </div>
          <input
            type="range"
            min={100}
            max={15000}
            step={50}
            value={capexMillions}
            onChange={(e) => setCapexMillions(Number(e.target.value))}
            className="w-full accent-aurora h-1.5 bg-card rounded-lg cursor-pointer"
          />
        </div>

        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between text-xs font-semibold text-text-muted mb-2">
            <span className="flex items-center gap-1.5">
              <TrendingUp className="h-3.5 w-3.5 text-blue-400" /> Projected Annual Revenue
            </span>
            <span className="font-mono text-sm font-black text-text-main">${revenueEstMillions.toLocaleString()}M CAD</span>
          </div>
          <input
            type="range"
            min={20}
            max={4000}
            step={20}
            value={revenueEstMillions}
            onChange={(e) => setRevenueEstMillions(Number(e.target.value))}
            className="w-full accent-blue-400 h-1.5 bg-card rounded-lg cursor-pointer"
          />
        </div>

        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between text-xs font-semibold text-text-muted mb-2">
            <span className="flex items-center gap-1.5">
              <Percent className="h-3.5 w-3.5 text-emerald-400" /> EBITDA Operating Margin
            </span>
            <span className="font-mono text-sm font-black text-text-main">{ebitdaMarginPct}%</span>
          </div>
          <input
            type="range"
            min={10}
            max={75}
            step={1}
            value={ebitdaMarginPct}
            onChange={(e) => setEbitdaMarginPct(Number(e.target.value))}
            className="w-full accent-emerald-400 h-1.5 bg-card rounded-lg cursor-pointer"
          />
        </div>
      </div>

      {/* Visual Tranche Stack Bar */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-xs font-semibold text-text-muted">
          <span>Capital Stack Tranche Allocation ({metrics.totalShare}% Allocated)</span>
          {metrics.totalShare !== 100 && (
            <span className="flex items-center gap-1 text-amber-400 font-mono text-[11px]">
              <AlertCircle className="h-3 w-3" /> Allocation must equal 100% (currently {metrics.totalShare}%)
            </span>
          )}
        </div>
        <div className="flex h-7 w-full overflow-hidden rounded-xl border border-borderSubtle bg-card shadow-inner">
          {metrics.breakdown.map((t) => (
            <div
              key={t.id}
              style={{ width: `${(t.sharePct / (metrics.totalShare || 1)) * 100}%`, backgroundColor: t.color }}
              className="group relative flex items-center justify-center transition-all duration-300 hover:brightness-110"
              title={`${t.name}: ${t.sharePct}% ($${Math.round(t.dollarValue / 1e6)}M CAD)`}
            >
              {t.sharePct >= 8 && (
                <span className="truncate px-1 font-mono text-[10px] font-bold text-white drop-shadow">
                  {t.badge} {t.sharePct}%
                </span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Tranche Configuration Table */}
      <div className="overflow-x-auto rounded-xl border border-borderSubtle bg-surface">
        <table className="w-full text-left text-xs">
          <thead className="border-b border-borderSubtle bg-card/60 font-mono text-text-muted">
            <tr>
              <th className="py-2.5 px-3">Tranche / Facility</th>
              <th className="py-2.5 px-3">Share %</th>
              <th className="py-2.5 px-3">Tranche CapEx (CAD)</th>
              <th className="py-2.5 px-3">Cost / Rate %</th>
              <th className="py-2.5 px-3">Amortization</th>
              <th className="py-2.5 px-3">Annual Service</th>
              <th className="py-2.5 px-3">Sovereign Protection</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-borderSubtle">
            {metrics.breakdown.map((t) => (
              <tr key={t.id} className="hover:bg-card/40 transition-colors">
                <td className="py-2.5 px-3">
                  <div className="flex items-center gap-2">
                    <span className="h-2.5 w-2.5 rounded-full shrink-0" style={{ backgroundColor: t.color }} />
                    <span className="font-semibold text-text-main">{t.name}</span>
                  </div>
                </td>
                <td className="py-2.5 px-3 font-mono font-bold text-text-main">
                  <div className="flex items-center gap-1.5">
                    <input
                      type="number"
                      min={0}
                      max={100}
                      value={t.sharePct}
                      onChange={(e) => handleShareChange(t.id, Number(e.target.value))}
                      className="w-14 rounded border border-borderSubtle bg-card px-1.5 py-0.5 text-center font-mono text-xs font-bold text-text-main focus:border-primary focus:outline-none"
                    />
                    <span>%</span>
                  </div>
                </td>
                <td className="py-2.5 px-3 font-mono font-bold text-aurora">
                  ${(t.dollarValue / 1e6).toFixed(1)}M
                </td>
                <td className="py-2.5 px-3 font-mono">
                  <div className="flex items-center gap-1.5">
                    <input
                      type="number"
                      step={0.1}
                      min={0}
                      max={30}
                      value={t.interestRatePct}
                      onChange={(e) => handleRateChange(t.id, Number(e.target.value))}
                      className="w-14 rounded border border-borderSubtle bg-card px-1.5 py-0.5 text-center font-mono text-xs font-bold text-text-main focus:border-primary focus:outline-none"
                    />
                    <span>%</span>
                  </div>
                </td>
                <td className="py-2.5 px-3 font-mono text-text-muted">
                  {t.amortizationYears > 0 ? `${t.amortizationYears} yrs` : "N/A (Equity/ITC)"}
                </td>
                <td className="py-2.5 px-3 font-mono text-text-muted">
                  {t.annualService > 0 ? `$${(t.annualService / 1e6).toFixed(1)}M/yr` : "—"}
                </td>
                <td className="py-2.5 px-3">
                  {t.sovereignGuaranteed ? (
                    <span className="inline-flex items-center gap-1 rounded border border-emerald-500/40 bg-emerald-500/10 px-2 py-0.5 font-mono text-[10px] font-bold text-emerald-400">
                      <ShieldCheck className="h-3 w-3" /> CROWN / AAA
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1 rounded border border-borderSubtle bg-card px-2 py-0.5 font-mono text-[10px] text-text-subtle">
                      Commercial Risk
                    </span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* KPI & Valuation Engine Output Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Blended WACC */}
        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-text-muted">Blended WACC</span>
            <span className="rounded bg-primary/10 px-1.5 py-0.5 font-mono text-[10px] font-bold text-aurora">
              HURDLE RATE
            </span>
          </div>
          <div className="mt-2 text-2xl font-black font-mono text-aurora">
            {metrics.wacc.toFixed(2)}%
          </div>
          <p className="mt-1 text-[11px] text-text-subtle">
            Weighted Average Cost of Capital blending senior debt, CIB & equity.
          </p>
        </div>

        {/* DSCR */}
        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-text-muted">Min Debt Service Coverage (DSCR)</span>
            <span className={`rounded px-1.5 py-0.5 font-mono text-[10px] font-bold ${
              metrics.dscr >= 1.4 ? "bg-emerald-500/10 text-emerald-400" : "bg-amber-500/10 text-amber-400"
            }`}>
              {metrics.dscr >= 1.4 ? "BANKABLE" : "TIGHT"}
            </span>
          </div>
          <div className={`mt-2 text-2xl font-black font-mono ${
            metrics.dscr >= 1.4 ? "text-emerald-400" : "text-amber-400"
          }`}>
            {metrics.dscr > 50 ? ">50.0x" : `${metrics.dscr.toFixed(2)}x`}
          </div>
          <p className="mt-1 text-[11px] text-text-subtle">
            Annual EBITDA (${(metrics.ebitda / 1e6).toFixed(0)}M) vs Total Debt Service (${(metrics.totalDebtServiceAnnual / 1e6).toFixed(1)}M).
          </p>
        </div>

        {/* Sovereign Capital Shield */}
        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-text-muted">Sovereign Shield Share</span>
            <span className="rounded bg-emerald-500/10 px-1.5 py-0.5 font-mono text-[10px] font-bold text-emerald-400">
              CROWN / ILGP
            </span>
          </div>
          <div className="mt-2 text-2xl font-black font-mono text-emerald-400">
            {metrics.sovereignCoveragePct.toFixed(1)}%
          </div>
          <p className="mt-1 text-[11px] text-text-subtle">
            Capital de-risked by Federal ITC, CIB concessionary debt, or ILGP guarantee.
          </p>
        </div>

        {/* Annual Interest Savings */}
        <div className="rounded-xl border border-borderSubtle bg-surface p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-text-muted">Annual Sovereign Spread Alpha</span>
            <span className="rounded bg-blue-500/10 px-1.5 py-0.5 font-mono text-[10px] font-bold text-blue-400">
              SAVINGS
            </span>
          </div>
          <div className="mt-2 text-2xl font-black font-mono text-blue-400">
            ${(metrics.annualInterestSavings / 1e6).toFixed(1)}M/yr
          </div>
          <p className="mt-1 text-[11px] text-text-subtle">
            Annual debt service saved vs standard commercial borrowing without federal support.
          </p>
        </div>
      </div>
    </div>
  );
}
