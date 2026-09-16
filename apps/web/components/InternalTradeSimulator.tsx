"use client";

import React, { useState, useMemo } from "react";
import {
  TrendingUp,
  AlertOctagon,
  CheckCircle2,
  Sliders,
  Scale,
  ArrowRightLeft,
  DollarSign,
  Clock,
  Sparkles,
  ShieldCheck,
} from "lucide-react";

interface CorridorData {
  id: string;
  name: string;
  provinces: string;
  volumeCAD: number; // in Billions
  frictionRatePct: number; // Tariff equivalent %
  primaryCommodities: string[];
  regulatoryChokepoints: string[];
}

const CORRIDORS: CorridorData[] = [
  {
    id: "ab-bc",
    name: "Western Energy & Gateway Corridor",
    provinces: "Alberta ↔ British Columbia",
    volumeCAD: 48.5,
    frictionRatePct: 7.4,
    primaryCommodities: ["Clean Hydrogen & LNG", "Agricultural Grain", "Forestry", "Rare Earth Elements"],
    regulatoryChokepoints: ["Duplicative Provincial Environmental Review", "Commercial Trucking Axle Discrepancies"],
  },
  {
    id: "on-qc",
    name: "St. Lawrence - Great Lakes Manufacturing Spine",
    provinces: "Ontario ↔ Quebec",
    volumeCAD: 94.2,
    frictionRatePct: 6.8,
    primaryCommodities: ["EV Battery Materials", "Hydroelectricity", "Automotive Assemblies", "Refined Nickel"],
    regulatoryChokepoints: ["Bilingual Safety Permitting Duplication", "Interprovincial High-Voltage Intertie Rules"],
  },
  {
    id: "mb-sk",
    name: "Prairie Agricultural & Northern Rail Gateway",
    provinces: "Manitoba ↔ Saskatchewan",
    volumeCAD: 28.6,
    frictionRatePct: 7.9,
    primaryCommodities: ["Agricultural Potash", "Specialty Crops", "Uranium Fuels", "Rail Maintenance Machinery"],
    regulatoryChokepoints: ["Differing Heavy-Haul Transport Limits", "Labour Credential Recognition Lags"],
  },
  {
    id: "atlantic",
    name: "Atlantic Marine & Energy Intertie",
    provinces: "Nova Scotia ↔ New Brunswick ↔ NL",
    volumeCAD: 19.8,
    frictionRatePct: 8.5,
    primaryCommodities: ["Offshore Wind", "Seafood & Aquaculture", "Direct-Reduction Pellets"],
    regulatoryChokepoints: ["Fragmented Procurement Rules", "Inter-Utility Transmission Tariffs"],
  },
];

export default function InternalTradeSimulator() {
  const [harmonizationPct, setHarmonizationPct] = useState<number>(40);
  const [oneProjectOneReview, setOneProjectOneReview] = useState<boolean>(true);
  const [transportReciprocity, setTransportReciprocity] = useState<boolean>(true);

  // Economic calculations
  const sim = useMemo(() => {
    const baseFrictionTaxTotal = 130.0; // $130B CAD annual friction according to StatsCan / IMF
    const baseVelocity = CORRIDORS.reduce((acc, c) => acc + c.volumeCAD, 0);

    // One Project One Review adds 15% bonus efficiency
    const reviewBonus = oneProjectOneReview ? 0.15 : 0;
    // Transport reciprocity adds 10% bonus efficiency
    const transportBonus = transportReciprocity ? 0.10 : 0;

    const totalEfficiencyFactor = Math.min(1.0, (harmonizationPct / 100) + reviewBonus + transportBonus);

    // Friction saved
    const frictionSavedCAD = baseFrictionTaxTotal * totalEfficiencyFactor * 0.72;
    const remainingFrictionCAD = baseFrictionTaxTotal - frictionSavedCAD;

    // Projected GDP dividend (IMF models up to ~3.8% GDP increase from dismantling internal barriers)
    const gdpDividendCAD = frictionSavedCAD * 0.92;
    const gdpGrowthPct = (gdpDividendCAD / 2800) * 100; // Canadian nominal GDP ~$2.8T

    // Inter-provincial trade velocity expansion
    const velocityExpansionCAD = baseVelocity * (totalEfficiencyFactor * 0.28);

    // Timeline compression for approvals
    const timelineReductionMonths = Math.round(6 + totalEfficiencyFactor * 18);

    return {
      totalEfficiencyFactor,
      frictionSavedCAD,
      remainingFrictionCAD,
      gdpDividendCAD,
      gdpGrowthPct,
      velocityExpansionCAD,
      timelineReductionMonths,
    };
  }, [harmonizationPct, oneProjectOneReview, transportReciprocity]);

  return (
    <div className="glass-card rounded-2xl border border-primary/40 bg-[#040806]/90 p-6 sm:p-8 space-y-6 shadow-2xl">
      {/* Title */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border/80 pb-5">
        <div>
          <div className="inline-flex items-center gap-2 text-xs font-mono font-bold text-aurora uppercase tracking-wider">
            <Scale className="h-4 w-4" />
            Pillar C: Internal Trade Sovereignty & Regulatory Reform
          </div>
          <h2 className="text-xl sm:text-2xl font-black text-text-main mt-1">
            Inter-Provincial Trade & Regulatory Friction Simulator
          </h2>
          <p className="text-xs text-text-muted mt-1 max-w-2xl">
            Internal trade barriers impose an estimated 6.7%–8.5% implicit tariff on goods moving across provincial boundaries, costing the Canadian economy up to $130B/year. Model the national dividend of the Canadian Free Trade Agreement (CFTA) modernization and One Project, One Review harmonization.
          </p>
        </div>
        <div className="rounded-xl border border-aurora/40 bg-[#0C1812] px-4 py-2 text-right">
          <div className="text-[10px] font-mono uppercase text-text-subtle">Baseline Friction Tax</div>
          <div className="text-xl font-black text-red-400 font-mono">~$130.0B CAD/yr</div>
        </div>
      </div>

      {/* Simulator Inputs & Telemetry Grid */}
      <div className="grid gap-6 lg:grid-cols-[1fr_1.1fr]">
        {/* Controls */}
        <div className="space-y-5 bg-[#0C1812]/70 rounded-xl border border-border/70 p-5 font-mono text-xs">
          <div className="flex items-center justify-between font-bold border-b border-border/50 pb-2">
            <span className="flex items-center gap-2 text-white">
              <Sliders className="h-4 w-4 text-aurora" /> Harmonization Policy Levers
            </span>
            <span className="text-aurora">{(sim.totalEfficiencyFactor * 100).toFixed(0)}% Accord Factor</span>
          </div>

          <label className="block space-y-2 text-text-muted">
            <div className="flex justify-between">
              <span>Provincial Regulatory Accord Level</span>
              <strong className="text-aurora font-bold text-sm">{harmonizationPct}%</strong>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              step="5"
              value={harmonizationPct}
              onChange={(e) => setHarmonizationPct(Number(e.target.value))}
              className="w-full accent-emerald-400"
            />
            <div className="flex justify-between text-[10px] text-text-subtle">
              <span>Status Quo (Fragmented)</span>
              <span>Full Mutual Recognition</span>
            </div>
          </label>

          {/* Toggle 1: One Project, One Review */}
          <div
            onClick={() => setOneProjectOneReview(!oneProjectOneReview)}
            className={`p-3 rounded-xl border cursor-pointer transition-colors flex items-start gap-3 ${
              oneProjectOneReview
                ? "bg-aurora/10 border-aurora/60 text-white"
                : "bg-surface/50 border-border text-text-muted hover:border-border/80"
            }`}
          >
            <div className={`mt-0.5 rounded p-0.5 ${oneProjectOneReview ? "text-aurora" : "text-text-subtle"}`}>
              <CheckCircle2 className="h-4 w-4" />
            </div>
            <div>
              <div className="font-bold text-[11px] text-white flex items-center gap-1.5">
                <span>&quot;One Project, One Review&quot; Accord</span>
                {oneProjectOneReview && <span className="text-[9px] px-1.5 py-0.2 rounded bg-aurora/20 text-aurora font-mono">+15% EFF</span>}
              </div>
              <p className="text-[10px] text-text-subtle mt-0.5 leading-relaxed">
                Eliminates dual federal-provincial permitting review delays for clean infrastructure, interties, and critical mineral mines.
              </p>
            </div>
          </div>

          {/* Toggle 2: Transport & Axle Reciprocity */}
          <div
            onClick={() => setTransportReciprocity(!transportReciprocity)}
            className={`p-3 rounded-xl border cursor-pointer transition-colors flex items-start gap-3 ${
              transportReciprocity
                ? "bg-sky-500/10 border-sky-500/60 text-white"
                : "bg-surface/50 border-border text-text-muted hover:border-border/80"
            }`}
          >
            <div className={`mt-0.5 rounded p-0.5 ${transportReciprocity ? "text-sky-400" : "text-text-subtle"}`}>
              <CheckCircle2 className="h-4 w-4" />
            </div>
            <div>
              <div className="font-bold text-[11px] text-white flex items-center gap-1.5">
                <span>Unified Heavy Haul & Axle-Load Standards</span>
                {transportReciprocity && <span className="text-[9px] px-1.5 py-0.2 rounded bg-sky-500/20 text-sky-400 font-mono">+10% EFF</span>}
              </div>
              <p className="text-[10px] text-text-subtle mt-0.5 leading-relaxed">
                Standardizes trucking axle configurations across provincial boundaries, eliminating break-bulk transfers at borders.
              </p>
            </div>
          </div>
        </div>

        {/* Real-time Economic Telemetry */}
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <div className="bg-[#0C1812] rounded-xl border border-aurora/40 p-4 space-y-1">
              <div className="text-[10px] font-mono text-text-subtle uppercase">Annual Friction Eliminated</div>
              <div className="text-2xl font-black text-aurora font-mono">
                ${sim.frictionSavedCAD.toFixed(1)}B <span className="text-xs font-normal text-text-muted">CAD/yr</span>
              </div>
              <div className="text-[10px] text-emerald-400">Direct business cost savings</div>
            </div>

            <div className="bg-[#0C1812] rounded-xl border border-gold/40 p-4 space-y-1">
              <div className="text-[10px] font-mono text-text-subtle uppercase">National GDP Dividend</div>
              <div className="text-2xl font-black text-gold font-mono">
                +${sim.gdpDividendCAD.toFixed(1)}B <span className="text-xs font-normal text-text-muted">CAD</span>
              </div>
              <div className="text-[10px] text-gold font-bold">+{sim.gdpGrowthPct.toFixed(2)}% permanent GDP boost</div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="bg-[#0C1812] rounded-xl border border-sky-500/30 p-4 space-y-1">
              <div className="text-[10px] font-mono text-text-subtle uppercase">Trade Velocity Gain</div>
              <div className="text-xl font-black text-sky-300 font-mono">
                +${sim.velocityExpansionCAD.toFixed(1)}B <span className="text-xs font-normal text-text-muted">CAD</span>
              </div>
              <div className="text-[10px] text-text-muted">Interprovincial freight surge</div>
            </div>

            <div className="bg-[#0C1812] rounded-xl border border-purple-500/30 p-4 space-y-1">
              <div className="text-[10px] font-mono text-text-subtle uppercase">Approval Timeline Saved</div>
              <div className="text-xl font-black text-purple-300 font-mono">
                -{sim.timelineReductionMonths} <span className="text-xs font-normal text-text-muted">Months</span>
              </div>
              <div className="text-[10px] text-text-muted">Accelerated project COD</div>
            </div>
          </div>

          <div className="bg-[#040806] rounded-xl border border-border/80 p-4 text-xs font-mono space-y-2">
            <div className="flex items-center justify-between text-text-muted">
              <span>Remaining Interprovincial Friction Tax:</span>
              <span className="text-red-400 font-bold">${sim.remainingFrictionCAD.toFixed(1)}B CAD/yr</span>
            </div>
            <div className="w-full bg-surface rounded-full h-2 overflow-hidden border border-border">
              <div
                className="bg-aurora h-full transition-all duration-300"
                style={{ width: `${(sim.frictionSavedCAD / 130) * 100}%` }}
              />
            </div>
            <div className="flex justify-between text-[10px] text-text-subtle">
              <span>Friction Reduced: {((sim.frictionSavedCAD / 130) * 100).toFixed(0)}%</span>
              <span>Potential Max: $130.0B</span>
            </div>
          </div>
        </div>
      </div>

      {/* Strategic Corridor Matrix */}
      <div className="space-y-3 pt-2">
        <h3 className="font-mono text-xs uppercase font-bold text-text-muted tracking-wider flex items-center gap-1.5">
          <ArrowRightLeft className="h-4 w-4 text-aurora" /> Key Inter-Provincial Trade Corridors
        </h3>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {CORRIDORS.map((c) => {
            const corridorFrictionSaved = (c.volumeCAD * (c.frictionRatePct / 100) * sim.totalEfficiencyFactor * 0.75);
            return (
              <div key={c.id} className="rounded-xl border border-border/80 bg-[#0C1812]/60 p-4 space-y-2 font-mono text-xs">
                <div className="text-[10px] font-bold text-aurora uppercase">{c.provinces}</div>
                <div className="font-bold text-sm text-white">{c.name}</div>
                <div className="flex items-center justify-between text-[11px] pt-1 border-t border-border/40">
                  <span className="text-text-subtle">Annual Commerce:</span>
                  <span className="text-white font-bold">${c.volumeCAD}B CAD</span>
                </div>
                <div className="flex items-center justify-between text-[11px]">
                  <span className="text-text-subtle">Friction Saved:</span>
                  <span className="text-emerald-400 font-bold">+${corridorFrictionSaved.toFixed(2)}B</span>
                </div>
                <div className="pt-2 border-t border-border/40 text-[10px] text-text-subtle">
                  <span className="font-bold text-text-muted uppercase">Target Commodities:</span>
                  <div className="truncate mt-0.5 text-white">{c.primaryCommodities.join(", ")}</div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
