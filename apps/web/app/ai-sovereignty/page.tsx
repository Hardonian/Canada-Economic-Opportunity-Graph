"use client";

import { useMemo, useState } from "react";
import {
  CheckCircle2,
  Cpu,
  Info,
  ShieldAlert,
  Sliders,
  XCircle,
  Zap,
  Server,
  Activity,
  Gauge,
  Sparkles,
  Lock,
} from "lucide-react";

type Criterion = {
  id: string;
  label: string;
  description: string;
  weight: number;
};

const criteria: Criterion[] = [
  { id: "data", label: "Canadian data residency", description: "Workload data is stored in Canada.", weight: 15 },
  { id: "compute", label: "Canadian compute residency", description: "The workload executes on infrastructure in Canada.", weight: 15 },
  { id: "control", label: "Operational control", description: "Documented Canadian control exists for privileged operations and keys.", weight: 20 },
  { id: "legal", label: "Legal exposure reviewed", description: "Qualified counsel has assessed applicable domestic and foreign disclosure obligations.", weight: 20 },
  { id: "continuity", label: "Continuity and portability", description: "Exit, recovery, and provider-portability controls have been tested.", weight: 15 },
  { id: "privacy", label: "Privacy controls assessed", description: "A current privacy impact and jurisdictional review exists.", weight: 10 },
  { id: "language", label: "Official-language needs tested", description: "Required English and French service quality has been evaluated.", weight: 5 },
];

interface GridOperatorProfile {
  id: string;
  name: string;
  province: string;
  operator: string;
  baseloadSource: string;
  cleanCapacityMW: number;
  availableAiHeadroomMW: number;
  gridCarbonIntensityGramsKWh: number; // gCO2eq/kWh
  cleanFlopsRatioExaGW: number; // ExaFLOPs per GW
}

const GRID_PROFILES: GridOperatorProfile[] = [
  {
    id: "ieso-on",
    name: "Ontario Nuclear & SMR Backbone",
    province: "Ontario",
    operator: "IESO",
    baseloadSource: "CANDU Nuclear & BWRX-300 SMRs",
    cleanCapacityMW: 3500,
    availableAiHeadroomMW: 850,
    gridCarbonIntensityGramsKWh: 28,
    cleanFlopsRatioExaGW: 7.00,
  },
  {
    id: "hq-qc",
    name: "Quebec Hydroelectric Hyperscale Megacampus",
    province: "Quebec",
    operator: "Hydro-Québec",
    baseloadSource: "Northern Reservoir & Run-of-River Hydro",
    cleanCapacityMW: 1900,
    availableAiHeadroomMW: 600,
    gridCarbonIntensityGramsKWh: 1.2,
    cleanFlopsRatioExaGW: 6.94,
  },
  {
    id: "bchydro-bc",
    name: "British Columbia Clean Hydro Dispatch",
    province: "British Columbia",
    operator: "BC Hydro",
    baseloadSource: "Site C & Peace River Hydroelectric",
    cleanCapacityMW: 1100,
    availableAiHeadroomMW: 450,
    gridCarbonIntensityGramsKWh: 9.5,
    cleanFlopsRatioExaGW: 6.95,
  },
  {
    id: "aeso-ab",
    name: "Alberta Industrial Heartland Cogeneration",
    province: "Alberta",
    operator: "AESO",
    baseloadSource: "Industrial Cogeneration + Carbon Capture (CCS)",
    cleanCapacityMW: 900,
    availableAiHeadroomMW: 350,
    gridCarbonIntensityGramsKWh: 85,
    cleanFlopsRatioExaGW: 6.46,
  },
];

export default function AISovereigntyPage() {
  // Checklist State
  const [answers, setAnswers] = useState<Record<string, boolean>>(() =>
    Object.fromEntries(criteria.map((criterion) => [criterion.id, false])),
  );

  // Compute Nexus State
  const [selectedGridId, setSelectedGridId] = useState<string>("ieso-on");
  const [requestedPowerMW, setRequestedPowerMW] = useState<number>(120);
  const [datacenterPUE, setDatacenterPUE] = useState<number>(1.18);
  const [gpuArch, setGpuArch] = useState<"B200" | "H100">("B200");

  const score = useMemo(
    () => criteria.reduce((total, criterion) => total + (answers[criterion.id] ? criterion.weight : 0), 0),
    [answers],
  );

  const selectedGrid = useMemo(
    () => GRID_PROFILES.find((g) => g.id === selectedGridId) || GRID_PROFILES[0],
    [selectedGridId],
  );

  // Compute Nexus Simulation Engine
  const computeModel = useMemo(() => {
    // Net IT Power dedicated to accelerator compute = Total Power / PUE
    const itPowerMW = requestedPowerMW / datacenterPUE;

    // GPU density estimation (assuming ~1.2 kW per GPU server board including cooling & network)
    const gpuTdpKW = gpuArch === "B200" ? 1.2 : 0.85;
    const totalAccelerators = Math.round((itPowerMW * 1000) / gpuTdpKW);

    // Compute throughput in ExaFLOPs (FP8 Tensor Core)
    // B200 ~4.5 PFLOPs FP8; H100 ~2.0 PFLOPs FP8
    const flopsPerGpuPflops = gpuArch === "B200" ? 4.5 : 2.0;
    const totalComputeExaFlops = (totalAccelerators * flopsPerGpuPflops) / 1000;

    // Headroom utilization percentage
    const headroomPct = Math.min(100, (requestedPowerMW / selectedGrid.availableAiHeadroomMW) * 100);

    // Clean energy factor & carbon footprint avoidance
    const isCleanBaseload = selectedGrid.gridCarbonIntensityGramsKWh < 40;
    const annualGWh = (requestedPowerMW * 8760) / 1000;
    const annualCarbonTonnes = (annualGWh * 1_000_000 * selectedGrid.gridCarbonIntensityGramsKWh) / 1e6;

    // Avoided carbon compared to coal/heavy gas baseline (~500g/kWh)
    const avoidedCarbonTonnes = Math.max(0, (annualGWh * 1_000_000 * (500 - selectedGrid.gridCarbonIntensityGramsKWh)) / 1e6);

    return {
      itPowerMW,
      totalAccelerators,
      totalComputeExaFlops,
      headroomPct,
      isCleanBaseload,
      annualGWh,
      annualCarbonTonnes,
      avoidedCarbonTonnes,
    };
  }, [requestedPowerMW, datacenterPUE, gpuArch, selectedGrid]);

  return (
    <div className="mx-auto max-w-6xl space-y-12 px-4 py-8 sm:px-6 lg:px-8">
      {/* Official Header */}
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
            <span className="font-bold text-white tracking-tight">ISED / NRCAN</span>
          </div>
          <span className="text-border">|</span>
          <Cpu className="h-3.5 w-3.5" />
          <span>SOVEREIGN AI & CLEAN BASELOAD NEXUS</span>
        </div>
        <h1 className="text-2xl font-black tracking-tight text-text-main sm:text-4xl">
          Sovereign AI Compute & <span className="text-aurora">Clean Baseload Power</span>
        </h1>
        <p className="mt-2 max-w-4xl text-sm leading-relaxed text-text-muted">
          Canada possesses one of the world&apos;s lowest-carbon, high-capacity electricity grids. Evaluate regional
          baseload headroom, clean FLOPs efficiency, and domestic data sovereignty controls across provincial utilities.
        </p>
      </header>

      {/* SECTION 1: Clean Baseload Power vs Sovereign AI Hyperscale Nexus */}
      <section className="space-y-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <div className="inline-flex items-center gap-2 text-xs font-mono font-bold text-sky-400 uppercase tracking-wider">
              <Zap className="h-4 w-4" />
              Pillar D: Clean FLOPs & Energy Sovereignty
            </div>
            <h2 className="text-xl font-black text-text-main mt-1">
              Hyperscale Power Allocation & Clean FLOPs Simulator
            </h2>
          </div>
          <div className="rounded-xl border border-sky-500/40 bg-sky-950/20 px-3 py-1.5 text-xs font-mono text-sky-300">
            IESO · Hydro-Québec · BC Hydro · AESO
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[1fr_1.1fr]">
          {/* Controls */}
          <div className="glass-card space-y-5 rounded-2xl border border-border/80 p-6">
            <h3 className="font-mono text-xs uppercase font-bold text-text-muted tracking-wider border-b border-border/60 pb-2 flex items-center gap-1.5">
              <Server className="h-3.5 w-3.5 text-aurora" /> Datacenter Interconnection Inputs
            </h3>

            <label className="block text-xs font-mono text-text-muted">
              Target Provincial Utility & Clean Grid Intertie
              <select
                value={selectedGridId}
                onChange={(e) => setSelectedGridId(e.target.value)}
                className="mt-1.5 w-full rounded-xl border border-border bg-[#040806] px-3 py-2 text-xs font-bold text-text-main outline-none focus:border-aurora"
              >
                {GRID_PROFILES.map((g) => (
                  <option key={g.id} value={g.id}>
                    {g.name} ({g.operator} — {g.availableAiHeadroomMW} MW Headroom)
                  </option>
                ))}
              </select>
            </label>

            <label className="block text-xs font-mono text-text-muted">
              <div className="flex justify-between">
                <span>Requested Campus Power Capacity</span>
                <strong className="text-aurora font-bold text-sm">{requestedPowerMW} MW</strong>
              </div>
              <input
                type="range"
                min="20"
                max={Math.min(500, selectedGrid.availableAiHeadroomMW)}
                step="10"
                value={requestedPowerMW}
                onChange={(e) => setRequestedPowerMW(Number(e.target.value))}
                className="mt-2 w-full accent-emerald-400"
              />
              <div className="flex justify-between text-[10px] text-text-subtle">
                <span>20 MW Edge Node</span>
                <span>{selectedGrid.availableAiHeadroomMW} MW Max Utility Headroom</span>
              </div>
            </label>

            <div className="grid grid-cols-2 gap-3">
              <label className="block text-xs font-mono text-text-muted">
                <span>Architecture Target</span>
                <select
                  value={gpuArch}
                  onChange={(e) => setGpuArch(e.target.value as "B200" | "H100")}
                  className="mt-1.5 w-full rounded-xl border border-border bg-[#040806] px-3 py-2 text-xs font-bold text-white outline-none focus:border-aurora"
                >
                  <option value="B200">Blackwell B200 (4.5 PF)</option>
                  <option value="H100">Hopper H100 (2.0 PF)</option>
                </select>
              </label>

              <label className="block text-xs font-mono text-text-muted">
                <span>Power Usage Effectiveness (PUE)</span>
                <input
                  type="number"
                  min="1.05"
                  max="1.50"
                  step="0.01"
                  value={datacenterPUE}
                  onChange={(e) => setDatacenterPUE(Number(e.target.value) || 1.15)}
                  className="mt-1.5 w-full rounded-xl border border-border bg-[#040806] px-3 py-2 text-xs font-bold text-white outline-none focus:border-aurora"
                />
              </label>
            </div>

            <div className="bg-[#0C1812] rounded-xl border border-border/70 p-3.5 font-mono text-[11px] space-y-1.5">
              <div className="text-text-subtle font-bold uppercase text-[9px]">Grid Operator Profile</div>
              <div className="text-white font-bold">{selectedGrid.operator} — {selectedGrid.province}</div>
              <div className="text-aurora">{selectedGrid.baseloadSource}</div>
              <div className="text-text-muted">
                Grid Carbon Intensity: <span className="text-emerald-400 font-bold">{selectedGrid.gridCarbonIntensityGramsKWh} gCO2/kWh</span>
              </div>
            </div>
          </div>

          {/* Real-time Telemetry Results */}
          <div className="glass-card space-y-4 rounded-2xl border border-border/80 p-6">
            <div className="flex items-start justify-between gap-4 border-b border-borderSubtle pb-4">
              <div>
                <div className="text-[10px] font-mono uppercase text-text-subtle">Total Sovereign AI Compute Throughput</div>
                <div className="mt-1 text-3xl font-black text-aurora font-mono">
                  {computeModel.totalComputeExaFlops.toFixed(2)} <span className="text-base text-text-muted font-normal">ExaFLOPs (FP8)</span>
                </div>
              </div>
              <span className="rounded-lg bg-emerald-500/20 border border-emerald-500/40 px-2.5 py-1 text-xs font-mono font-bold text-emerald-300">
                100% Sovereign Compute
              </span>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-xl border border-borderSubtle bg-surface/80 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-text-subtle">Installed Accelerator Cluster</div>
                <div className="text-xl font-black text-white font-mono">
                  {computeModel.totalAccelerators.toLocaleString()} <span className="text-xs text-text-muted font-normal">GPUs</span>
                </div>
                <div className="text-[10px] text-text-muted">Net IT Load: {computeModel.itPowerMW.toFixed(1)} MW</div>
              </div>

              <div className="rounded-xl border border-sky-500/30 bg-sky-950/30 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-sky-400">Regional AI Headroom Taken</div>
                <div className="text-xl font-black text-sky-300 font-mono">
                  {computeModel.headroomPct.toFixed(1)}% <span className="text-xs text-text-muted font-normal">of {selectedGrid.availableAiHeadroomMW} MW</span>
                </div>
                <div className="text-[10px] text-sky-400/80">Available on {selectedGrid.operator}</div>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-xl border border-emerald-500/30 bg-emerald-950/30 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-emerald-400">Avoided Annual Emissions</div>
                <div className="text-lg font-black text-emerald-300 font-mono">
                  {Math.round(computeModel.avoidedCarbonTonnes).toLocaleString()} <span className="text-xs font-normal">tCO2e/yr</span>
                </div>
                <div className="text-[10px] text-text-muted">vs fossil thermal grids</div>
              </div>

              <div className="rounded-xl border border-borderSubtle bg-surface/80 p-3.5 space-y-1">
                <div className="text-[10px] font-mono text-text-subtle">Annual Power Consumption</div>
                <div className="text-lg font-black text-white font-mono">
                  {Math.round(computeModel.annualGWh).toLocaleString()} <span className="text-xs font-normal">GWh/yr</span>
                </div>
                <div className="text-[10px] text-text-muted">Carbon: ~{Math.round(computeModel.annualCarbonTonnes)} tCO2e/yr</div>
              </div>
            </div>

            <div className="rounded-xl border border-aurora/40 bg-[#040806] p-4 text-xs font-mono space-y-2">
              <div className="flex items-center justify-between font-bold">
                <span className="text-text-muted">Clean FLOPs / Megawatt Sovereign Standard</span>
                <span className="text-aurora">{selectedGrid.cleanFlopsRatioExaGW} ExaFLOPs / GW</span>
              </div>
              <p className="text-[11px] text-text-subtle leading-relaxed">
                Positions Canada as a premier G7 clean intelligence powerhouse: computing complex foundation models,
                autonomous robotics, and quantum-classical workloads powered entirely by zero-emission baseload.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* SECTION 2: AI Sovereignty Control Checklist */}
      <section className="glass-card rounded-2xl border border-border/80 p-6 shadow-xl sm:p-8 space-y-6">
        <div className="flex flex-col justify-between gap-4 border-b border-borderSubtle pb-5 sm:flex-row sm:items-center">
          <div>
            <div className="inline-flex items-center gap-1.5 font-mono text-[11px] text-gold">
              <Sliders className="h-3.5 w-3.5" /> INTERACTIVE HEURISTIC
            </div>
            <h2 className="mt-1 text-xl font-bold text-text-main">AI Sovereignty Governance Checklist</h2>
            <p className="mt-1 text-xs text-text-muted">Verify operational control, encryption key sovereignty, and Canadian data residency.</p>
          </div>
          <div className="rounded-xl border border-borderSubtle bg-surface px-5 py-3 text-right">
            <div className="font-mono text-[10px] uppercase text-text-subtle">Checklist coverage</div>
            <div className="text-3xl font-black text-aurora">{score}<span className="text-xs font-normal text-text-subtle">/100</span></div>
          </div>
        </div>

        <div className="grid gap-3 md:grid-cols-2">
          {criteria.map((criterion) => {
            const selected = answers[criterion.id];
            return (
              <button
                key={criterion.id}
                type="button"
                aria-pressed={selected}
                onClick={() => setAnswers((current) => ({ ...current, [criterion.id]: !current[criterion.id] }))}
                className={`flex items-start justify-between gap-4 rounded-xl border p-4 text-left transition-colors ${
                  selected ? "border-primary/40 bg-primary/10" : "border-borderSubtle bg-surface hover:border-border"
                }`}
              >
                <div>
                  <div className="flex items-center gap-2 text-sm font-bold text-text-main">
                    {criterion.label}
                    <span className="font-mono text-[10px] font-normal text-aurora">{criterion.weight} pts</span>
                  </div>
                  <p className="mt-1 text-xs leading-relaxed text-text-muted">{criterion.description}</p>
                </div>
                {selected ? <CheckCircle2 className="h-5 w-5 shrink-0 text-aurora" /> : <XCircle className="h-5 w-5 shrink-0 text-text-subtle" />}
              </button>
            );
          })}
        </div>
      </section>

      {/* Disclaimers & Advice */}
      <section className="grid gap-4 md:grid-cols-2">
        <div className="rounded-2xl border border-borderSubtle bg-surface p-5">
          <div className="flex items-center gap-2 font-bold text-text-main">
            <Info className="h-4 w-4 text-aurora" /> How to use the result
          </div>
          <p className="mt-2 text-xs leading-relaxed text-text-muted">
            Treat unchecked items as diligence questions. Attach evidence, name the accountable owner, record the review date, and reassess when architecture or providers change.
          </p>
        </div>
        <div className="rounded-2xl border border-gold/30 bg-gold/5 p-5">
          <div className="flex items-center gap-2 font-bold text-text-main">
            <ShieldAlert className="h-4 w-4 text-gold" /> Statutory Limits
          </div>
          <p className="mt-2 text-xs leading-relaxed text-text-muted">
            The score is a user-entered coverage total, not a risk probability or independent certification. Interconnection approvals and transmission connection queues remain subject to utility tariffs and system operator study processes.
          </p>
        </div>
      </section>
    </div>
  );
}

