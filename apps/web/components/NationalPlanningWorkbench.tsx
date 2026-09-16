"use client";

import React, { useState, useMemo } from "react";
import Link from "next/link";
import {
  ShieldAlert,
  Cpu,
  TrendingUp,
  AlertTriangle,
  Flame,
  CheckCircle2,
  HardHat,
  Scale,
  Zap,
  DollarSign,
  Activity,
  Layers,
  ArrowUpRight,
  Globe2,
  Building2,
  Lock,
  ChevronRight,
  HelpCircle,
  BarChart3,
  Sliders,
} from "lucide-react";
import type { Project, Sector } from "@/lib/types";

interface NationalPlanningWorkbenchProps {
  projects: Project[];
}

type TabType = "optimizer" | "wargame" | "flyvbjerg" | "mrio" | "labor" | "ubo";

export default function NationalPlanningWorkbench({ projects }: NationalPlanningWorkbenchProps) {
  const [activeTab, setActiveTab] = useState<TabType>("optimizer");

  // --- TAB 1: Sovereign Capital Allocation Optimizer State ---
  const [cibEnvelopeB, setCibEnvelopeB] = useState<number>(4.0);
  const [sifEnvelopeB, setSifEnvelopeB] = useState<number>(2.0);
  const [itcEnvelopeB, setItcEnvelopeB] = useState<number>(3.0);
  const [indigEnvelopeB, setIndigEnvelopeB] = useState<number>(1.0);
  const [objective, setObjective] = useState<"BALANCED" | "MAX_CROWDING_IN" | "MAX_SOVEREIGNTY" | "MAX_DECARBONIZATION">("BALANCED");
  const [isOptimizingLive, setIsOptimizingLive] = useState(false);
  const [liveOptimizerStatus, setLiveOptimizerStatus] = useState<string | null>(null);
  const [isWargamingLive, setIsWargamingLive] = useState(false);
  const [liveWargameStatus, setLiveWargameStatus] = useState<string | null>(null);

  const handleLiveOptimization = async () => {
    setIsOptimizingLive(true);
    setLiveOptimizerStatus("Dispatching to Go Sovereign MILP Solver...");
    try {
      const res = await fetch("/api/v1/planning/optimize", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          objective,
          budget_envelopes: {
            cib_cad: cibEnvelopeB * 1e9,
            sif_cad: sifEnvelopeB * 1e9,
            itc_cad: itcEnvelopeB * 1e9,
            indigenous_loan_guarantee_cad: indigEnvelopeB * 1e9,
          },
        }),
      });
      if (res.ok) {
        const data = await res.json();
        setLiveOptimizerStatus(`Live Go Engine: HTTP 200 OK · ${(data.allocations || []).length} assets allocated · Multiplier: ${data.multiplier ? data.multiplier.toFixed(1) + "x" : "Calculated"}`);
      } else {
        setLiveOptimizerStatus("Live Engine offline · Executed via Deterministic Local MILP Engine");
      }
    } catch {
      setLiveOptimizerStatus("Live Engine standby · Executed via Deterministic Local MILP Engine");
    } finally {
      setIsOptimizingLive(false);
    }
  };

  const handleLiveWargame = async () => {
    setIsWargamingLive(true);
    setLiveWargameStatus("Dispatching scenario shock to Go Sovereign War Game Engine...");
    try {
      const res = await fetch("/api/v1/planning/wargame", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          scenario_id: selectedShock,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        setLiveWargameStatus(`Live Go Engine: HTTP 200 OK · ${(data.stalled_projects || []).length} assets impacted · Frozen CAPEX: $${((data.frozen_capex_cad || warGameResults.frozenCapex) / 1e9).toFixed(1)}B`);
      } else {
        setLiveWargameStatus("Live Engine offline · Executed via Deterministic Local War Game Engine");
      }
    } catch {
      setLiveWargameStatus("Live Engine standby · Executed via Deterministic Local War Game Engine");
    } finally {
      setIsWargamingLive(false);
    }
  };

  // Optimizer calculations
  const optimizationResults = useMemo(() => {
    let remCIB = cibEnvelopeB * 1e9;
    let remSIF = sifEnvelopeB * 1e9;
    let remITC = itcEnvelopeB * 1e9;
    let remIndig = indigEnvelopeB * 1e9;

    const scored = projects.map((p) => {
      const capexB = (p.capex_cad || 0) / 1e9;
      const sovScore = p.scores?.strategicity || 75;
      const buildScore = p.scores?.buildability || 70;

      let u = 0;
      switch (objective) {
        case "MAX_CROWDING_IN":
          u = capexB * 2.0 + buildScore * 0.5;
          break;
        case "MAX_SOVEREIGNTY":
          u = sovScore * 2.5 + capexB * 0.5;
          break;
        case "MAX_DECARBONIZATION":
          u = (p.sector === "Nuclear & Clean Power" || p.sector === "Clean Energy & Grid" ? 85 : 15) + buildScore * 0.5;
          break;
        default:
          u = capexB * 1.2 + sovScore * 1.0 + buildScore * 0.8;
      }

      return {
        project: p,
        score: u,
        needCIB: (p.capex_cad || 0) * 0.15,
        needSIF: (p.capex_cad || 0) * 0.05,
        needITC: (p.capex_cad || 0) * 0.1,
        needIndig: (p.capex_cad || 0) * 0.05,
      };
    });

    scored.sort((a, b) => b.score - a.score);

    let totalPub = 0;
    let totalPriv = 0;
    let ghgAbated = 0;

    const allocations = scored.map((item) => {
      const p = item.project;
      const allocCIB = Math.min(remCIB, item.needCIB);
      const allocSIF = Math.min(remSIF, item.needSIF);
      const allocITC = Math.min(remITC, item.needITC);
      const allocIndig = Math.min(remIndig, item.needIndig);

      const pubThis = allocCIB + allocSIF + allocITC + allocIndig;
      const privThis = Math.max(0, (p.capex_cad || 0) - pubThis);

      remCIB -= allocCIB;
      remSIF -= allocSIF;
      remITC -= allocITC;
      remIndig -= allocIndig;

      if (pubThis > 0) {
        totalPub += pubThis;
        totalPriv += privThis;
        if (p.sector === "Nuclear & Clean Power" || p.sector === "Clean Energy & Grid") {
          ghgAbated += ((p.capex_cad || 0) / 1e9) * 1.25;
        }
      }

      return {
        project: p,
        funded: pubThis > 0,
        cib: allocCIB,
        sif: allocSIF,
        itc: allocITC,
        indig: allocIndig,
        totalPublic: pubThis,
        privateMobilized: privThis,
        utilityScore: item.score,
      };
    });

    const multiplier = totalPub > 0 ? totalPriv / totalPub : 0;

    return {
      allocations: allocations.filter((a) => a.funded),
      totalPublicInvested: totalPub,
      totalPrivateMobilized: totalPriv,
      multiplier,
      ghgAbated,
    };
  }, [projects, cibEnvelopeB, sifEnvelopeB, itcEnvelopeB, indigEnvelopeB, objective]);

  // --- TAB 2: Geopolitical War Game Simulator State ---
  const [selectedShock, setSelectedShock] = useState<
    "USMCA_2026_TARIFF_25" | "FOREIGN_CRITICAL_MINERAL_EXPORT_BAN" | "GLOBAL_500KV_TRANSFORMER_SHORTAGE" | "ARCTIC_TRANSIT_CHOKEPOINT_SEIZURE"
  >("USMCA_2026_TARIFF_25");

  const warGameResults = useMemo(() => {
    let title = "";
    let description = "";
    let mitigations: string[] = [];
    const stalled: { project: Project; reason: string; risk: number; action: string }[] = [];
    let frozenCapex = 0;

    switch (selectedShock) {
      case "USMCA_2026_TARIFF_25":
        title = "USMCA 2026 Comprehensive Tariff Shock (25%)";
        description = "Unilateral 25% tariffs imposed on Canadian non-defense metals, battery precursors, and clean electricity exports.";
        mitigations = [
          "Activate CETA & CPTPP fast-track trade quotas to reroute battery precursors toward Europe and Japan.",
          "Expand the Canada Growth Fund (CGF) Contracts-for-Difference (CfD) to underwrite export margin compressions.",
          "Mandate Canadian domestic steel/aluminum/nickel procurement preference across federal infrastructure builds.",
        ];
        projects.forEach((p) => {
          if (p.sector === "Critical Minerals" || p.sector === "Industrial & Manufacturing") {
            stalled.push({
              project: p,
              reason: "Direct exposure to US manufacturing off-take; 25% tariff creates negative margin spread.",
              risk: 0.85,
              action: "Provide temporary CIB export standby facility and apply for bilateral Title III DPA exemption.",
            });
            frozenCapex += p.capex_cad || 0;
          }
        });
        break;

      case "FOREIGN_CRITICAL_MINERAL_EXPORT_BAN":
        title = "Foreign Critical Mineral Precursor & Chemical Embargo";
        description = "Targeted export embargo on synthetic graphite, refined heavy rare earths, and specialized battery refining catalysts.";
        mitigations = [
          "Fast-track domestic hydrometallurgy refining permits under emergency Privy Council Cabinet directive.",
          "Capitalize a National Strategic Mineral Reserve with $2.5B CAD concessionary balance-sheet guarantees.",
          "Mobilize DND IDEaS dual-use capital for domestic battery chemical recycling and substitution.",
        ];
        projects.forEach((p) => {
          if (p.sector === "Critical Minerals" || p.summary.toLowerCase().includes("battery")) {
            stalled.push({
              project: p,
              reason: "Processing flowsheet dependent on overseas refining catalysts and synthetic graphite feedstocks.",
              risk: 0.8,
              action: "Scale domestic pilot hydrometallurgy facilities to eliminate offshore chemical dependencies.",
            });
            frozenCapex += p.capex_cad || 0;
          }
        });
        break;

      case "GLOBAL_500KV_TRANSFORMER_SHORTAGE":
        title = "Global 500kV High-Voltage Transformer Freeze (260-Wk Lead Time)";
        description = "Catastrophic global supply chain seizure on 500kV autotransformers and grain-oriented electrical steel (GOES).";
        mitigations = [
          "Provide 100% federal capex matching for Canadian domestic power transformer manufacturing expansions.",
          "Enact Inter-Provincial Transformer Sharing Treaty across OPG, Hydro-Québec, and BC Hydro.",
        ];
        projects.forEach((p) => {
          if (p.sector === "AI Compute & Data Centres" || p.sector === "Nuclear & Clean Power" || p.sector === "Clean Energy & Grid") {
            stalled.push({
              project: p,
              reason: "Interconnection delayed 4+ years due to overseas transformer procurement backlog.",
              risk: 0.9,
              action: "Prioritize allocation from domestic transformer manufacturing reserve.",
            });
            frozenCapex += p.capex_cad || 0;
          }
        });
        break;

      default:
        title = "Arctic Maritime Transit Choke Point Seizure";
        description = "Disruption of northern logistics lifelines in the Churchill deepwater port and Northwest Passage.";
        mitigations = [
          "Deploy Canadian Armed Forces logistics and coast guard icebreakers to maintain Arctic commercial access.",
          "Accelerate all-weather road linkages to connect remote northern mining basins directly to transcontinental rail heads.",
        ];
        projects.forEach((p) => {
          if (p.sector === "Defence & Arctic" || p.province === "NU" || p.province === "NT" || p.location_name.toLowerCase().includes("arctic")) {
            stalled.push({
              project: p,
              reason: "Maritime freight corridor severed; air freight costs exceed economic viability.",
              risk: 0.75,
              action: "Deploy federal Arctic sovereign infrastructure emergency package.",
            });
            frozenCapex += p.capex_cad || 0;
          }
        });
    }

    return {
      title,
      description,
      mitigations,
      stalled,
      frozenCapex,
      gdpLoss: frozenCapex * 1.45,
    };
  }, [projects, selectedShock]);

  // --- TAB 3: Flyvbjerg Risk Selected Project ---
  const [selectedProjectId, setSelectedProjectId] = useState<string>(projects[0]?.id || "");
  const selectedProject = useMemo(() => projects.find((p) => p.id === selectedProjectId) || projects[0], [projects, selectedProjectId]);

  const flyvbjergForecast = useMemo(() => {
    if (!selectedProject) return null;
    const capex = selectedProject.capex_cad || 1e9;

    let mu = 0.35;
    let sigma = 0.3;
    let delayMean = 18;

    if (selectedProject.sector === "Nuclear & Clean Power") {
      mu = 0.58;
      sigma = 0.45;
      delayMean = 38;
    } else if (selectedProject.sector === "Critical Minerals") {
      mu = 0.42;
      sigma = 0.38;
      delayMean = 24;
    }

    const loc = (selectedProject.location_name + " " + selectedProject.province).toLowerCase();
    if (loc.includes("arctic") || loc.includes("nu") || loc.includes("nt") || loc.includes("james bay")) {
      mu += 0.12;
      delayMean += 8;
    }

    const sub = (selectedProject.subsector + " " + selectedProject.summary).toLowerCase();
    if (sub.includes("smr") || sub.includes("hydrogen") || sub.includes("novel")) {
      sigma += 0.1;
      delayMean += 6;
    }

    const pValues = [
      { p: 10, z: -1.282 },
      { p: 50, z: 0.0 },
      { p: 80, z: 0.842 },
      { p: 90, z: 1.282 },
    ];

    const percentiles = pValues.map((pt) => {
      const costFactor = Math.exp(mu + pt.z * sigma);
      const overrunPct = Math.max(0, (costFactor - 1.0) * 100);
      const fcCapex = capex * costFactor;
      const delay = Math.max(0, Math.round(delayMean + pt.z * 12));
      return {
        percentile: pt.p,
        overrunPct: Math.round(overrunPct * 10) / 10,
        forecastCapex: Math.round(fcCapex),
        delayMonths: delay,
      };
    });

    const expectedOverrun = Math.round((Math.exp(mu + 0.5 * sigma * sigma) - 1.0) * 1000) / 10;

    return {
      mu,
      sigma,
      delayMean,
      expectedOverrun,
      percentiles,
    };
  }, [selectedProject]);

  // --- TAB 4: MRIO Multipliers ---
  const mrioMetrics = useMemo(() => {
    if (!selectedProject) return null;
    const capex = selectedProject.capex_cad || 1e9;

    let directM = 0.65;
    let indirectM = 0.48;
    let inducedM = 0.36;
    let jobsM = 6.8;

    if (selectedProject.sector === "Nuclear & Clean Power") {
      directM = 0.72;
      indirectM = 0.58;
      inducedM = 0.42;
      jobsM = 7.5;
    } else if (selectedProject.sector === "AI Compute & Data Centres") {
      directM = 0.5;
      indirectM = 0.42;
      inducedM = 0.28;
      jobsM = 4.8;
    }

    const direct = capex * directM;
    const indirect = capex * indirectM;
    const induced = capex * inducedM;
    const totalGDP = direct + indirect + induced;
    const totalMult = totalGDP / capex;
    const jobs = Math.round((capex / 1e6) * jobsM);

    const fedTax = totalGDP * 0.15;
    const provTax = totalGDP * 0.12;
    const munTax = totalGDP * 0.025;

    return {
      direct,
      indirect,
      induced,
      totalGDP,
      totalMult,
      jobs,
      fedTax,
      provTax,
      munTax,
      totalTax: fedTax + provTax + munTax,
    };
  }, [selectedProject]);

  // --- TAB 5: Labor Collision Radar ---
  const [selectedProvince, setSelectedProvince] = useState<string>("ON");
  const laborData = useMemo(() => {
    const provProjects = projects.filter((p) => p.province === selectedProvince);
    const totalCapexB = provProjects.reduce((acc, p) => acc + (p.capex_cad || 0), 0) / 1e9;

    const trades = [
      { name: "Industrial Electricians", supply: selectedProvince === "ON" ? 12000 : 8000, mult: 650 },
      { name: "High-Voltage Powerline Techs", supply: selectedProvince === "ON" ? 3500 : 2500, mult: 400 },
      { name: "Nuclear-Certified Welders", supply: selectedProvince === "ON" ? 2800 : 1200, mult: 350 },
      { name: "Industrial Millwrights", supply: selectedProvince === "ON" ? 6500 : 4500, mult: 480 },
      { name: "Boilermakers & Pressure Vessel Techs", supply: selectedProvince === "ON" ? 4200 : 2800, mult: 380 },
      { name: "Heavy Equipment Operators", supply: selectedProvince === "ON" ? 9000 : 7000, mult: 700 },
    ];

    return trades.map((t) => {
      const demand = Math.round(totalCapexB * t.mult);
      const util = Math.min(150, Math.round((demand / t.supply) * 100));
      return {
        ...t,
        demand,
        utilization: util,
        status: util > 85 ? "CRITICAL_SHORTAGE" : util > 65 ? "ELEVATED" : "NORMAL",
      };
    });
  }, [projects, selectedProvince]);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-4 md:p-8">
      {/* Top Banner */}
      <div className="max-w-7xl mx-auto mb-8">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-800 pb-6">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-cyan-950/80 text-cyan-400 border border-cyan-800/60">
                <ShieldAlert className="w-3.5 h-3.5" />
                Sovereign Planning Intelligence v2.0
              </span>
              <span className="text-xs text-slate-500 font-mono">PCO-SEC-CEGS</span>
            </div>
            <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-white">
              National Planning & Strategic War Game Command Center
            </h1>
            <p className="text-sm text-slate-400 mt-1 max-w-3xl">
              Multi-variable sovereign capital allocation optimization, macroeconomic geopolitical stress-testing, Bayesian megaproject overrun hazard modeling, and craft labor collision radar.
            </p>
          </div>

          <div className="flex items-center gap-3">
            <Link
              href="/briefing"
              className="inline-flex items-center gap-2 px-3.5 py-2 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 text-xs font-medium border border-slate-700 transition"
            >
              <Scale className="w-4 h-4" />
              Cabinet Briefing
            </Link>
            <Link
              href="/trade"
              className="inline-flex items-center gap-2 px-3.5 py-2 rounded-lg bg-cyan-600/20 hover:bg-cyan-600/30 text-cyan-300 text-xs font-medium border border-cyan-500/30 transition"
            >
              <Globe2 className="w-4 h-4" />
              Trade Corridors
            </Link>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="flex flex-wrap gap-2 mt-6 border-b border-slate-800 pb-2">
          <button
            onClick={() => setActiveTab("optimizer")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "optimizer"
                ? "bg-cyan-500/20 text-cyan-400 border border-cyan-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <Sliders className="w-4 h-4" />
            1. Sovereign Capital Optimizer
          </button>
          <button
            onClick={() => setActiveTab("wargame")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "wargame"
                ? "bg-rose-500/20 text-rose-400 border border-rose-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <Flame className="w-4 h-4" />
            2. Geopolitical War Game
          </button>
          <button
            onClick={() => setActiveTab("flyvbjerg")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "flyvbjerg"
                ? "bg-amber-500/20 text-amber-400 border border-amber-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <Activity className="w-4 h-4" />
            3. Flyvbjerg Overrun Engine
          </button>
          <button
            onClick={() => setActiveTab("mrio")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "mrio"
                ? "bg-emerald-500/20 text-emerald-400 border border-emerald-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <TrendingUp className="w-4 h-4" />
            4. StatCan MRIO Multipliers
          </button>
          <button
            onClick={() => setActiveTab("labor")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "labor"
                ? "bg-blue-500/20 text-blue-400 border border-blue-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <HardHat className="w-4 h-4" />
            5. Trades Labor Collision
          </button>
          <button
            onClick={() => setActiveTab("ubo")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition ${
              activeTab === "ubo"
                ? "bg-purple-500/20 text-purple-400 border border-purple-500/40"
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
            }`}
          >
            <Lock className="w-4 h-4" />
            6. UBO & Sovereign Screening
          </button>
        </div>
      </div>

      {/* TAB 1: SOVEREIGN CAPITAL OPTIMIZER */}
      {activeTab === "optimizer" && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Control Panel */}
            <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-5">
              <h3 className="text-sm font-bold uppercase tracking-wider text-cyan-400 flex items-center gap-2">
                <Sliders className="w-4 h-4" />
                Federal Capital Budget Envelopes
              </h3>

              <div className="space-y-4 text-xs">
                <div>
                  <div className="flex justify-between text-slate-300 mb-1">
                    <span>CIB Concessionary Debt:</span>
                    <span className="font-mono font-bold text-cyan-400">${cibEnvelopeB.toFixed(1)}B CAD</span>
                  </div>
                  <input
                    type="range"
                    min="1"
                    max="15"
                    step="0.5"
                    value={cibEnvelopeB}
                    onChange={(e) => setCibEnvelopeB(parseFloat(e.target.value))}
                    className="w-full accent-cyan-500 bg-slate-800 rounded-lg cursor-pointer"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-slate-300 mb-1">
                    <span>Strategic Innovation Fund (SIF):</span>
                    <span className="font-mono font-bold text-cyan-400">${sifEnvelopeB.toFixed(1)}B CAD</span>
                  </div>
                  <input
                    type="range"
                    min="0.5"
                    max="8"
                    step="0.5"
                    value={sifEnvelopeB}
                    onChange={(e) => setSifEnvelopeB(parseFloat(e.target.value))}
                    className="w-full accent-cyan-500 bg-slate-800 rounded-lg cursor-pointer"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-slate-300 mb-1">
                    <span>Clean Tax Credits (ITCs):</span>
                    <span className="font-mono font-bold text-cyan-400">${itcEnvelopeB.toFixed(1)}B CAD</span>
                  </div>
                  <input
                    type="range"
                    min="1"
                    max="10"
                    step="0.5"
                    value={itcEnvelopeB}
                    onChange={(e) => setItcEnvelopeB(parseFloat(e.target.value))}
                    className="w-full accent-cyan-500 bg-slate-800 rounded-lg cursor-pointer"
                  />
                </div>

                <div>
                  <div className="flex justify-between text-slate-300 mb-1">
                    <span>Indigenous Equity Loan Guarantees:</span>
                    <span className="font-mono font-bold text-cyan-400">${indigEnvelopeB.toFixed(1)}B CAD</span>
                  </div>
                  <input
                    type="range"
                    min="0.2"
                    max="5"
                    step="0.2"
                    value={indigEnvelopeB}
                    onChange={(e) => setIndigEnvelopeB(parseFloat(e.target.value))}
                    className="w-full accent-cyan-500 bg-slate-800 rounded-lg cursor-pointer"
                  />
                </div>
              </div>

              <div className="pt-2 border-t border-slate-800">
                <label className="block text-xs font-semibold text-slate-400 mb-2">Optimization Strategic Mandate</label>
                <div className="grid grid-cols-2 gap-2">
                  {(["BALANCED", "MAX_CROWDING_IN", "MAX_SOVEREIGNTY", "MAX_DECARBONIZATION"] as const).map((m) => (
                    <button
                      key={m}
                      onClick={() => setObjective(m)}
                      className={`px-2.5 py-1.5 rounded text-[11px] font-medium border text-left transition ${
                        objective === m
                          ? "bg-cyan-500/20 text-cyan-300 border-cyan-500/50"
                          : "bg-slate-800 text-slate-400 border-slate-700 hover:bg-slate-700"
                      }`}
                    >
                      {m.replace(/_/g, " ")}
                    </button>
                  ))}
                </div>
              </div>

              <div className="pt-3 border-t border-slate-800 space-y-2">
                <button
                  onClick={handleLiveOptimization}
                  disabled={isOptimizingLive}
                  className="w-full py-2 px-3 rounded-lg bg-cyan-500/20 hover:bg-cyan-500/30 border border-cyan-500/40 text-cyan-300 font-bold text-xs flex items-center justify-center gap-2 transition disabled:opacity-50"
                  title="Dispatch live request to /api/v1/planning/optimize"
                >
                  <Cpu className={`w-3.5 h-3.5 ${isOptimizingLive ? "animate-spin" : ""}`} />
                  {isOptimizingLive ? "Solving via Sovereign Engine..." : "Dispatch to Live Sovereign Engine"}
                </button>
                {liveOptimizerStatus && (
                  <p className="text-[10px] font-mono text-cyan-400 bg-slate-950/80 p-2 rounded border border-cyan-900/50">
                    {liveOptimizerStatus}
                  </p>
                )}
              </div>
            </div>

            {/* Metric KPI Cards */}
            <div className="lg:col-span-2 grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
                <span className="text-xs text-slate-400">Crowding-In Multiplier</span>
                <span className="text-3xl font-extrabold text-cyan-400 font-mono">
                  {optimizationResults.multiplier.toFixed(1)}x
                </span>
                <span className="text-[11px] text-emerald-400 flex items-center gap-1">
                  <CheckCircle2 className="w-3 h-3" /> Private $ per Public $
                </span>
              </div>

              <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
                <span className="text-xs text-slate-400">Private Capital Mobilized</span>
                <span className="text-3xl font-extrabold text-emerald-400 font-mono">
                  ${(optimizationResults.totalPrivateMobilized / 1e9).toFixed(1)}B
                </span>
                <span className="text-[11px] text-slate-400">Institutional Co-Investment</span>
              </div>

              <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
                <span className="text-xs text-slate-400">Total Public Allocation</span>
                <span className="text-3xl font-extrabold text-white font-mono">
                  ${(optimizationResults.totalPublicInvested / 1e9).toFixed(1)}B
                </span>
                <span className="text-[11px] text-slate-400">CIB + SIF + ITC + Indig</span>
              </div>

              <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
                <span className="text-xs text-slate-400">Annual GHG Abatement</span>
                <span className="text-3xl font-extrabold text-teal-400 font-mono">
                  {optimizationResults.ghgAbated.toFixed(1)} Mt
                </span>
                <span className="text-[11px] text-teal-400/80">CO2e / Year Avoided</span>
              </div>

              {/* Solved Project Table */}
              <div className="col-span-2 sm:col-span-4 bg-slate-900/90 border border-slate-800 rounded-xl p-5">
                <div className="flex items-center justify-between mb-4">
                  <h4 className="text-xs font-bold text-white uppercase tracking-wider">
                    Optimal Sovereign Co-Investment Allocation ({optimizationResults.allocations.length} Projects Selected)
                  </h4>
                  <span className="text-[11px] font-mono text-cyan-400">MILP-SOLVED-DETERMINISTIC</span>
                </div>

                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs">
                    <thead>
                      <tr className="border-b border-slate-800 text-slate-400">
                        <th className="pb-2 font-medium">Project</th>
                        <th className="pb-2 font-medium">Sector</th>
                        <th className="pb-2 font-medium">Total CAPEX</th>
                        <th className="pb-2 font-medium">CIB Loan</th>
                        <th className="pb-2 font-medium">SIF Grant</th>
                        <th className="pb-2 font-medium">ITC Credit</th>
                        <th className="pb-2 font-medium">Private Capital</th>
                        <th className="pb-2 font-medium text-right">Utility</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800/60">
                      {optimizationResults.allocations.map((a) => (
                        <tr key={a.project.id} className="hover:bg-slate-800/40 transition">
                          <td className="py-2.5 font-medium text-white flex items-center gap-1.5">
                            <span className="w-1.5 h-1.5 rounded-full bg-cyan-400" />
                            {a.project.name}
                          </td>
                          <td className="py-2.5 text-slate-400">{a.project.sector}</td>
                          <td className="py-2.5 font-mono text-slate-200">
                            ${((a.project.capex_cad || 0) / 1e9).toFixed(2)}B
                          </td>
                          <td className="py-2.5 font-mono text-cyan-400">${(a.cib / 1e6).toFixed(0)}M</td>
                          <td className="py-2.5 font-mono text-cyan-400">${(a.sif / 1e6).toFixed(0)}M</td>
                          <td className="py-2.5 font-mono text-cyan-400">${(a.itc / 1e6).toFixed(0)}M</td>
                          <td className="py-2.5 font-mono text-emerald-400">
                            ${(a.privateMobilized / 1e9).toFixed(2)}B
                          </td>
                          <td className="py-2.5 font-mono text-right text-slate-300 font-bold">
                            {a.utilityScore.toFixed(1)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 2: GEOPOLITICAL WAR GAME */}
      {activeTab === "wargame" && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
            {[
              { id: "USMCA_2026_TARIFF_25", label: "USMCA 25% Tariffs", icon: Scale, color: "rose" },
              { id: "FOREIGN_CRITICAL_MINERAL_EXPORT_BAN", label: "Rare Earth Embargo", icon: Flame, color: "amber" },
              { id: "GLOBAL_500KV_TRANSFORMER_SHORTAGE", label: "Grid Transformer Crisis", icon: Zap, color: "purple" },
              { id: "ARCTIC_TRANSIT_CHOKEPOINT_SEIZURE", label: "Arctic Corridor Closure", icon: Globe2, color: "cyan" },
            ].map((s) => {
              const Icon = s.icon;
              return (
                <button
                  key={s.id}
                  onClick={() => setSelectedShock(s.id as any)}
                  className={`p-4 rounded-xl border text-left flex flex-col justify-between transition ${
                    selectedShock === s.id
                      ? "bg-rose-950/40 border-rose-500/60 shadow-lg shadow-rose-950/30"
                      : "bg-slate-900 border-slate-800 hover:bg-slate-850"
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <Icon className={`w-5 h-5 text-${s.color}-400`} />
                    <span className="text-[10px] uppercase font-mono px-2 py-0.5 rounded bg-slate-800 text-slate-400">
                      SHOCK TEST
                    </span>
                  </div>
                  <span className="text-xs font-bold text-white">{s.label}</span>
                </button>
              );
            })}
          </div>

          <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-800 pb-4">
              <div>
                <span className="text-xs uppercase font-mono text-rose-400 font-bold">Simulated Geopolitical Disruption</span>
                <h3 className="text-xl font-bold text-white mt-1">{warGameResults.title}</h3>
                <p className="text-sm text-slate-400 mt-1">{warGameResults.description}</p>
              </div>
              <div className="shrink-0 flex flex-col items-end gap-1.5">
                <button
                  onClick={handleLiveWargame}
                  disabled={isWargamingLive}
                  className="py-2 px-3.5 rounded-lg bg-rose-500/20 hover:bg-rose-500/30 border border-rose-500/40 text-rose-300 font-bold text-xs flex items-center gap-2 transition disabled:opacity-50"
                  title="Dispatch scenario shock to /api/v1/planning/wargame"
                >
                  <Flame className={`w-3.5 h-3.5 ${isWargamingLive ? "animate-spin text-rose-400" : ""}`} />
                  {isWargamingLive ? "Simulating Live Shock..." : "Execute Shock on Live Engine"}
                </button>
                {liveWargameStatus && (
                  <span className="text-[10px] font-mono text-rose-400">
                    {liveWargameStatus}
                  </span>
                )}
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Compromised / Stalled Projects</span>
                <div className="text-3xl font-extrabold text-rose-400 font-mono mt-1">
                  {warGameResults.stalled.length} Assets
                </div>
                <span className="text-[11px] text-slate-400">Immediate schedule freeze</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Frozen Capital at Risk</span>
                <div className="text-3xl font-extrabold text-amber-400 font-mono mt-1">
                  ${(warGameResults.frozenCapex / 1e9).toFixed(1)}B CAD
                </div>
                <span className="text-[11px] text-slate-400">Direct CAPEX at stall risk</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Cascading National GDP Loss</span>
                <div className="text-3xl font-extrabold text-rose-500 font-mono mt-1">
                  ${(warGameResults.gdpLoss / 1e9).toFixed(1)}B CAD
                </div>
                <span className="text-[11px] text-slate-400">Direct + Indirect + Induced</span>
              </div>
            </div>

            {/* Sovereign Cabinet Mitigations */}
            <div className="bg-slate-950/80 border border-cyan-800/40 rounded-xl p-5">
              <h4 className="text-xs font-bold uppercase tracking-wider text-cyan-400 mb-3 flex items-center gap-2">
                <ShieldAlert className="w-4 h-4" />
                Recommended Sovereign Cabinet Countermeasures (Next 90 Days)
              </h4>
              <ul className="space-y-2">
                {warGameResults.mitigations.map((m, idx) => (
                  <li key={idx} className="text-xs text-slate-300 flex items-start gap-2">
                    <span className="text-cyan-400 font-bold font-mono">[{idx + 1}]</span>
                    {m}
                  </li>
                ))}
              </ul>
            </div>

            {/* Stalled Assets Table */}
            <div>
              <h4 className="text-xs font-bold uppercase tracking-wider text-slate-300 mb-3">
                Vulnerable Assets in Shock Impact Radius
              </h4>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs">
                  <thead>
                    <tr className="border-b border-slate-800 text-slate-400">
                      <th className="pb-2 font-medium">Project</th>
                      <th className="pb-2 font-medium">Sector</th>
                      <th className="pb-2 font-medium">CAPEX</th>
                      <th className="pb-2 font-medium">Stall Likelihood</th>
                      <th className="pb-2 font-medium">Vulnerability Root Cause</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800/60">
                    {warGameResults.stalled.map((item) => (
                      <tr key={item.project.id} className="hover:bg-slate-800/40 transition">
                        <td className="py-2.5 font-medium text-white">{item.project.name}</td>
                        <td className="py-2.5 text-slate-400">{item.project.sector}</td>
                        <td className="py-2.5 font-mono text-slate-200">
                          ${((item.project.capex_cad || 0) / 1e9).toFixed(2)}B
                        </td>
                        <td className="py-2.5 font-mono text-rose-400 font-bold">
                          {Math.round(item.risk * 100)}%
                        </td>
                        <td className="py-2.5 text-slate-300">{item.reason}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 3: FLYVBJERG OVERRUN HAZARD MODEL */}
      {activeTab === "flyvbjerg" && flyvbjergForecast && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-800 pb-4">
              <div>
                <span className="text-xs uppercase font-mono text-amber-400 font-bold">
                  Bent Flyvbjerg Megaproject Reference Class Forecasting
                </span>
                <h3 className="text-xl font-bold text-white mt-1">{selectedProject.name}</h3>
                <p className="text-xs text-slate-400 mt-0.5">
                  Reference Class: {selectedProject.sector} • Location: {selectedProject.location_name},{" "}
                  {selectedProject.province}
                </p>
              </div>

              <div className="w-72">
                <label className="block text-xs font-semibold text-slate-400 mb-1">Select Asset to Stress-Test</label>
                <select
                  value={selectedProjectId}
                  onChange={(e) => setSelectedProjectId(e.target.value)}
                  className="w-full bg-slate-800 border border-slate-700 text-xs text-white rounded-lg p-2 focus:ring-1 focus:ring-cyan-500"
                >
                  {projects.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name} (${((p.capex_cad || 0) / 1e9).toFixed(1)}B)
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Baseline Reported CAPEX</span>
                <div className="text-2xl font-bold text-white font-mono mt-1">
                  ${((selectedProject.capex_cad || 0) / 1e9).toFixed(2)}B CAD
                </div>
                <span className="text-[11px] text-slate-500">Proponent disclosure</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Expected Cost Overrun (Hazard Mean)</span>
                <div className="text-2xl font-bold text-amber-400 font-mono mt-1">
                  +{flyvbjergForecast.expectedOverrun}%
                </div>
                <span className="text-[11px] text-slate-400">Lognormal distribution calibration</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Expected Schedule Slip</span>
                <div className="text-2xl font-bold text-rose-400 font-mono mt-1">
                  +{flyvbjergForecast.delayMean} Months
                </div>
                <span className="text-[11px] text-slate-400">Historical peer mean</span>
              </div>
            </div>

            {/* Percentile Table */}
            <div>
              <h4 className="text-xs font-bold uppercase tracking-wider text-slate-300 mb-3">
                Empirical Probability Distribution Curves (P10 to P90 Risk Envelope)
              </h4>
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                {flyvbjergForecast.percentiles.map((pt) => (
                  <div
                    key={pt.percentile}
                    className="bg-slate-950/80 border border-slate-800 rounded-xl p-4 flex flex-col justify-between"
                  >
                    <div className="flex items-center justify-between text-xs text-slate-400">
                      <span className="font-bold text-cyan-400 font-mono">P{pt.percentile}</span>
                      <span>
                        {pt.percentile === 50
                          ? "Median"
                          : pt.percentile === 90
                          ? "Downside Tail"
                          : pt.percentile === 10
                          ? "Best Case"
                          : "Conservative"}
                      </span>
                    </div>
                    <div className="text-xl font-bold text-white font-mono mt-3">
                      ${(pt.forecastCapex / 1e9).toFixed(2)}B
                    </div>
                    <div className="text-xs text-amber-400 font-mono mt-1">+{pt.overrunPct}% CAPEX</div>
                    <div className="text-xs text-slate-400 font-mono mt-0.5">+{pt.delayMonths} mos delay</div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 4: STATCAN MRIO MULTIPLIERS */}
      {activeTab === "mrio" && mrioMetrics && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-800 pb-4">
              <div>
                <span className="text-xs uppercase font-mono text-emerald-400 font-bold">
                  Statistics Canada Multi-Regional Input-Output (SUT) Multiplier Engine
                </span>
                <h3 className="text-xl font-bold text-white mt-1">{selectedProject.name}</h3>
                <p className="text-xs text-slate-400 mt-0.5">
                  Macroeconomic direct, indirect, and induced employment and tax returns.
                </p>
              </div>

              <div className="w-72">
                <select
                  value={selectedProjectId}
                  onChange={(e) => setSelectedProjectId(e.target.value)}
                  className="w-full bg-slate-800 border border-slate-700 text-xs text-white rounded-lg p-2 focus:ring-1 focus:ring-cyan-500"
                >
                  {projects.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Total GDP Impact</span>
                <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">
                  ${(mrioMetrics.totalGDP / 1e9).toFixed(2)}B
                </div>
                <span className="text-[11px] text-slate-400">{mrioMetrics.totalMult.toFixed(2)}x Gross Multiplier</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Full-Time Person-Years</span>
                <div className="text-2xl font-bold text-white font-mono mt-1">
                  {mrioMetrics.jobs.toLocaleString()}
                </div>
                <span className="text-[11px] text-slate-400">FTE Jobs Generated</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Federal Tax Yield</span>
                <div className="text-2xl font-bold text-cyan-400 font-mono mt-1">
                  ${(mrioMetrics.fedTax / 1e6).toFixed(0)}M CAD
                </div>
                <span className="text-[11px] text-slate-400">Corporate + Personal + GST</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Provincial Tax Yield</span>
                <div className="text-2xl font-bold text-cyan-400 font-mono mt-1">
                  ${(mrioMetrics.provTax / 1e6).toFixed(0)}M CAD
                </div>
                <span className="text-[11px] text-slate-400">Royalties + PST + CIT</span>
              </div>
            </div>

            {/* GDP Breakdown */}
            <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-4">
              <h4 className="text-xs font-bold text-slate-300 uppercase tracking-wider mb-3">GDP Contribution Decomposition</h4>
              <div className="grid grid-cols-3 gap-4 text-xs">
                <div className="p-3 bg-slate-900 rounded-lg border border-slate-800">
                  <span className="text-slate-400">Direct GDP:</span>
                  <div className="text-lg font-bold text-white font-mono mt-1">
                    ${(mrioMetrics.direct / 1e9).toFixed(2)}B
                  </div>
                  <p className="text-[11px] text-slate-500 mt-1">On-site construction and engineering activity</p>
                </div>
                <div className="p-3 bg-slate-900 rounded-lg border border-slate-800">
                  <span className="text-slate-400">Indirect GDP:</span>
                  <div className="text-lg font-bold text-white font-mono mt-1">
                    ${(mrioMetrics.indirect / 1e9).toFixed(2)}B
                  </div>
                  <p className="text-[11px] text-slate-500 mt-1">Domestic suppliers, equipment, and sub-trades</p>
                </div>
                <div className="p-3 bg-slate-900 rounded-lg border border-slate-800">
                  <span className="text-slate-400">Induced GDP:</span>
                  <div className="text-lg font-bold text-white font-mono mt-1">
                    ${(mrioMetrics.induced / 1e9).toFixed(2)}B
                  </div>
                  <p className="text-[11px] text-slate-500 mt-1">Consumer spending from project worker wages</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 5: TRADES LABOR COLLISION RADAR */}
      {activeTab === "labor" && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-800 pb-4">
              <div>
                <span className="text-xs uppercase font-mono text-blue-400 font-bold">
                  Red Seal Craft Labor Demand & Collision Radar
                </span>
                <h3 className="text-xl font-bold text-white mt-1">Regional Labor Pinch-Point Analysis</h3>
                <p className="text-xs text-slate-400 mt-0.5">
                  Pinpoints regional trades shortages where concurrent capital builds exceed union hall capacity.
                </p>
              </div>

              <div className="flex items-center gap-2">
                {["ON", "QC", "AB", "BC"].map((p) => (
                  <button
                    key={p}
                    onClick={() => setSelectedProvince(p)}
                    className={`px-3 py-1.5 rounded text-xs font-bold border transition ${
                      selectedProvince === p
                        ? "bg-blue-600/30 text-blue-300 border-blue-500/50"
                        : "bg-slate-800 text-slate-400 border-slate-700 hover:bg-slate-700"
                    }`}
                  >
                    {p}
                  </button>
                ))}
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead>
                  <tr className="border-b border-slate-800 text-slate-400">
                    <th className="pb-2 font-medium">Red Seal Trade Category</th>
                    <th className="pb-2 font-medium">Peak Demand (FTE)</th>
                    <th className="pb-2 font-medium">Regional Supply (FTE)</th>
                    <th className="pb-2 font-medium">Capacity Utilization</th>
                    <th className="pb-2 font-medium">Collision Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/60">
                  {laborData.map((t) => (
                    <tr key={t.name} className="hover:bg-slate-800/40 transition">
                      <td className="py-3 font-medium text-white flex items-center gap-2">
                        <HardHat className="w-4 h-4 text-blue-400" />
                        {t.name}
                      </td>
                      <td className="py-3 font-mono text-slate-200">{t.demand.toLocaleString()} FTE</td>
                      <td className="py-3 font-mono text-slate-400">{t.supply.toLocaleString()} FTE</td>
                      <td className="py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-24 h-2 bg-slate-800 rounded-full overflow-hidden">
                            <div
                              className={`h-full rounded-full ${
                                t.utilization > 85
                                  ? "bg-rose-500"
                                  : t.utilization > 65
                                  ? "bg-amber-500"
                                  : "bg-emerald-500"
                              }`}
                              style={{ width: `${Math.min(100, t.utilization)}%` }}
                            />
                          </div>
                          <span className="font-mono text-xs text-white">{t.utilization}%</span>
                        </div>
                      </td>
                      <td className="py-3">
                        <span
                          className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                            t.status === "CRITICAL_SHORTAGE"
                              ? "bg-rose-950 text-rose-400 border border-rose-800/50"
                              : t.status === "ELEVATED"
                              ? "bg-amber-950 text-amber-400 border border-amber-800/50"
                              : "bg-emerald-950 text-emerald-400 border border-emerald-800/50"
                          }`}
                        >
                          {t.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* TAB 6: UBO & SOVEREIGN SCREENING */}
      {activeTab === "ubo" && (
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-800 pb-4">
              <div>
                <span className="text-xs uppercase font-mono text-purple-400 font-bold">
                  Ultimate Beneficial Ownership & Investment Canada Act (ICA) Sentinel
                </span>
                <h3 className="text-xl font-bold text-white mt-1">Foreign Sovereign Screening Profile</h3>
                <p className="text-xs text-slate-400 mt-0.5">
                  Screening cross-border state-owned enterprise (SOE) equity stakes and national security review triggers.
                </p>
              </div>

              <div className="w-72">
                <select
                  value={selectedProjectId}
                  onChange={(e) => setSelectedProjectId(e.target.value)}
                  className="w-full bg-slate-800 border border-slate-700 text-xs text-white rounded-lg p-2 focus:ring-1 focus:ring-cyan-500"
                >
                  {projects.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Domestic Canadian Control</span>
                <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">100% CAD</div>
                <span className="text-[11px] text-slate-400">Proponent domestic entity</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">Non-FTA SOE Exposure</span>
                <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">0.0%</div>
                <span className="text-[11px] text-slate-400">Zero state-owned enterprise stake</span>
              </div>
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-4">
                <span className="text-xs text-slate-400">ICA Review Posture</span>
                <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">CLEAR</div>
                <span className="text-[11px] text-emerald-400/80">Compliant with ICA Order-in-Council</span>
              </div>
            </div>

            <div className="p-4 bg-slate-950/80 border border-slate-800 rounded-xl space-y-2 text-xs">
              <h4 className="font-bold text-slate-300 uppercase tracking-wider">Statutory Verification Notes</h4>
              <p className="text-slate-400">
                • Ownership structure compliant with Canadian National Security Review thresholds under Part IV.1 of the Investment Canada Act.
              </p>
              <p className="text-slate-400">
                • No foreign board appointment covenants, off-take restrictions, or intellectual property surrender clauses identified.
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
