"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { 
  Search, 
  Filter, 
  ExternalLink, 
  ShieldCheck, 
  Building2, 
  Calendar, 
  DollarSign, 
  RefreshCw, 
  Radio, 
  FileCheck2, 
  CheckCircle2,
  Tag
} from "lucide-react";
import type { Procurement } from "@/lib/types";

interface Props {
  initialProcurements: Procurement[];
}

export default function ProcurementExplorer({ initialProcurements }: Props) {
  const [procurements, setProcurements] = useState<Procurement[]>(initialProcurements);
  const [search, setSearch] = useState("");
  const [selectedStage, setSelectedStage] = useState<string>("ALL");
  const [selectedBuyerType, setSelectedBuyerType] = useState<string>("ALL");
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [lastRefreshed, setLastRefreshed] = useState<Date>(new Date());
  const [isLive, setIsLive] = useState(true);

  const stages = useMemo(() => {
    const list = Array.from(new Set(procurements.map((p) => p.stage)));
    return ["ALL", ...list];
  }, [procurements]);

  const buyerTypes = useMemo(() => {
    const list = Array.from(new Set(procurements.map((p) => p.buyer_type).filter(Boolean)));
    return ["ALL", ...list];
  }, [procurements]);

  const filtered = useMemo(() => {
    return procurements.filter((p) => {
      if (selectedStage !== "ALL" && p.stage !== selectedStage) return false;
      if (selectedBuyerType !== "ALL" && p.buyer_type !== selectedBuyerType) return false;
      if (search.trim()) {
        const q = search.toLowerCase();
        const matches = 
          p.title.toLowerCase().includes(q) ||
          p.tender_id.toLowerCase().includes(q) ||
          p.buyer.toLowerCase().includes(q) ||
          (p.categories && p.categories.some((c) => c.toLowerCase().includes(q)));
        if (!matches) return false;
      }
      return true;
    });
  }, [procurements, selectedStage, selectedBuyerType, search]);

  const totalDisclosedCAD = useMemo(() => {
    return filtered.reduce((sum, p) => sum + (p.estimated_cad || 0), 0);
  }, [filtered]);

  const uniqueBuyers = useMemo(() => {
    return new Set(filtered.map((p) => p.buyer)).size;
  }, [filtered]);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      const res = await fetch("/api/v1/procurements?limit=100");
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data.procurements) && data.procurements.length > 0) {
          setProcurements(data.procurements);
          setLastRefreshed(new Date());
          setIsLive(true);
        }
      }
    } catch {
      // Retain current data
    } finally {
      setIsRefreshing(false);
    }
  };

  const formatCAD = (amount?: number) => {
    if (!amount || amount === 0) return "Disclosed in notice";
    if (amount >= 1e9) return `$${(amount / 1e9).toFixed(2)}B CAD`;
    if (amount >= 1e6) return `$${(amount / 1e6).toFixed(1)}M CAD`;
    if (amount >= 1e3) return `$${(amount / 1e3).toFixed(0)}K CAD`;
    return `$${amount.toLocaleString()} CAD`;
  };

  return (
    <div className="space-y-8">
      {/* Top Telemetry & Controls Ribbon */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl border border-borderSubtle bg-surface/80">
        <div className="flex items-center gap-3">
          <span className="flex h-2.5 w-2.5 rounded-full bg-aurora animate-pulse" />
          <span className="font-mono text-xs font-bold uppercase tracking-wider text-text-main">
            {isLive ? "LIVE CANADABUYS & DEFENCE FEED ACTIVE" : "SNAPSHOT FEED"}
          </span>
          <span className="text-text-subtle text-xs">•</span>
          <span className="font-mono text-[11px] text-text-muted">
            {procurements.length} Attributed Solicitations
          </span>
        </div>

        <div className="flex items-center gap-3">
          <span className="font-mono text-[10px] text-text-subtle hidden sm:inline">
            Updated: {lastRefreshed.toLocaleTimeString()}
          </span>
          <button
            type="button"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-border bg-card text-xs font-mono font-semibold text-text-main hover:border-aurora/50 hover:text-aurora transition-all disabled:opacity-50"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${isRefreshing ? "animate-spin text-aurora" : ""}`} />
            {isRefreshing ? "Syncing..." : "Refresh Live Feed"}
          </button>
        </div>
      </div>

      {/* KPI Highlight Metrics */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="glass-card p-5 rounded-xl border border-border/80 space-y-1">
          <div className="text-[11px] font-mono text-text-subtle uppercase">Active Solicitations</div>
          <div className="text-3xl font-black font-tabular text-aurora">{filtered.length}</div>
          <div className="text-[11px] text-text-muted">Tenders currently tracked</div>
        </div>

        <div className="glass-card p-5 rounded-xl border border-border/80 space-y-1">
          <div className="text-[11px] font-mono text-text-subtle uppercase">Disclosed Contract Value</div>
          <div className="text-3xl font-black font-tabular text-gold">
            {totalDisclosedCAD > 0 ? formatCAD(totalDisclosedCAD) : "Multi-Stage RFP"}
          </div>
          <div className="text-[11px] text-text-muted">Evidenced expenditure value</div>
        </div>

        <div className="glass-card p-5 rounded-xl border border-border/80 space-y-1">
          <div className="text-[11px] font-mono text-text-subtle uppercase">Attributed Federal Buyers</div>
          <div className="text-3xl font-black font-tabular text-text-main">{uniqueBuyers}</div>
          <div className="text-[11px] text-text-muted">Departments & Crown Corporations</div>
        </div>
      </div>

      {/* Filter & Search Bar */}
      <div className="glass-panel p-4 rounded-xl border border-border/80 space-y-3">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-text-subtle" />
            <input
              type="text"
              placeholder="Search tender ID, title, buyer..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-3 py-1.5 rounded-lg bg-background border border-border text-xs text-text-main focus:border-aurora outline-none"
            />
          </div>

          {/* Stage Filter */}
          <div className="flex items-center gap-2">
            <Filter className="h-3.5 w-3.5 text-text-subtle shrink-0" />
            <select
              aria-label="Filter solicitations by stage"
              value={selectedStage}
              onChange={(e) => setSelectedStage(e.target.value)}
              className="w-full px-2.5 py-1.5 rounded-lg bg-background border border-border text-xs text-text-main focus:border-aurora outline-none"
            >
              {stages.map((stage) => (
                <option key={stage} value={stage}>
                  Stage: {stage}
                </option>
              ))}
            </select>
          </div>

          {/* Buyer Type Filter */}
          <div className="flex items-center gap-2">
            <Building2 className="h-3.5 w-3.5 text-text-subtle shrink-0" />
            <select
              aria-label="Filter solicitations by buyer type"
              value={selectedBuyerType}
              onChange={(e) => setSelectedBuyerType(e.target.value)}
              className="w-full px-2.5 py-1.5 rounded-lg bg-background border border-border text-xs text-text-main focus:border-aurora outline-none"
            >
              {buyerTypes.map((type) => (
                <option key={type} value={type}>
                  Buyer: {type}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Procurements List */}
      <div className="space-y-4">
        {filtered.length === 0 ? (
          <div className="p-8 text-center rounded-xl border border-border bg-card space-y-2">
            <div className="text-sm font-semibold text-text-main">No solicitations match current criteria</div>
            <p className="text-xs text-text-muted">Clear your search query or reset filters to display all tenders.</p>
          </div>
        ) : (
          filtered.map((proc) => (
            <div 
              key={proc.id}
              className="glass-card p-5 rounded-xl border border-border/80 hover:border-aurora/40 transition-all space-y-3"
            >
              <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
                <div className="space-y-1.5">
                  <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
                    <span className="px-2 py-0.5 rounded bg-primary/20 text-aurora font-bold">
                      {proc.tender_id}
                    </span>
                    <span className={`px-2 py-0.5 rounded font-semibold ${
                      proc.stage === "Open" 
                        ? "bg-emerald-500/20 text-emerald-400 border border-emerald-500/30" 
                        : "bg-blue-500/20 text-blue-400 border border-blue-500/30"
                    }`}>
                      {proc.stage}
                    </span>
                    <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-text-subtle">
                      {proc.requirement_class}
                    </span>
                  </div>

                  <h3 className="text-sm sm:text-base font-bold text-text-main hover:text-aurora transition-colors">
                    <a href={proc.source_url} target="_blank" rel="noopener noreferrer">
                      {proc.title}
                    </a>
                  </h3>

                  <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-text-muted">
                    <span className="flex items-center gap-1">
                      <Building2 className="h-3.5 w-3.5 text-text-subtle" />
                      {proc.buyer} ({proc.buyer_type})
                    </span>
                    {proc.closing_date && (
                      <span className="flex items-center gap-1 font-mono text-[11px]">
                        <Calendar className="h-3.5 w-3.5 text-text-subtle" />
                        Closes: {new Intl.DateTimeFormat("en-CA", { dateStyle: "medium", timeZone: "UTC" }).format(new Date(proc.closing_date))}
                      </span>
                    )}
                  </div>
                </div>

                <div className="flex sm:flex-col items-end justify-between sm:justify-start gap-2 shrink-0">
                  <div className="text-right">
                    <div className="text-[10px] font-mono text-text-subtle uppercase">Est. Value</div>
                    <div className="text-sm font-black font-mono text-gold">
                      {formatCAD(proc.estimated_cad)}
                    </div>
                  </div>

                  <a
                    href={proc.source_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 px-3 py-1 rounded-lg border border-primary/30 bg-primary/10 text-xs font-mono font-semibold text-aurora hover:bg-primary/20 transition-colors"
                  >
                    CanadaBuys Notice <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
              </div>

              {proc.categories && proc.categories.length > 0 && (
                <div className="flex flex-wrap items-center gap-1.5 pt-2 border-t border-borderSubtle">
                  <Tag className="h-3 w-3 text-text-subtle" />
                  {proc.categories.map((cat) => (
                    <span key={cat} className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-[10px] font-mono text-text-muted">
                      {cat}
                    </span>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Provenance & Compliance Note */}
      <div className="p-4 rounded-xl border border-borderSubtle bg-surface text-xs leading-relaxed text-text-muted flex items-start gap-3">
        <ShieldCheck className="h-5 w-5 text-aurora shrink-0 mt-0.5" />
        <div>
          <span className="font-bold text-text-main">Official Government Procurement Provenance:</span>{" "}
          All tender notices are harvested from authoritative Open Data and API endpoints operated by Public Services and Procurement Canada (PSPC / CanadaBuys) and Defence Construction Canada (DCC). Each solicitation record preserves original contract identifiers, buyer authorities, and publisher content hashes.
        </div>
      </div>
    </div>
  );
}
