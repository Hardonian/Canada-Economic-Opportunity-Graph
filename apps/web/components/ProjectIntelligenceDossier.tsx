import React from "react";
import { 
  BarChart3, 
  ShieldAlert, 
  Zap, 
  Satellite, 
  TrendingUp, 
  Building2, 
  AlertTriangle, 
  CheckCircle2, 
  Landmark,
  Scale,
  Gauge
} from "lucide-react";
import { Project } from "@/lib/types";
import { 
  calculateMRIO, 
  calculateFlyvbjerg, 
  calculateUBOScreening, 
  calculateGridAssessment, 
  calculateEarthObs,
  MRIOResult,
  FlyvbjergResult,
  UBOScreeningResult,
  GridFeasibilityResult,
  EarthObsResult
} from "@/lib/intelligence";

interface Props {
  project: Project;
  initialMrio?: MRIOResult;
  initialFlyvbjerg?: FlyvbjergResult;
  initialUbo?: UBOScreeningResult;
  initialGrid?: GridFeasibilityResult;
  initialEarthobs?: EarthObsResult;
}

export default function ProjectIntelligenceDossier({ 
  project, 
  initialMrio, 
  initialFlyvbjerg, 
  initialUbo, 
  initialGrid, 
  initialEarthobs 
}: Props) {
  const mrio = initialMrio ?? calculateMRIO(project);
  const flyvbjerg = initialFlyvbjerg ?? calculateFlyvbjerg(project);
  const ubo = initialUbo ?? calculateUBOScreening(project);
  const grid = initialGrid ?? calculateGridAssessment(project);
  const earthobs = initialEarthobs ?? calculateEarthObs(project);

  const formatCAD = (n: number) => {
    if (n >= 1e9) return `$${(n / 1e9).toFixed(2)}B`;
    if (n >= 1e6) return `$${(n / 1e6).toFixed(1)}M`;
    if (n >= 1e3) return `$${(n / 1e3).toFixed(0)}K`;
    return `$${n.toLocaleString()}`;
  };

  return (
    <section className="space-y-6" aria-labelledby="institutional-intelligence-title">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-border/80 pb-3">
        <div>
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-aurora animate-pulse"></span>
            <h2 id="institutional-intelligence-title" className="text-sm font-bold font-mono uppercase tracking-wider text-text-main">
              Institutional Intelligence & Risk Dossier
            </h2>
          </div>
          <p className="text-xs text-text-muted mt-0.5">
            Deterministic macro-econometric, Bayesian overrun, sovereign screening, electrical grid, and satellite ground-truth analytics.
          </p>
        </div>
        <div className="flex items-center gap-2 text-[11px] font-mono text-text-subtle">
          <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-aurora">
            CEGS-AI v2.0
          </span>
          <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle">
            Zero-Hallucination
          </span>
        </div>
      </div>

      {/* Grid: 2 columns on large screens */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        
        {/* 1. StatCan MRIO Macro Multipliers */}
        <div className="glass-card p-6 rounded-2xl border border-border/80 space-y-4 shadow-xl relative overflow-hidden">
          <div className="flex items-center justify-between border-b border-borderSubtle pb-3">
            <div className="flex items-center gap-2">
              <BarChart3 className="h-4 w-4 text-aurora" />
              <h3 className="font-mono text-xs font-bold uppercase text-text-main">
                StatCan Input-Output Macro Multipliers
              </h3>
            </div>
            <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-primary/20 text-aurora font-bold">
              {mrio.multiplier.toFixed(2)}x Multiplier
            </span>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[10px] font-mono uppercase text-text-subtle">Total GDP Generated</div>
              <div className="text-lg font-bold font-tabular text-aurora mt-0.5">{formatCAD(mrio.totalGDPCAD)}</div>
              <div className="text-[9px] text-text-muted">Direct + Indirect + Induced</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[10px] font-mono uppercase text-text-subtle">Employment Created</div>
              <div className="text-lg font-bold font-tabular text-text-main mt-0.5">{mrio.personYearsFTE.toLocaleString()}</div>
              <div className="text-[9px] text-text-muted">Person-Years (FTE)</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle col-span-2 sm:col-span-1">
              <div className="text-[10px] font-mono uppercase text-text-subtle">Total Fiscal Return</div>
              <div className="text-lg font-bold font-tabular text-gold mt-0.5">{formatCAD(mrio.totalFiscalReturnCAD)}</div>
              <div className="text-[9px] text-text-muted">3-Tier Public Tax Yield</div>
            </div>
          </div>

          {/* Detailed Tax Yields */}
          <div className="space-y-2 pt-1">
            <div className="text-[11px] font-mono text-text-subtle uppercase">Tax Return Breakdown (StatCan SUT):</div>
            <div className="grid grid-cols-3 gap-2 text-center text-xs">
              <div className="p-2 rounded-lg bg-surface/40 border border-borderSubtle">
                <div className="text-[9px] text-text-subtle">Federal Tax</div>
                <div className="font-bold text-text-main font-tabular">{formatCAD(mrio.federalTaxCAD)}</div>
              </div>
              <div className="p-2 rounded-lg bg-surface/40 border border-borderSubtle">
                <div className="text-[9px] text-text-subtle">Provincial Tax</div>
                <div className="font-bold text-text-main font-tabular">{formatCAD(mrio.provincialTaxCAD)}</div>
              </div>
              <div className="p-2 rounded-lg bg-surface/40 border border-borderSubtle">
                <div className="text-[9px] text-text-subtle">Municipal Tax</div>
                <div className="font-bold text-text-main font-tabular">{formatCAD(mrio.municipalTaxCAD)}</div>
              </div>
            </div>
          </div>
        </div>

        {/* 2. Bayesian Flyvbjerg Overrun Hazard Curve */}
        <div className="glass-card p-6 rounded-2xl border border-border/80 space-y-4 shadow-xl">
          <div className="flex items-center justify-between border-b border-borderSubtle pb-3">
            <div className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4 text-gold" />
              <h3 className="font-mono text-xs font-bold uppercase text-text-main">
                Bayesian Overrun Hazard (Flyvbjerg)
              </h3>
            </div>
            <span className="text-[10px] font-mono text-text-subtle truncate max-w-[180px]">
              {flyvbjerg.referenceClass}
            </span>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[10px] font-mono uppercase text-text-subtle">Expected Cost Overrun</div>
              <div className="text-lg font-bold font-tabular text-gold mt-0.5">+{flyvbjerg.expectedCostOverrunPct}%</div>
              <div className="text-[9px] text-text-muted">Sample: {flyvbjerg.historicalSampleSize} projects</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[10px] font-mono uppercase text-text-subtle">Expected Schedule Slip</div>
              <div className="text-lg font-bold font-tabular text-text-main mt-0.5">+{flyvbjerg.expectedDelayMonths} mo</div>
              <div className="text-[9px] text-text-muted">
                {flyvbjerg.remotePenaltyPct > 0 ? `+${flyvbjerg.remotePenaltyPct}% Remote geographic penalty` : "Standard geography"}
              </div>
            </div>
          </div>

          {/* Hazard Quantiles Table */}
          <div className="space-y-1.5">
            <div className="text-[11px] font-mono text-text-subtle uppercase">Lognormal Stress Quantiles:</div>
            <div className="grid grid-cols-4 gap-1.5 text-center text-[10px] font-mono">
              {flyvbjerg.quantiles.map((q) => (
                <div key={q.percentile} className="p-2 rounded-lg bg-surface/50 border border-borderSubtle">
                  <div className="text-text-subtle font-bold">P{q.percentile}</div>
                  <div className="font-bold text-text-main mt-0.5">+{q.costOverrunPct}%</div>
                  <div className="text-[9px] text-aurora">{formatCAD(q.forecastCapexCAD)}</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* 3. Sovereign Screening & Investment Canada Act */}
        <div className="glass-card p-6 rounded-2xl border border-border/80 space-y-4 shadow-xl">
          <div className="flex items-center justify-between border-b border-borderSubtle pb-3">
            <div className="flex items-center gap-2">
              <Scale className="h-4 w-4 text-aurora" />
              <h3 className="font-mono text-xs font-bold uppercase text-text-main">
                Sovereign Screening & ICA Review
              </h3>
            </div>
            <span className={`text-[10px] font-mono px-2.5 py-0.5 rounded-full font-bold border ${
              ubo.icaRisk === "CLEAR" 
                ? "bg-emerald-500/10 border-emerald-500/40 text-emerald-400" 
                : "bg-amber-500/10 border-amber-500/40 text-amber-400"
            }`}>
              ICA: {ubo.icaRisk}
            </span>
          </div>

          <div className="grid grid-cols-3 gap-2 text-center text-xs font-mono">
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">Domestic Control</div>
              <div className="text-base font-bold text-aurora mt-0.5">{ubo.domesticControlPct}%</div>
              <div className="text-[8px] text-text-muted">Canadian Control</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">FTA Partner Share</div>
              <div className="text-base font-bold text-text-main mt-0.5">{ubo.ftaPartnerPct}%</div>
              <div className="text-[8px] text-text-muted">USMCA / CETA</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">Foreign SOE Share</div>
              <div className="text-base font-bold text-text-main mt-0.5">{ubo.soeExposurePct}%</div>
              <div className="text-[8px] text-text-muted">State-Owned</div>
            </div>
          </div>

          <div className="space-y-1.5 text-xs">
            {ubo.notes.map((note, idx) => (
              <div key={idx} className="flex items-start gap-2 text-[11px] text-text-muted bg-surface/30 p-2 rounded-lg border border-borderSubtle">
                <CheckCircle2 className="h-3.5 w-3.5 text-aurora shrink-0 mt-0.5" />
                <span>{note}</span>
              </div>
            ))}
          </div>
        </div>

        {/* 4. Electrical Grid Feasibility & Interconnect */}
        <div className="glass-card p-6 rounded-2xl border border-border/80 space-y-4 shadow-xl">
          <div className="flex items-center justify-between border-b border-borderSubtle pb-3">
            <div className="flex items-center gap-2">
              <Zap className="h-4 w-4 text-primary" />
              <h3 className="font-mono text-xs font-bold uppercase text-text-main">
                Electrical Grid Feasibility & Headroom
              </h3>
            </div>
            <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-surface border border-borderSubtle text-text-main">
              {grid.systemOperator}
            </span>
          </div>

          <div className="grid grid-cols-3 gap-2 text-center text-xs font-mono">
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">Feasibility Score</div>
              <div className="text-base font-bold text-aurora mt-0.5">{grid.feasibilityScore}/100</div>
              <div className="text-[8px] text-text-muted">Transmission Rating</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">Clean Power Purity</div>
              <div className="text-base font-bold text-text-main mt-0.5">{grid.cleanPurityPct}%</div>
              <div className="text-[8px] text-text-muted">Hydro / Nuclear</div>
            </div>
            <div className="p-3 rounded-xl bg-surface/70 border border-borderSubtle">
              <div className="text-[9px] text-text-subtle uppercase">Queue Estimated</div>
              <div className="text-base font-bold text-gold mt-0.5">{grid.queueMonths} mo</div>
              <div className="text-[8px] text-text-muted">ISO Interconnect</div>
            </div>
          </div>

          <div className="p-3 rounded-xl bg-surface/40 border border-borderSubtle text-xs space-y-1">
            <div className="flex justify-between text-text-muted">
              <span>Estimated Load / Gen:</span>
              <span className="font-mono text-text-main font-bold">{grid.estimatedMW} MW ({grid.interconnectVoltageKV} kV)</span>
            </div>
            <div className="flex justify-between text-text-muted">
              <span>Substation Headroom:</span>
              <span className="font-mono text-text-main font-bold">{grid.headroomMW} MW</span>
            </div>
            <div className="flex justify-between text-text-muted">
              <span>Dedicated Substation Required:</span>
              <span className="font-mono text-text-main font-bold">{grid.dedicatedSubstation ? "Yes" : "No"}</span>
            </div>
            <div className="flex justify-between text-text-muted">
              <span>Reinforcement Capital:</span>
              <span className="font-mono text-aurora font-bold">{formatCAD(grid.reinforcementCapexCAD)}</span>
            </div>
          </div>
        </div>

      </div>

      {/* 5. Full-Width Satellite Ground-Truth Telemetry (SAR & Optical) */}
      <div className="glass-card p-6 rounded-2xl border border-border/80 space-y-4 shadow-xl">
        <div className="flex items-center justify-between border-b border-borderSubtle pb-3">
          <div className="flex items-center gap-2">
            <Satellite className="h-4 w-4 text-aurora" />
            <h3 className="font-mono text-xs font-bold uppercase text-text-main">
              Earth Observation SAR & Optical Ground-Truth Telemetry
            </h3>
          </div>
          <div className="flex items-center gap-2">
            <span className={`text-[10px] font-mono px-2.5 py-0.5 rounded-full font-bold border ${
              earthobs.corroborationStatus === "VERIFIED"
                ? "bg-emerald-500/10 border-emerald-500/40 text-emerald-400"
                : "bg-blue-500/10 border-blue-500/40 text-blue-400"
            }`}>
              {earthobs.corroborationStatus}
            </span>
            <span className="text-[10px] font-mono text-text-subtle hidden sm:inline">
              Pass: {earthobs.lastSatellitePass}
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
          <div className="p-4 rounded-xl bg-surface/70 border border-borderSubtle space-y-2">
            <div className="flex justify-between text-xs">
              <span className="font-mono text-text-subtle">Physical Progress Score:</span>
              <span className="font-mono text-aurora font-bold">{earthobs.physicalProgressScore.toFixed(1)}/100</span>
            </div>
            <div className="w-full h-2 rounded-full bg-surface overflow-hidden border border-borderSubtle">
              <div 
                className="h-full bg-primary rounded-full shadow-[0_0_8px_#00F5A0]" 
                style={{ width: `${earthobs.physicalProgressScore}%` }}
              ></div>
            </div>
          </div>

          <div className="sm:col-span-2 p-4 rounded-xl bg-surface/40 border border-borderSubtle text-xs space-y-1.5">
            <div className="flex items-center gap-2 text-[11px] text-text-subtle font-mono">
              <span className="text-aurora">Telemetry Source:</span>
              <span>{earthobs.sensorConstellation}</span>
            </div>
            <p className="text-text-muted text-xs leading-relaxed">
              {earthobs.telemetrySummary}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
