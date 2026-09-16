"use client";

import { useMemo, useState } from "react";
import {
  Calculator,
  CircleDollarSign,
  ShieldAlert,
  Landmark,
  Sparkles,
  Percent,
  TrendingUp,
  CheckCircle2,
  Building2,
  FileText,
} from "lucide-react";

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max);
}

function formatCad(value: number) {
  return new Intl.NumberFormat("en-CA", {
    style: "currency",
    currency: "CAD",
    maximumFractionDigits: 0,
  }).format(value);
}

export default function CapitalPage() {
  // General Capital Stack State
  const [capexMillions, setCapexMillions] = useState(500);
  const [seniorDebt, setSeniorDebt] = useState(45);
  const [sponsorEquity, setSponsorEquity] = useState(35);
  const [otherCapital, setOtherCapital] = useState(20);

  // Pillar B: Indigenous Loan Guarantee Program (ILGP) Simulator State
  const [ilgpCapexMillions, setIlgpCapexMillions] = useState(750);
  const [indigenousEquityPct, setIndigenousEquityPct] = useState(20);
  const [loanTermYears, setLoanTermYears] = useState(20);
  const [interestSpreadBps, setInterestSpreadBps] = useState(225); // 2.25% AAA Sovereign Spread Benefit
  const [projectReturnPct, setProjectReturnPct] = useState(10.5);

  const model = useMemo(() => {
    const capex = clamp(capexMillions, 0, 1_000_000) * 1_000_000;
    const percentages = [seniorDebt, sponsorEquity, otherCapital].map((value) => clamp(value, 0, 100));
    const totalPercent = percentages.reduce((total, value) => total + value, 0);
    return {
      capex,
      totalPercent,
      allocations: [
        { label: "Senior debt assumption", percent: percentages[0], value: (capex * percentages[0]) / 100 },
        { label: "Sponsor equity assumption", percent: percentages[1], value: (capex * percentages[1]) / 100 },
        { label: "Other capital assumption", percent: percentages[2], value: (capex * percentages[2]) / 100 },
      ],
    };
  }, [capexMillions, seniorDebt, sponsorEquity, otherCapital]);

  // Pillar B: ILGP Financial Engineering Engine
  const ilgpModel = useMemo(() => {
    const capex = clamp(ilgpCapexMillions, 10, 50_000) * 1_000_000;
    const equityPct = clamp(indigenousEquityPct, 1, 100) / 100;
    const term = clamp(loanTermYears, 5, 35);
    const spreadPct = clamp(interestSpreadBps, 50, 500) / 10000;
    const returnPct = clamp(projectReturnPct, 2, 30) / 100;

    // Equity stake dollar value
    const equityStakeValue = capex * equityPct;

    // Under the $5B Federal ILGP, up to 100% of the Indigenous equity debt tranche is guaranteed
    const guaranteedDebt = equityStakeValue * 0.95; // 95% debt-financed equity acquisition
    const initialCashRequired = equityStakeValue * 0.05;

    // Annual interest savings from AAA Sovereign Crown guarantee vs sub-investment grade commercial loan
    const annualInterestSavings = guaranteedDebt * spreadPct;
    const cumulativeInterestSavings = annualInterestSavings * term;

    // Annual gross equity cash distribution to community
    const grossAnnualReturn = equityStakeValue * returnPct;

    // Estimated debt service (principal amort + subsidized interest ~3.75% sovereign base)
    const baseInterestRate = 0.0375;
    const annualDebtService = (guaranteedDebt / term) + (guaranteedDebt * baseInterestRate);

    // Net annual sovereign dividend to community
    const netAnnualDividend = Math.max(0, grossAnnualReturn - annualDebtService);
    const cumulativeNetDividends = netAnnualDividend * term + (grossAnnualReturn * 10); // 10 years unencumbered post-amortization

    return {
      capex,
      equityStakeValue,
      guaranteedDebt,
      initialCashRequired,
      annualInterestSavings,
      cumulativeInterestSavings,
      grossAnnualReturn,
      annualDebtService,
      netAnnualDividend,
      cumulativeNetDividends,
      term,
    };
  }, [ilgpCapexMillions, indigenousEquityPct, loanTermYears, interestSpreadBps, projectReturnPct]);

  const controls = [
    { label: "Senior debt", value: seniorDebt, setValue: setSeniorDebt },
    { label: "Sponsor equity", value: sponsorEquity, setValue: setSponsorEquity },
    { label: "Other capital", value: otherCapital, setValue: setOtherCapital },
  ];

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-4 py-8 sm:px-6 lg:px-8">
      {/* Header */}
      <header className="border-b border-border/80 pb-6">
        <div className="mb-3 inline-flex items-center gap-2 rounded-full border border-primary/40 bg-[#040806] px-3 py-1 font-mono text-[11px] text-aurora">
          {/* Canadian Government FIP Flaglet */}
          <div className="flex items-center gap-1.5 text-white">
            <svg className="h-3 w-5 shrink-0 rounded-[1px] overflow-hidden border border-white/20" viewBox="0 0 100 50">
              <rect width="25" height="50" fill="#D8292F" />
              <rect x="25" width="50" height="50" fill="#FFFFFF" />
              <rect x="75" width="25" height="50" fill="#D8292F" />
              <path
                d="M 50 10 L 52 18 L 59 15 L 56 22 L 64 22 L 59 27 L 66 33 L 57 33 L 54 36 L 53 43 L 51 43 L 50 41 L 49 43 L 47 43 L 46 36 L 43 33 L 34 33 L 41 27 L 36 22 L 44 22 L 41 15 L 48 18 Z"
                fill="#D8292F"
              />
            </svg>
            <span className="font-bold text-white tracking-tight">FINANCE CANADA / CIB / NRCAN</span>
          </div>
          <span className="text-border">|</span>
          <span>CAPITAL ARITHMETIC & SYNDICATION</span>
        </div>
        <h1 className="text-2xl font-black tracking-tight text-text-main sm:text-4xl">
          National Infrastructure <span className="text-aurora">Capital Architecture</span>
        </h1>
        <p className="mt-2 max-w-3xl text-sm leading-relaxed text-text-muted">
          Transparent debt syndication, private sponsor equity allocation, and sovereign loan guarantee arithmetic.
          Includes the simulation framework for the $5.0 Billion Canada Indigenous Loan Guarantee Program (ILGP).
        </p>
      </header>

      {/* SECTION 1: Indigenous Loan Guarantee Program (ILGP) Simulator */}
      <section className="space-y-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <div className="inline-flex items-center gap-2 text-xs font-mono font-bold text-emerald-400 uppercase tracking-wider">
              <Landmark className="h-4 w-4" />
              Pillar B: Indigenous Economic Sovereignty
            </div>
            <h2 className="text-xl font-black text-text-main mt-1">
              $5B Federal Indigenous Loan Guarantee Simulator
            </h2>
          </div>
          <div className="rounded-xl border border-emerald-500/40 bg-emerald-950/20 px-3 py-1.5 text-xs font-mono text-emerald-300">
            CIB / NRCAN Sovereign Crown Facility
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
          {/* Controls */}
          <div className="glass-card space-y-5 rounded-2xl border border-border/80 p-6">
            <h3 className="font-mono text-xs uppercase font-bold text-text-muted tracking-wider border-b border-border/60 pb-2 flex items-center gap-1.5">
              <Building2 className="h-3.5 w-3.5 text-aurora" /> Project Parameters
            </h3>

            <label className="block text-xs font-mono text-text-muted">
              Total Project CAPEX (CAD Millions)
              <input
                type="number"
                min="10"
                max="50000"
                step="50"
                value={ilgpCapexMillions}
                onChange={(e) => setIlgpCapexMillions(Number(e.target.value) || 10)}
                className="mt-1.5 w-full rounded-xl border border-border bg-[#040806] px-3 py-2 text-base font-bold text-text-main outline-none focus:border-aurora"
              />
            </label>

            <label className="block text-xs font-mono text-text-muted">
              <div className="flex justify-between">
                <span>Indigenous Equity Stake</span>
                <strong className="text-aurora">{indigenousEquityPct}%</strong>
              </div>
              <input
                type="range"
                min="1"
                max="50"
                step="1"
                value={indigenousEquityPct}
                onChange={(e) => setIndigenousEquityPct(Number(e.target.value))}
                className="mt-2 w-full accent-emerald-400"
              />
              <span className="text-[10px] text-text-subtle">Target equity acquired by First Nations / Métis / Inuit partners</span>
            </label>

            <label className="block text-xs font-mono text-text-muted">
              <div className="flex justify-between">
                <span>Loan Amortization Tenure</span>
                <strong className="text-white">{loanTermYears} Years</strong>
              </div>
              <input
                type="range"
                min="5"
                max="30"
                step="1"
                value={loanTermYears}
                onChange={(e) => setLoanTermYears(Number(e.target.value))}
                className="mt-2 w-full accent-emerald-400"
              />
            </label>

            <label className="block text-xs font-mono text-text-muted">
              <div className="flex justify-between">
                <span>AAA Sovereign Spread Advantage</span>
                <strong className="text-emerald-400">{interestSpreadBps} bps ({(interestSpreadBps / 100).toFixed(2)}%)</strong>
              </div>
              <input
                type="range"
                min="50"
                max="400"
                step="25"
                value={interestSpreadBps}
                onChange={(e) => setInterestSpreadBps(Number(e.target.value))}
                className="mt-2 w-full accent-emerald-400"
              />
              <span className="text-[10px] text-text-subtle">Federal guarantee reduces borrowing rate from commercial high-yield down to Government of Canada AAA equivalent</span>
            </label>

            <label className="block text-xs font-mono text-text-muted">
              <div className="flex justify-between">
                <span>Projected Asset Return (Cash Flow Yield)</span>
                <strong className="text-gold">{projectReturnPct}%</strong>
              </div>
              <input
                type="range"
                min="4"
                max="20"
                step="0.5"
                value={projectReturnPct}
                onChange={(e) => setProjectReturnPct(Number(e.target.value))}
                className="mt-2 w-full accent-amber-400"
              />
            </label>
          </div>

          {/* Syndication Results & Dividends */}
          <div className="glass-card space-y-4 rounded-2xl border border-border/80 p-6">
            <div className="flex items-start justify-between gap-4 border-b border-borderSubtle pb-4">
              <div>
                <div className="text-[10px] font-mono uppercase text-text-subtle">Total Equity Value Acquired</div>
                <div className="mt-1 text-2xl font-black text-aurora">{formatCad(ilgpModel.equityStakeValue)}</div>
              </div>
              <span className="rounded-lg bg-emerald-500/20 border border-emerald-500/40 px-2.5 py-1 text-xs font-mono font-bold text-emerald-300">
                100% CIB Guaranteed
              </span>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-xl border border-borderSubtle bg-surface/80 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-text-subtle">Crown Loan Guarantee Amount</div>
                <div className="text-lg font-black text-white">{formatCad(ilgpModel.guaranteedDebt)}</div>
                <div className="text-[10px] text-text-muted">Non-recourse to community assets</div>
              </div>

              <div className="rounded-xl border border-emerald-500/30 bg-emerald-950/30 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-emerald-400">Annual Interest Savings</div>
                <div className="text-lg font-black text-emerald-300">{formatCad(ilgpModel.annualInterestSavings)}/yr</div>
                <div className="text-[10px] text-emerald-400/80">From AAA Crown backing</div>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-xl border border-borderSubtle bg-surface/80 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-text-subtle">Net Annual Community Cash Flow</div>
                <div className="text-lg font-black text-gold">{formatCad(ilgpModel.netAnnualDividend)}/yr</div>
                <div className="text-[10px] text-text-muted">After full debt servicing</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/80 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-text-subtle">Cumulative Interest Saved</div>
                <div className="text-lg font-black text-aurora">{formatCad(ilgpModel.cumulativeInterestSavings)}</div>
                <div className="text-[10px] text-text-muted">Across {ilgpModel.term}-year tenure</div>
              </div>
            </div>

            <div className="rounded-xl border border-aurora/40 bg-[#040806] p-4 space-y-2">
              <div className="flex items-center justify-between text-xs font-mono font-bold">
                <span className="text-text-muted">Intergenerational Sovereign Wealth</span>
                <span className="text-aurora">{formatCad(ilgpModel.cumulativeNetDividends)}</span>
              </div>
              <p className="text-[11px] text-text-subtle leading-relaxed">
                Projected cumulative unencumbered cash flows generated for Indigenous community governance, healthcare, and infrastructure endowments over asset lifetime.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* SECTION 2: General Capital Stack Scenario */}
      <section className="space-y-6 pt-4 border-t border-border/80">
        <div className="flex items-center gap-2">
          <Calculator className="h-4 w-4 text-aurora" />
          <h2 className="text-xl font-black text-text-main">
            Project Capital Stack Breakdown
          </h2>
        </div>

        <div className="grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
          <div className="glass-card space-y-6 rounded-2xl border border-border/80 p-6">
            <label className="block text-xs font-mono text-text-muted">
              Total project CAPEX (CAD millions)
              <input
                type="number"
                min="0"
                max="1000000"
                step="10"
                value={capexMillions}
                onChange={(event) => setCapexMillions(Number(event.target.value) || 0)}
                className="mt-2 w-full rounded-xl border border-border bg-[#040806] px-3 py-2 text-base font-bold text-text-main outline-none focus:border-aurora"
              />
            </label>

            {controls.map((control) => (
              <label key={control.label} className="block text-xs font-mono text-text-muted">
                <span className="flex justify-between">
                  <span>{control.label}</span>
                  <strong className="text-text-main">{control.value}%</strong>
                </span>
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="1"
                  value={control.value}
                  onChange={(event) => control.setValue(Number(event.target.value))}
                  className="mt-2 w-full accent-emerald-400"
                />
              </label>
            ))}
          </div>

          <div className="glass-card space-y-5 rounded-2xl border border-border/80 p-6">
            <div className="flex items-start justify-between gap-4 border-b border-borderSubtle pb-4">
              <div>
                <div className="text-[10px] font-mono uppercase text-text-subtle">Scenario CAPEX</div>
                <div className="mt-1 text-2xl font-black text-gold">{formatCad(model.capex)}</div>
              </div>
              <CircleDollarSign className="h-7 w-7 text-aurora" />
            </div>

            {model.allocations.map((allocation) => (
              <div key={allocation.label} className="rounded-xl border border-borderSubtle bg-surface p-4">
                <div className="flex items-center justify-between gap-3 text-xs">
                  <span className="text-text-muted">{allocation.label}</span>
                  <span className="font-mono font-bold text-text-main">{allocation.percent}%</span>
                </div>
                <div className="mt-1 text-lg font-black text-aurora">{formatCad(allocation.value)}</div>
              </div>
            ))}

            <div
              aria-live="polite"
              className={`rounded-xl border p-4 text-xs ${
                model.totalPercent === 100
                  ? "border-primary/40 bg-primary/10 text-aurora"
                  : "border-gold/40 bg-gold/10 text-gold"
              }`}
            >
              Allocation total: <strong>{model.totalPercent}%</strong>.{" "}
              {model.totalPercent === 100
                ? "The arithmetic balances to total CAPEX."
                : "Adjust assumptions until the stack totals 100%."}
            </div>
          </div>
        </div>
      </section>

      {/* Institutional Disclaimer */}
      <div className="flex gap-3 rounded-2xl border border-borderSubtle bg-surface p-5 text-xs leading-relaxed text-text-subtle">
        <ShieldAlert className="h-5 w-5 shrink-0 text-gold" />
        <p>
          All outputs are arithmetic scenarios for strategic decision support. Actual project underwriting under the
          Canada Indigenous Loan Guarantee Program (ILGP), Alberta Indigenous Opportunities Corporation (AIOC), Ontario
          Aboriginal Loan Guarantee Program (ALGP), or private debt markets requires formal application, credit committee review,
          and agreement with respective sovereign and First Nation governing bodies.
        </p>
      </div>
    </div>
  );
}
