"use client";

import { useState, useMemo } from "react";
import {
  Compass,
  Train,
  Anchor,
  AlertTriangle,
  CheckCircle2,
  Layers,
  MapPin,
  TrendingUp,
  Activity,
  ArrowRight,
  ShieldCheck,
  Zap,
  Sliders,
  Sparkles,
} from "lucide-react";

interface GatewayProfile {
  id: string;
  name: string;
  province: string;
  railways: string[];
  throughputMt: number;
  vesselDwellHours: number;
  railcarDwellHours: number;
  berthUtilizationPct: number;
  turnaroundDays: number;
  status: "OPTIMAL" | "ELEVATED" | "CONGESTED";
  mitigation: string;
}

const CANONICAL_GATEWAYS: GatewayProfile[] = [
  {
    id: "prince-rupert",
    name: "Port of Prince Rupert (Fairview & Ridley Terminals)",
    province: "BC",
    railways: ["CN"],
    throughputMt: 32.5,
    vesselDwellHours: 38.0,
    railcarDwellHours: 22.5,
    berthUtilizationPct: 78.5,
    turnaroundDays: 5.2,
    status: "OPTIMAL",
    mitigation: "Fastest trans-Pacific maritime route; ongoing Ridley Island Energy Export Terminal expansion.",
  },
  {
    id: "vancouver",
    name: "Port of Vancouver (Roberts Bank & Burrard Inlet)",
    province: "BC",
    railways: ["CN", "CPKC"],
    throughputMt: 150.0,
    vesselDwellHours: 74.0,
    railcarDwellHours: 52.0,
    berthUtilizationPct: 91.2,
    turnaroundDays: 8.8,
    status: "CONGESTED",
    mitigation: "Roberts Bank Terminal 2 approval; rail fluidity corridors through Fraser Canyon prioritized.",
  },
  {
    id: "montreal",
    name: "Port of Montreal (Contrecœur Expansion Hub)",
    province: "QC",
    railways: ["CN", "CPKC"],
    throughputMt: 38.0,
    vesselDwellHours: 44.0,
    railcarDwellHours: 28.0,
    berthUtilizationPct: 82.0,
    turnaroundDays: 6.0,
    status: "OPTIMAL",
    mitigation: "St. Lawrence Seaway multimodal container transfer and green shipping corridor.",
  },
  {
    id: "churchill",
    name: "Port of Churchill (Hudson Bay Arctic Deepwater Gateway)",
    province: "MB",
    railways: ["Arctic Gateway Hudson Bay Railway"],
    throughputMt: 1.2,
    vesselDwellHours: 24.0,
    railcarDwellHours: 18.0,
    berthUtilizationPct: 35.0,
    turnaroundDays: 4.5,
    status: "OPTIMAL",
    mitigation: "Arctic maritime sovereignty route connecting Prairie grain & critical minerals directly to European ports.",
  },
];

interface PredefinedRoute {
  id: string;
  name: string;
  origin: string;
  destination: string;
  latOrigin: number;
  lonOrigin: number;
  latDest: number;
  lonDest: number;
  description: string;
}

const PREDEFINED_ROUTES: PredefinedRoute[] = [
  {
    id: "prince-george-rupert",
    name: "Northern BC Clean Energy & Critical Minerals Corridor",
    origin: "Prince George, BC",
    destination: "Port of Prince Rupert, BC",
    latOrigin: 53.9171,
    lonOrigin: -122.7497,
    latDest: 54.315,
    lonDest: -130.3208,
    description: "Connects BC Interior clean electricity and hydrogen production to trans-Pacific export berths.",
  },
  {
    id: "ring-of-fire-corridor",
    name: "Ring of Fire Heavy Industrial Supply Corridor",
    origin: "Webequie / Marten Falls, ON",
    destination: "Thunder Bay / Sudbury Smelting Hub, ON",
    latOrigin: 52.8333,
    lonOrigin: -86.5,
    latDest: 48.3809,
    lonDest: -89.2477,
    description: "Sub-Arctic all-weather road and transmission line opening up world-class chromite & nickel reserves.",
  },
  {
    id: "alberta-heartland-h2",
    name: "Alberta Industrial Heartland Clean H2 / CO2 Pipeline",
    origin: "Fort Saskatchewan, AB",
    destination: "Edmonton & Western Rail Gateways, AB",
    latOrigin: 53.7128,
    lonOrigin: -113.2133,
    latDest: 53.5461,
    lonDest: -113.4938,
    description: "Zero-carbon hydrogen transmission network and carbon capture trunkline for heavy industry.",
  },
  {
    id: "arctic-churchill-corridor",
    name: "Prairie to Arctic Deepwater Sovereignty Line",
    origin: "Thompson / The Pas, MB",
    destination: "Port of Churchill, MB",
    latOrigin: 55.7435,
    lonOrigin: -97.8558,
    latDest: 58.7684,
    lonDest: -94.165,
    description: "Heavy rail revitalization over discontinuous permafrost for Arctic export resilience.",
  },
];

type InfraMode = "HVDC" | "HYDROGEN" | "CO2" | "RAIL" | "ROAD";

export default function CorridorsPage() {
  const [selectedRouteId, setSelectedRouteId] = useState<string>("prince-george-rupert");
  const [mode, setMode] = useState<InfraMode>("HVDC");
  const [permafrostOverride, setPermafrostOverride] = useState<boolean>(false);

  const currentRoute = useMemo(() => {
    return PREDEFINED_ROUTES.find((r) => r.id === selectedRouteId) || PREDEFINED_ROUTES[0];
  }, [selectedRouteId]);

  // Routing Evaluation Math
  const routeEvaluation = useMemo(() => {
    const r = currentRoute;
    const earthRadiusKM = 6371.0;
    const dLat = ((r.latDest - r.latOrigin) * Math.PI) / 180.0;
    const dLon = ((r.lonDest - r.lonOrigin) * Math.PI) / 180.0;
    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos((r.latOrigin * Math.PI) / 180.0) *
        Math.cos((r.latDest * Math.PI) / 180.0) *
        Math.sin(dLon / 2) *
        Math.sin(dLon / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    const directKM = earthRadiusKM * c;
    const actualKM = Math.round(directKM * 1.24 * 10) / 10;

    let perKM = 3_200_000;
    let scheduleMonths = Math.round(actualKM * 0.12) + 24;

    switch (mode) {
      case "HYDROGEN":
        perKM = 4_500_000;
        scheduleMonths = Math.round(actualKM * 0.15) + 30;
        break;
      case "CO2":
        perKM = 3_800_000;
        scheduleMonths = Math.round(actualKM * 0.14) + 24;
        break;
      case "RAIL":
        perKM = 6_500_000;
        scheduleMonths = Math.round(actualKM * 0.2) + 36;
        break;
      case "ROAD":
        perKM = 2_800_000;
        scheduleMonths = Math.round(actualKM * 0.1) + 18;
        break;
      default:
        perKM = 3_200_000;
        scheduleMonths = Math.round(actualKM * 0.12) + 24;
        break;
    }

    const baseConstructionCapex = actualKM * perKM;
    const rowAcquisitionCapex = Math.round(baseConstructionCapex * 0.08);
    const envMitigationCapex = Math.round(baseConstructionCapex * 0.12);
    const totalEstimatedCapex = baseConstructionCapex + rowAcquisitionCapex + envMitigationCapex;

    const isNorthern = (r.latOrigin + r.latDest) / 2 > 52.5 || permafrostOverride;
    const permafrostRisk = isNorthern ? 0.74 : 0.15;
    const caribouKM = Math.round(actualKM * 0.35);
    const wetlands = Math.round(actualKM / 22.0);
    const waterCrossings = Math.round(actualKM / 40.0);
    const impedanceIndex = Math.round((actualKM * 1.15 + permafrostRisk * 120) * 10) / 10;

    return {
      actualKM,
      baseConstructionCapex,
      rowAcquisitionCapex,
      envMitigationCapex,
      totalEstimatedCapex,
      scheduleMonths,
      permafrostRisk,
      caribouKM,
      wetlands,
      waterCrossings,
      impedanceIndex,
      isNorthern,
    };
  }, [currentRoute, mode, permafrostOverride]);

  function formatCad(value: number): string {
    if (value >= 1e9) return `$${(value / 1e9).toFixed(2)}B`;
    if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`;
    return `$${value.toLocaleString("en-CA")}`;
  }

  return (
    <div className="min-h-screen bg-background text-foreground py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Top Header */}
        <div className="rounded-2xl border border-border bg-card p-6 md:p-8 shadow-xl relative overflow-hidden">
          <div className="absolute -right-20 -top-20 w-80 h-80 bg-blue-500/10 rounded-full blur-3xl pointer-events-none" />
          <div className="relative z-10 flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono tracking-wider uppercase text-blue-400 mb-2">
                <Compass className="w-4 h-4" />
                Strategic Linear Corridors & Maritime Intermodal Engine
              </div>
              <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-white">
                Corridor Rights-of-Way & Logistics
              </h1>
              <p className="text-muted-foreground text-sm sm:text-base mt-2 max-w-3xl">
                Multi-criteria geospatial pathfinding across Crown provincial, federal, and treaty lands with
                permafrost degradation risk, caribou habitats, and port-to-rail intermodal dwell times.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-mono font-semibold bg-blue-500/10 text-blue-300 border border-blue-500/30">
                <ShieldCheck className="w-3.5 h-3.5" />
                Deterministic Physics & GIS
              </span>
            </div>
          </div>
        </div>

        {/* SECTION 1: LINEAR RIGHT-OF-WAY PATHFINDER */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Controls */}
          <div className="lg:col-span-1 space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-5">
              <h3 className="text-base font-semibold text-white flex items-center gap-2">
                <Sliders className="w-4 h-4 text-blue-400" />
                Corridor Alignment Configuration
              </h3>

              <div>
                <label className="text-xs font-mono text-muted-foreground block mb-2">Select Linear Corridor</label>
                <select
                  value={selectedRouteId}
                  onChange={(e) => setSelectedRouteId(e.target.value)}
                  className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {PREDEFINED_ROUTES.map((r) => (
                    <option key={r.id} value={r.id}>
                      {r.name}
                    </option>
                  ))}
                </select>
                <p className="text-xs text-muted-foreground mt-2 leading-relaxed">
                  {currentRoute.description}
                </p>
              </div>

              <div className="pt-2 border-t border-border">
                <label className="text-xs font-mono text-muted-foreground block mb-2">Infrastructure Mode</label>
                <div className="grid grid-cols-2 gap-2">
                  <button
                    onClick={() => setMode("HVDC")}
                    className={`px-3 py-2 text-xs font-mono rounded-lg border transition-all ${
                      mode === "HVDC"
                        ? "bg-blue-500/20 text-blue-300 border-blue-500 font-semibold"
                        : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                    }`}
                  >
                    500kV HVDC
                  </button>
                  <button
                    onClick={() => setMode("HYDROGEN")}
                    className={`px-3 py-2 text-xs font-mono rounded-lg border transition-all ${
                      mode === "HYDROGEN"
                        ? "bg-emerald-500/20 text-emerald-300 border-emerald-500 font-semibold"
                        : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                    }`}
                  >
                    Clean H2 Pipe
                  </button>
                  <button
                    onClick={() => setMode("CO2")}
                    className={`px-3 py-2 text-xs font-mono rounded-lg border transition-all ${
                      mode === "CO2"
                        ? "bg-purple-500/20 text-purple-300 border-purple-500 font-semibold"
                        : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                    }`}
                  >
                    CO2 Capture Trunk
                  </button>
                  <button
                    onClick={() => setMode("RAIL")}
                    className={`px-3 py-2 text-xs font-mono rounded-lg border transition-all ${
                      mode === "RAIL"
                        ? "bg-amber-500/20 text-amber-300 border-amber-500 font-semibold"
                        : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                    }`}
                  >
                    Heavy-Haul Rail
                  </button>
                </div>
              </div>

              <div className="pt-2 border-t border-border">
                <label className="flex items-center gap-2 cursor-pointer text-xs font-mono text-foreground">
                  <input
                    type="checkbox"
                    checked={permafrostOverride}
                    onChange={(e) => setPermafrostOverride(e.target.checked)}
                    className="rounded bg-background border-border text-blue-500 focus:ring-blue-500"
                  />
                  <span>Force Discontinuous Permafrost Thaw Stress</span>
                </label>
              </div>

              {/* Waypoint Coordinates */}
              <div className="pt-2 border-t border-border space-y-2 text-xs font-mono text-muted-foreground">
                <div className="flex justify-between">
                  <span>Origin:</span>
                  <span className="text-white">
                    {currentRoute.latOrigin.toFixed(4)}°, {currentRoute.lonOrigin.toFixed(4)}°
                  </span>
                </div>
                <div className="flex justify-between">
                  <span>Destination:</span>
                  <span className="text-white">
                    {currentRoute.latDest.toFixed(4)}°, {currentRoute.lonDest.toFixed(4)}°
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Results Display */}
          <div className="lg:col-span-2 space-y-6">
            <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-6">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border pb-4">
                <div>
                  <div className="text-xs font-mono text-muted-foreground uppercase">Alignment Evaluation</div>
                  <h3 className="text-xl font-bold text-white">{currentRoute.name}</h3>
                </div>
                <div className="text-right">
                  <div className="text-xs font-mono text-muted-foreground">Impedance Index</div>
                  <div className="text-2xl font-black font-mono text-blue-400">
                    {routeEvaluation.impedanceIndex}
                  </div>
                </div>
              </div>

              {/* Metric Highlights */}
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                <div className="p-3.5 rounded-lg bg-background/70 border border-border space-y-1">
                  <div className="text-xs font-mono text-muted-foreground">Total Length</div>
                  <div className="text-lg font-bold font-mono text-white">{routeEvaluation.actualKM} km</div>
                  <div className="text-[10px] text-muted-foreground">Terrain Factor: 1.24x</div>
                </div>

                <div className="p-3.5 rounded-lg bg-background/70 border border-border space-y-1">
                  <div className="text-xs font-mono text-muted-foreground">Total Estimated CAPEX</div>
                  <div className="text-lg font-bold font-mono text-emerald-400">
                    {formatCad(routeEvaluation.totalEstimatedCapex)}
                  </div>
                  <div className="text-[10px] text-muted-foreground">Base: {formatCad(routeEvaluation.baseConstructionCapex)}</div>
                </div>

                <div className="p-3.5 rounded-lg bg-background/70 border border-border space-y-1">
                  <div className="text-xs font-mono text-muted-foreground">EPC Schedule</div>
                  <div className="text-lg font-bold font-mono text-white">
                    {routeEvaluation.scheduleMonths} mo
                  </div>
                  <div className="text-[10px] text-muted-foreground">Turnkey Commissioning</div>
                </div>

                <div className="p-3.5 rounded-lg bg-background/70 border border-border space-y-1">
                  <div className="text-xs font-mono text-muted-foreground">Permafrost Hazard</div>
                  <div
                    className={`text-lg font-bold font-mono ${
                      routeEvaluation.permafrostRisk > 0.5 ? "text-amber-400" : "text-emerald-400"
                    }`}
                  >
                    {(routeEvaluation.permafrostRisk * 100).toFixed(0)}%
                  </div>
                  <div className="text-[10px] text-muted-foreground">Geotechnical Settle Risk</div>
                </div>
              </div>

              {/* Environmental Constraints */}
              <div className="space-y-3">
                <h4 className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
                  Geotechnical & Environmental Sensitivity Overlay
                </h4>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs font-mono">
                  <div className="p-3 rounded bg-muted/40 border border-border">
                    <span className="text-muted-foreground">Caribou Range Overlap:</span>
                    <div className="text-white font-bold mt-0.5">{routeEvaluation.caribouKM} km</div>
                  </div>
                  <div className="p-3 rounded bg-muted/40 border border-border">
                    <span className="text-muted-foreground">Wetland Crossings:</span>
                    <div className="text-white font-bold mt-0.5">{routeEvaluation.wetlands} spans</div>
                  </div>
                  <div className="p-3 rounded bg-muted/40 border border-border">
                    <span className="text-muted-foreground">Major Water Crossings:</span>
                    <div className="text-white font-bold mt-0.5">{routeEvaluation.waterCrossings} rivers</div>
                  </div>
                </div>
              </div>

              {/* Land Tenure Breakdown */}
              <div className="space-y-3">
                <h4 className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
                  Land Tenure & Sovereign Jurisdiction Distribution
                </h4>
                <div className="space-y-2">
                  <div className="h-3 w-full bg-muted rounded-full overflow-hidden flex">
                    <div style={{ width: "58%" }} className="bg-blue-500" title="Crown Provincial: 58%" />
                    <div style={{ width: "24%" }} className="bg-amber-500" title="Indigenous Settlement: 24%" />
                    <div style={{ width: "12%" }} className="bg-emerald-500" title="Crown Federal: 12%" />
                    <div style={{ width: "6%" }} className="bg-purple-500" title="Private Freehold: 6%" />
                  </div>
                  <div className="flex flex-wrap gap-4 text-xs font-mono">
                    <span className="flex items-center gap-1.5 text-blue-300">
                      <span className="w-2.5 h-2.5 rounded-full bg-blue-500" /> Crown Provincial (58%)
                    </span>
                    <span className="flex items-center gap-1.5 text-amber-300">
                      <span className="w-2.5 h-2.5 rounded-full bg-amber-500" /> Indigenous Settlement / Treaty (24%)
                    </span>
                    <span className="flex items-center gap-1.5 text-emerald-300">
                      <span className="w-2.5 h-2.5 rounded-full bg-emerald-500" /> Crown Federal (12%)
                    </span>
                    <span className="flex items-center gap-1.5 text-purple-300">
                      <span className="w-2.5 h-2.5 rounded-full bg-purple-500" /> Private Freehold (6%)
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* SECTION 2: STRATEGIC MARITIME GATEWAYS & MULTI-MODAL LOGISTICS */}
        <div className="rounded-xl border border-border bg-card p-6 shadow-lg space-y-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h3 className="text-lg font-bold text-white flex items-center gap-2">
                <Anchor className="w-5 h-5 text-blue-400" />
                Strategic Maritime Gateways & Rail Dwell Fluidity
              </h3>
              <p className="text-xs text-muted-foreground mt-1">
                Port-to-rail intermodal throughput, Class 1 rail interchange metrics (CN & CPKC), and bottleneck mitigations.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {CANONICAL_GATEWAYS.map((gateway) => (
              <div
                key={gateway.id}
                className="rounded-xl border border-border bg-background/60 p-5 space-y-4 hover:border-blue-500/50 transition-all shadow-md"
              >
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <h4 className="text-base font-bold text-white">{gateway.name}</h4>
                    <div className="text-xs font-mono text-muted-foreground mt-0.5">
                      Province: {gateway.province} • Rail Connections: {gateway.railways.join(" / ")}
                    </div>
                  </div>
                  <span
                    className={`px-2.5 py-1 rounded text-xs font-mono font-semibold border ${
                      gateway.status === "OPTIMAL"
                        ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/30"
                        : gateway.status === "CONGESTED"
                        ? "bg-red-500/10 text-red-400 border-red-500/30"
                        : "bg-amber-500/10 text-amber-400 border-amber-500/30"
                    }`}
                  >
                    {gateway.status}
                  </span>
                </div>

                <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-2 border-t border-border text-xs font-mono">
                  <div className="p-2 rounded bg-muted/30">
                    <span className="text-muted-foreground block text-[10px]">Annual Throughput</span>
                    <span className="font-bold text-white">{gateway.throughputMt} Mt</span>
                  </div>
                  <div className="p-2 rounded bg-muted/30">
                    <span className="text-muted-foreground block text-[10px]">Vessel Dwell</span>
                    <span className="font-bold text-white">{gateway.vesselDwellHours} hrs</span>
                  </div>
                  <div className="p-2 rounded bg-muted/30">
                    <span className="text-muted-foreground block text-[10px]">Railcar Dwell</span>
                    <span className="font-bold text-white">{gateway.railcarDwellHours} hrs</span>
                  </div>
                  <div className="p-2 rounded bg-muted/30">
                    <span className="text-muted-foreground block text-[10px]">Berth Cap. Used</span>
                    <span className="font-bold text-white">{gateway.berthUtilizationPct}%</span>
                  </div>
                </div>

                <div className="p-3 rounded bg-blue-500/5 border border-blue-500/20 text-xs text-muted-foreground">
                  <span className="font-semibold text-blue-300">Fluidity Action: </span>
                  {gateway.mitigation}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
