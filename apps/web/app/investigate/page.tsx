"use client";

import { useState, useMemo, useRef, useEffect } from "react";
import Link from "next/link";
import {
  Network,
  Search,
  ShieldAlert,
  ShieldCheck,
  AlertTriangle,
  Compass,
  Cpu,
  Globe2,
  Layers,
  Maximize2,
  Minimize2,
  ZoomIn,
  ZoomOut,
  RefreshCw,
  Sliders,
  ChevronRight,
  Filter,
  Eye,
  Crosshair,
  Lock,
  Building2,
  Coins,
  Anchor,
  Radio,
  FileText,
  Activity,
  Zap,
} from "lucide-react";
import CapitalStackBuilder from "@/components/finance/CapitalStackBuilder";

// Entity Node Type Definitions
export interface GraphNode {
  id: string;
  label: string;
  type: "PROJECT" | "SOE_ENTITY" | "INDIGENOUS_NATION" | "INFRASTRUCTURE" | "FINANCIAL_VEHICLE" | "OFFSHORE_SHELL";
  category: string;
  jurisdiction: string;
  classification: "UNCLASSIFIED" | "PROTECTED_B" | "SECRET" | "TOP_SECRET";
  riskLevel: "CRITICAL" | "ELEVATED" | "NOMINAL" | "CLEARED";
  x: number;
  y: number;
  vx?: number;
  vy?: number;
  capexCad?: string;
  uboUltimateOwner?: string;
  sasacAffiliation?: string;
  c69ClockDaysRemaining?: number;
  insarSubsidenceMm?: number;
  description: string;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  label: string;
  type: "OWNS" | "FINANCES" | "ICA_NATIONAL_SECURITY_REVIEW" | "CORRIDOR_LINK" | "CONSULTATION_TERRITORY" | "OFFSHORE_CONTROL";
  flagged: boolean;
}

// Canonical Investigation Knowledge Graph (Gotham Intelligence Mesh)
const INITIAL_NODES: GraphNode[] = [
  {
    id: "crawford_nickel",
    label: "Crawford Nickel Sulphide Project",
    type: "PROJECT",
    category: "Critical Minerals",
    jurisdiction: "ON (Timmins)",
    classification: "PROTECTED_B",
    riskLevel: "ELEVATED",
    x: 480,
    y: 280,
    capexCad: "$3.5B CAD",
    uboUltimateOwner: "Canada Nickel Corp (Public / Dispersed)",
    sasacAffiliation: "None Direct (Minority Offtake Under ICA Section 25.3 Review)",
    c69ClockDaysRemaining: 184,
    insarSubsidenceMm: -1.2,
    description: "World's largest nickel sulphide discovery in 20+ years. High strategic value for North American battery cell cathode supply chain.",
  },
  {
    id: "noront_eagles_nest",
    label: "Eagle's Nest Ni-Cu-PGE (Wyloo)",
    type: "PROJECT",
    category: "Critical Minerals",
    jurisdiction: "ON (Ring of Fire)",
    classification: "PROTECTED_B",
    riskLevel: "NOMINAL",
    x: 320,
    y: 200,
    capexCad: "$2.2B CAD",
    uboUltimateOwner: "Wyloo Metals (Tattarang / Andrew Forrest, Australia)",
    sasacAffiliation: "None (Five Eyes Allied Capital)",
    c69ClockDaysRemaining: 310,
    insarSubsidenceMm: -0.4,
    description: "High-grade underground mine anchoring James Bay Lowlands road and grid infrastructure.",
  },
  {
    id: "martik_resources",
    label: "Northern Horizon Minerals Ltd",
    type: "OFFSHORE_SHELL",
    category: "Offshore Shell",
    jurisdiction: "Cayman Islands / BVI",
    classification: "SECRET",
    riskLevel: "CRITICAL",
    x: 680,
    y: 180,
    capexCad: "$450M CAD",
    uboUltimateOwner: "Chinalco / Sinomine Group (Beijing SASAC Tier 2)",
    sasacAffiliation: "Direct SASAC Oversight (Order in Council Divestment Candidate)",
    description: "Multi-layered corporate shell holding minority shares and preemptive offtake agreements in Ontario mineral claims.",
  },
  {
    id: "cib_anchor",
    label: "Canada Infrastructure Bank (CIB)",
    type: "FINANCIAL_VEHICLE",
    category: "Crown Corporation",
    jurisdiction: "Federal (Canada)",
    classification: "UNCLASSIFIED",
    riskLevel: "CLEARED",
    x: 360,
    y: 420,
    capexCad: "$35.0B Portfolio",
    description: "Federal crown bank deploying low-cost concessionary debt and clean energy transmission loans.",
  },
  {
    id: "matawa_first_nations",
    label: "Matawa First Nations Council",
    type: "INDIGENOUS_NATION",
    category: "Treaty / Indigenous Nation",
    jurisdiction: "ON (Treaty 9)",
    classification: "UNCLASSIFIED",
    riskLevel: "CLEARED",
    x: 180,
    y: 270,
    capexCad: "ILGP Backed Tranche",
    description: "Tribal council representing 9 First Nations in Northern Ontario holding ancestral territorial stewardship over Ring of Fire road corridors.",
  },
  {
    id: "port_prince_rupert",
    label: "Port of Prince Rupert (Fairview Terminal)",
    type: "INFRASTRUCTURE",
    category: "Maritime Gateway",
    jurisdiction: "BC (North Coast)",
    classification: "PROTECTED_B",
    riskLevel: "NOMINAL",
    x: 140,
    y: 460,
    capexCad: "$1.8B Terminal Expansion",
    description: "Deepest natural harbor in North America, closest gateway to Asian trade corridors with CN Rail transcontinental connectivity.",
  },
  {
    id: "sasac_foreign_state",
    label: "State-owned Assets Supervision & Administration Commission (SASAC)",
    type: "SOE_ENTITY",
    category: "Foreign State Organ",
    jurisdiction: "Beijing, PRC",
    classification: "TOP_SECRET",
    riskLevel: "CRITICAL",
    x: 880,
    y: 130,
    capexCad: "Sovereign Holding",
    uboUltimateOwner: "State Council of the People's Republic of China",
    sasacAffiliation: "Direct Ministerial Organ",
    description: "Foreign state ministry coordinating strategic mineral stockpiling and critical infrastructure acquisitions across overseas jurisdictions.",
  },
  {
    id: "transmountain_expansion",
    label: "Trans Mountain Pipeline System (TMX)",
    type: "INFRASTRUCTURE",
    category: "Energy Corridor",
    jurisdiction: "AB / BC",
    classification: "PROTECTED_B",
    riskLevel: "NOMINAL",
    x: 260,
    y: 560,
    capexCad: "$34.0B CAD",
    uboUltimateOwner: "Trans Mountain Corp (Federal Crown / Indigenous Ownership Transition)",
    description: "890,000 bpd heavy crude and refined petroleum export system to Westridge Marine Terminal in Burnaby.",
  },
  {
    id: "churchill_arctic_corridor",
    label: "Port of Churchill & Hudson Bay Railway",
    type: "INFRASTRUCTURE",
    category: "Arctic Deepwater",
    jurisdiction: "MB (Hudson Bay)",
    classification: "SECRET",
    riskLevel: "NOMINAL",
    x: 520,
    y: 470,
    capexCad: "$450M Sovereign Upgrade",
    uboUltimateOwner: "Arctic Gateway Group (Consortium of 41 First Nation and Bayline Communities)",
    description: "Canada's only Arctic deep-water seaport linked to North American rail grid, critical for European grain and mineral exports.",
  },
];

const INITIAL_EDGES: GraphEdge[] = [
  { id: "e1", source: "sasac_foreign_state", target: "martik_resources", label: "BENEFICIAL_OWNERSHIP (72%)", type: "OFFSHORE_CONTROL", flagged: true },
  { id: "e2", source: "martik_resources", target: "crawford_nickel", label: "ICA S.25.3 NATIONAL SECURITY SCRUTINY", type: "ICA_NATIONAL_SECURITY_REVIEW", flagged: true },
  { id: "e3", source: "cib_anchor", target: "noront_eagles_nest", label: "CONCESSIONARY_FINANCE ($400M CAD)", type: "FINANCES", flagged: false },
  { id: "e4", source: "matawa_first_nations", target: "noront_eagles_nest", label: "FREE PRIOR & INFORMED CONSENT / TRC 92", type: "CONSULTATION_TERRITORY", flagged: false },
  { id: "e5", source: "noront_eagles_nest", target: "port_prince_rupert", label: "CN RAIL CORRIDOR TRANSIT", type: "CORRIDOR_LINK", flagged: false },
  { id: "e6", source: "cib_anchor", target: "churchill_arctic_corridor", label: "ARCTIC SOVEREIGNTY FACILITY ($150M)", type: "FINANCES", flagged: false },
  { id: "e7", source: "transmountain_expansion", target: "port_prince_rupert", label: "MARITIME LOGISTICS INTERCONNECT", type: "CORRIDOR_LINK", flagged: false },
  { id: "e8", source: "cib_anchor", target: "crawford_nickel", label: "CLEAN TECH ITC CO-FINANCE", type: "FINANCES", flagged: false },
];

export default function InvestigatePage() {
  const [nodes, setNodes] = useState<GraphNode[]>(INITIAL_NODES);
  const [edges] = useState<GraphEdge[]>(INITIAL_EDGES);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedNodeId, setSelectedNodeId] = useState<string>("martik_resources");
  const [filterRisk, setFilterRisk] = useState<string>("ALL");
  const [zoomLevel, setZoomLevel] = useState<number>(1);
  const [panOffset, setPanOffset] = useState({ x: 40, y: 30 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const [highlightedPath, setHighlightedPath] = useState<string[]>(["sasac_foreign_state", "martik_resources", "crawford_nickel"]);
  const [activeTab, setActiveTab] = useState<"DOSSIER" | "CAPITAL_STACK" | "SECURITY_TRACE">("DOSSIER");

  const selectedNode = useMemo(() => {
    return nodes.find((n) => n.id === selectedNodeId) || nodes[0];
  }, [nodes, selectedNodeId]);

  const filteredNodes = useMemo(() => {
    return nodes.filter((n) => {
      const matchesSearch =
        n.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
        n.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
        n.jurisdiction.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesRisk = filterRisk === "ALL" || n.riskLevel === filterRisk;
      return matchesSearch && matchesRisk;
    });
  }, [nodes, searchQuery, filterRisk]);

  const filteredNodeIds = useMemo(() => new Set(filteredNodes.map((n) => n.id)), [filteredNodes]);

  // Mouse pan handlers for SVG viewport
  const handleMouseDown = (e: React.MouseEvent<SVGSVGElement>) => {
    if (e.button !== 0) return;
    setIsDragging(true);
    setDragStart({ x: e.clientX - panOffset.x, y: e.clientY - panOffset.y });
  };

  const handleMouseMove = (e: React.MouseEvent<SVGSVGElement>) => {
    if (!isDragging) return;
    setPanOffset({ x: e.clientX - dragStart.x, y: e.clientY - dragStart.y });
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const getNodeColor = (node: GraphNode) => {
    switch (node.riskLevel) {
      case "CRITICAL":
        return "#EF4444"; // red
      case "ELEVATED":
        return "#F59E0B"; // amber
      case "CLEARED":
        return "#10B981"; // emerald
      default:
        return "#3B82F6"; // blue
    }
  };

  return (
    <div className="min-h-screen bg-[#060D0A] text-text-main flex flex-col font-sans selection:bg-primary/30">
      {/* Sovereign Intelligence Banner (Classification & Caveat) */}
      <div className="border-b border-red-500/30 bg-red-950/40 px-4 py-1.5 text-center text-xs font-mono font-bold tracking-widest text-red-400 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <ShieldAlert className="h-4 w-4 text-red-500 animate-pulse" />
          <span>TOP SECRET // CANADIAN EYES ONLY // DISSEMINATION CONTROLLED</span>
        </div>
        <div className="hidden sm:flex items-center gap-4 text-[11px] text-red-300">
          <span>ICA SECTION 25.3 NATIONAL SECURITY SURVEILLANCE MESH</span>
          <span>AUDIT HASH: #SOV-CA-2026-991A</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-red-500 animate-ping" />
          <span className="text-[10px]">LIVE WORKBENCH</span>
        </div>
      </div>

      {/* Top Controls Toolbar */}
      <div className="border-b border-border bg-surface/90 backdrop-blur-md px-6 py-3 flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg border border-primary/40 bg-primary/10 shadow-inner">
            <Network className="h-5 w-5 text-aurora" />
          </div>
          <div>
            <h1 className="text-base font-black tracking-tight text-white flex items-center gap-2">
              Investigation Workbench & Knowledge Mesh
              <span className="rounded border border-primary/40 bg-primary/20 px-1.5 py-0.2 font-mono text-[9px] text-aurora">
                GOTHAM 4.0 PARITY
              </span>
            </h1>
            <p className="text-[11px] text-text-muted">
              Interactive entity-link graph visualizer for beneficial ownership unraveling, SASAC foreign interference detection, and corridor choke-point analysis.
            </p>
          </div>
        </div>

        {/* Filter Controls & Search */}
        <div className="flex items-center gap-3">
          {/* Search Box */}
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-text-subtle" />
            <input
              type="text"
              placeholder="Search entities, owners, jurisdictions..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="h-8 w-64 rounded-lg border border-borderSubtle bg-card pl-8 pr-3 text-xs text-text-main placeholder-text-subtle focus:border-primary focus:outline-none"
            />
          </div>

          {/* Risk Level Filter */}
          <div className="flex items-center gap-1 bg-card p-1 rounded-lg border border-borderSubtle">
            <Filter className="h-3 w-3 text-text-subtle ml-1.5" />
            {(["ALL", "CRITICAL", "ELEVATED", "CLEARED"] as const).map((r) => (
              <button
                key={r}
                onClick={() => setFilterRisk(r)}
                className={`px-2 py-0.5 rounded text-[10px] font-mono font-bold transition-colors ${
                  filterRisk === r ? "bg-primary text-black" : "text-text-muted hover:text-text-main"
                }`}
              >
                {r}
              </button>
            ))}
          </div>

          {/* Zoom / Reset Buttons */}
          <div className="flex items-center gap-1 bg-card p-1 rounded-lg border border-borderSubtle">
            <button
              onClick={() => setZoomLevel((z) => Math.min(2.5, z + 0.15))}
              className="p-1 rounded hover:bg-surface text-text-muted hover:text-text-main"
              title="Zoom In"
            >
              <ZoomIn className="h-3.5 w-3.5" />
            </button>
            <button
              onClick={() => setZoomLevel((z) => Math.max(0.4, z - 0.15))}
              className="p-1 rounded hover:bg-surface text-text-muted hover:text-text-main"
              title="Zoom Out"
            >
              <ZoomOut className="h-3.5 w-3.5" />
            </button>
            <button
              onClick={() => {
                setZoomLevel(1);
                setPanOffset({ x: 40, y: 30 });
              }}
              className="p-1 rounded hover:bg-surface text-text-muted hover:text-text-main"
              title="Reset View"
            >
              <RefreshCw className="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>

      {/* Main Investigation Split Canvas */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left: Interactive SVG Graph Canvas */}
        <div className="flex-1 relative bg-[#040907] cursor-grab active:cursor-grabbing overflow-hidden">
          {/* Grid background styling */}
          <div
            className="absolute inset-0 opacity-15 pointer-events-none"
            style={{
              backgroundImage: `radial-gradient(circle at 1px 1px, #00F5A0 1px, transparent 0)`,
              backgroundSize: "28px 28px",
            }}
          />

          <svg
            className="w-full h-full"
            onMouseDown={handleMouseDown}
            onMouseMove={handleMouseMove}
            onMouseUp={handleMouseUp}
          >
            <g transform={`translate(${panOffset.x}, ${panOffset.y}) scale(${zoomLevel})`}>
              <defs>
                <marker
                  id="arrow"
                  viewBox="0 0 10 10"
                  refX="22"
                  refY="5"
                  markerWidth="6"
                  markerHeight="6"
                  orient="auto-start-reverse"
                >
                  <path d="M 0 0 L 10 5 L 0 10 z" fill="#00F5A0" opacity="0.6" />
                </marker>
                <marker
                  id="arrow-flagged"
                  viewBox="0 0 10 10"
                  refX="22"
                  refY="5"
                  markerWidth="6"
                  markerHeight="6"
                  orient="auto-start-reverse"
                >
                  <path d="M 0 0 L 10 5 L 0 10 z" fill="#EF4444" opacity="0.9" />
                </marker>
              </defs>

              {/* Render Edges */}
              {edges.map((edge) => {
                const src = nodes.find((n) => n.id === edge.source);
                const tgt = nodes.find((n) => n.id === edge.target);
                if (!src || !tgt) return null;

                const isHighlighted =
                  highlightedPath.includes(edge.source) && highlightedPath.includes(edge.target);

                return (
                  <g key={edge.id} className="transition-opacity">
                    <line
                      x1={src.x}
                      y1={src.y}
                      x2={tgt.x}
                      y2={tgt.y}
                      stroke={edge.flagged ? "#EF4444" : isHighlighted ? "#00F5A0" : "#244238"}
                      strokeWidth={edge.flagged || isHighlighted ? 2.5 : 1.2}
                      strokeDasharray={edge.flagged ? "6,4" : undefined}
                      markerEnd={edge.flagged ? "url(#arrow-flagged)" : "url(#arrow)"}
                    />
                    {/* Edge Label text at midpoint */}
                    <text
                      x={(src.x + tgt.x) / 2}
                      y={(src.y + tgt.y) / 2 - 6}
                      fill={edge.flagged ? "#FCA5A5" : "#6EE7B7"}
                      fontSize={9}
                      fontFamily="monospace"
                      textAnchor="middle"
                      className="select-none pointer-events-none drop-shadow"
                    >
                      {edge.label}
                    </text>
                  </g>
                );
              })}

              {/* Render Nodes */}
              {nodes.map((node) => {
                const isSelected = node.id === selectedNodeId;
                const isMatch = filteredNodeIds.has(node.id);
                const color = getNodeColor(node);

                return (
                  <g
                    key={node.id}
                    transform={`translate(${node.x}, ${node.y})`}
                    onClick={() => setSelectedNodeId(node.id)}
                    className="cursor-pointer group"
                    opacity={isMatch ? 1 : 0.25}
                  >
                    {/* Selection Pulse Ring */}
                    {isSelected && (
                      <circle
                        r={26}
                        fill="none"
                        stroke={color}
                        strokeWidth={2}
                        strokeDasharray="4,4"
                        className="animate-spin"
                        style={{ animationDuration: "12s" }}
                      />
                    )}

                    {/* Node Core Body */}
                    <circle
                      r={18}
                      fill="#0C1A14"
                      stroke={color}
                      strokeWidth={isSelected ? 3 : 1.8}
                      className="transition-all group-hover:scale-110"
                    />

                    {/* Inner Type Glyph */}
                    <circle r={6} fill={color} />

                    {/* Node Text Label */}
                    <text
                      y={30}
                      textAnchor="middle"
                      fill={isSelected ? "#00F5A0" : "#E2E8F0"}
                      fontSize={11}
                      fontWeight="bold"
                      fontFamily="sans-serif"
                      className="select-none pointer-events-none drop-shadow-md"
                    >
                      {node.label.length > 24 ? node.label.slice(0, 22) + "…" : node.label}
                    </text>

                    {/* Sub-label badge */}
                    <text
                      y={42}
                      textAnchor="middle"
                      fill="#94A3B8"
                      fontSize={8}
                      fontFamily="monospace"
                      className="select-none pointer-events-none"
                    >
                      {node.jurisdiction} • {node.category}
                    </text>
                  </g>
                );
              })}
            </g>
          </svg>

          {/* Canvas Overlay Legend */}
          <div className="absolute bottom-4 left-4 rounded-xl border border-borderSubtle bg-card/90 backdrop-blur-md p-3 text-xs shadow-xl pointer-events-auto">
            <div className="font-mono text-[10px] font-bold text-text-muted mb-2 tracking-wider">
              INVESTIGATION LEGEND
            </div>
            <div className="space-y-1.5 text-[11px]">
              <div className="flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-red-500" />
                <span className="text-text-main">Critical SASAC / SOE Scrutiny</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-amber-500" />
                <span className="text-text-main">Elevated ICA / Regulatory Scrutiny</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-emerald-500" />
                <span className="text-text-main">Cleared / Sovereign Allied</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-blue-500" />
                <span className="text-text-main">Domestic Anchor Project</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right: Gotham Entity Dossier & Investigation Inspector Drawer */}
        <div className="w-[480px] border-l border-border bg-[#08130E] flex flex-col overflow-y-auto">
          {/* Dossier Tabs */}
          <div className="border-b border-borderSubtle bg-surface/60 flex items-center">
            <button
              onClick={() => setActiveTab("DOSSIER")}
              className={`flex-1 py-3 text-xs font-bold border-b-2 transition-colors ${
                activeTab === "DOSSIER"
                  ? "border-primary text-aurora bg-card/50"
                  : "border-transparent text-text-muted hover:text-text-main"
              }`}
            >
              Entity Dossier
            </button>
            <button
              onClick={() => setActiveTab("CAPITAL_STACK")}
              className={`flex-1 py-3 text-xs font-bold border-b-2 transition-colors ${
                activeTab === "CAPITAL_STACK"
                  ? "border-primary text-aurora bg-card/50"
                  : "border-transparent text-text-muted hover:text-text-main"
              }`}
            >
              Capital Stack
            </button>
            <button
              onClick={() => setActiveTab("SECURITY_TRACE")}
              className={`flex-1 py-3 text-xs font-bold border-b-2 transition-colors ${
                activeTab === "SECURITY_TRACE"
                  ? "border-primary text-aurora bg-card/50"
                  : "border-transparent text-text-muted hover:text-text-main"
              }`}
            >
              ICA / UBO Trace
            </button>
          </div>

          <div className="p-6 space-y-6 flex-1">
            {activeTab === "DOSSIER" && (
              <>
                {/* Entity Header */}
                <div>
                  <div className="flex items-center justify-between">
                    <span
                      className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 font-mono text-[10px] font-bold ${
                        selectedNode.riskLevel === "CRITICAL"
                          ? "bg-red-500/20 text-red-400 border border-red-500/40"
                          : selectedNode.riskLevel === "ELEVATED"
                          ? "bg-amber-500/20 text-amber-400 border border-amber-500/40"
                          : "bg-emerald-500/20 text-emerald-400 border border-emerald-500/40"
                      }`}
                    >
                      {selectedNode.riskLevel} RISK
                    </span>
                    <span className="font-mono text-xs text-text-subtle">
                      SEC-ID: #{selectedNode.id.toUpperCase()}
                    </span>
                  </div>

                  <h2 className="mt-2 text-xl font-black tracking-tight text-white">
                    {selectedNode.label}
                  </h2>
                  <p className="mt-1 text-xs text-text-muted">{selectedNode.description}</p>
                </div>

                {/* Key Telemetry Matrix */}
                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-xl border border-borderSubtle bg-surface p-3">
                    <div className="text-[10px] font-mono text-text-subtle">JURISDICTION</div>
                    <div className="mt-1 text-xs font-bold text-text-main">{selectedNode.jurisdiction}</div>
                  </div>
                  <div className="rounded-xl border border-borderSubtle bg-surface p-3">
                    <div className="text-[10px] font-mono text-text-subtle">CLEARANCE LEVEL</div>
                    <div className="mt-1 text-xs font-bold text-aurora">{selectedNode.classification}</div>
                  </div>
                  <div className="rounded-xl border border-borderSubtle bg-surface p-3">
                    <div className="text-[10px] font-mono text-text-subtle">ESTIMATED VALUATION</div>
                    <div className="mt-1 text-xs font-bold text-text-main">{selectedNode.capexCad || "N/A"}</div>
                  </div>
                  <div className="rounded-xl border border-borderSubtle bg-surface p-3">
                    <div className="text-[10px] font-mono text-text-subtle">BILL C-69 CLOCK</div>
                    <div className="mt-1 text-xs font-bold text-text-main">
                      {selectedNode.c69ClockDaysRemaining ? `${selectedNode.c69ClockDaysRemaining} Days Rem.` : "EXEMPT / OPERATING"}
                    </div>
                  </div>
                </div>

                {/* Satellite Radar & InSAR Telemetry */}
                {selectedNode.insarSubsidenceMm !== undefined && (
                  <div className="rounded-xl border border-borderSubtle bg-surface p-4 space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="flex items-center gap-1.5 text-xs font-semibold text-text-muted">
                        <Radio className="h-4 w-4 text-aurora animate-pulse" />
                        Sentinel-1 InSAR Millimeter Stability
                      </span>
                      <span className="font-mono text-xs font-bold text-aurora">
                        {selectedNode.insarSubsidenceMm > 0 ? `+${selectedNode.insarSubsidenceMm}` : selectedNode.insarSubsidenceMm} mm/yr
                      </span>
                    </div>
                    <div className="w-full bg-card h-2 rounded-full overflow-hidden">
                      <div className="bg-emerald-500 h-full w-[85%]" />
                    </div>
                    <p className="text-[10px] text-text-subtle">
                      Interferometric SAR ground deformation analysis confirming zero anomalous crest slippage.
                    </p>
                  </div>
                )}

                {/* Ultimate Beneficial Ownership Card */}
                {selectedNode.uboUltimateOwner && (
                  <div className="rounded-xl border border-borderSubtle bg-surface p-4 space-y-2">
                    <div className="flex items-center gap-2 text-xs font-semibold text-text-main">
                      <Building2 className="h-4 w-4 text-amber-400" />
                      Ultimate Beneficial Ownership (UBO)
                    </div>
                    <div className="font-mono text-xs font-bold text-amber-300">
                      {selectedNode.uboUltimateOwner}
                    </div>
                    {selectedNode.sasacAffiliation && (
                      <div className="mt-2 rounded-lg bg-red-950/50 border border-red-500/30 p-2 text-[11px] text-red-300">
                        <span className="font-bold">Foreign SOE Risk: </span>
                        {selectedNode.sasacAffiliation}
                      </div>
                    )}
                  </div>
                )}
              </>
            )}

            {activeTab === "CAPITAL_STACK" && (
              <div className="space-y-4">
                <CapitalStackBuilder />
              </div>
            )}

            {activeTab === "SECURITY_TRACE" && (
              <div className="space-y-4">
                <div className="rounded-xl border border-red-500/40 bg-red-950/20 p-4">
                  <div className="flex items-center gap-2 text-red-400 font-bold text-xs">
                    <ShieldAlert className="h-4 w-4" />
                    Counter-Intelligence & ICA Section 25.3 Order
                  </div>
                  <p className="mt-1 text-xs text-text-muted">
                    Foreign interference audit identifying indirect voting equity or board observation rights held by entities under SASAC (State-owned Assets Supervision and Administration Commission) directive.
                  </p>
                </div>

                <div className="space-y-3">
                  <div className="text-xs font-mono text-text-subtle uppercase">Audited Ownership Chain:</div>
                  <div className="space-y-2 font-mono text-xs">
                    <div className="flex items-center gap-2 p-2.5 rounded-lg border border-red-500/30 bg-card">
                      <span className="font-bold text-red-400">1.</span>
                      <span className="text-white">SASAC (State Council, Beijing)</span>
                      <span className="ml-auto text-[10px] text-red-400 font-bold">SOE ORG</span>
                    </div>
                    <div className="flex items-center justify-center text-text-subtle">
                      <ChevronRight className="h-4 w-4 rotate-90" />
                    </div>
                    <div className="flex items-center gap-2 p-2.5 rounded-lg border border-amber-500/30 bg-card">
                      <span className="font-bold text-amber-400">2.</span>
                      <span className="text-white">Northern Horizon Minerals Ltd (BVI)</span>
                      <span className="ml-auto text-[10px] text-amber-400 font-bold">72% STAKE</span>
                    </div>
                    <div className="flex items-center justify-center text-text-subtle">
                      <ChevronRight className="h-4 w-4 rotate-90" />
                    </div>
                    <div className="flex items-center gap-2 p-2.5 rounded-lg border border-primary/30 bg-card">
                      <span className="font-bold text-aurora">3.</span>
                      <span className="text-white">Crawford Nickel Offtake Right</span>
                      <span className="ml-auto text-[10px] text-aurora font-bold">DIVESTMENT REC</span>
                    </div>
                  </div>
                </div>

                <button
                  onClick={() => alert("ICA Section 25.3 Formal Cabinet Referral Initiated. Case ref: #SOV-CABINET-2026.")}
                  className="w-full py-2.5 rounded-xl bg-red-600 hover:bg-red-500 text-white font-bold text-xs transition-colors flex items-center justify-center gap-2 shadow-lg"
                >
                  <Lock className="h-4 w-4" />
                  Issue Section 25.3 Order in Council Divestment
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
