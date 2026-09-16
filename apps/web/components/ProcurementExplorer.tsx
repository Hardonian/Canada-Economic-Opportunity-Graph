"use client";

import { useState, useMemo } from "react";
import { 
  Search, 
  Filter, 
  ExternalLink, 
  ShieldCheck, 
  Building2, 
  Calendar, 
  RefreshCw, 
  Tag,
  FileText,
  Scale,
  Clock,
  CheckCircle2,
  HardHat,
  Landmark,
  BadgeAlert
} from "lucide-react";
import type { 
  Procurement, 
  FilingsResponse, 
  FilingRecord, 
  EARecord, 
  TenderAmendment 
} from "@/lib/types";
import { CANONICAL_FILINGS_SNAPSHOT } from "@/lib/data";

interface Props {
  initialProcurements: Procurement[];
  initialFilings?: FilingsResponse;
}

export default function ProcurementExplorer({ initialProcurements, initialFilings }: Props) {
  const [procurements, setProcurements] = useState<Procurement[]>(initialProcurements);
  const [filings, setFilings] = useState<FilingsResponse>(initialFilings || CANONICAL_FILINGS_SNAPSHOT);
  const [activeTab, setActiveTab] = useState<"tenders" | "sedar" | "ea" | "amendments">("tenders");
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
      const [pRes, fRes] = await Promise.all([
        fetch("/api/v1/procurements?limit=100"),
        fetch("/api/v1/filings/recent")
      ]);
      if (pRes.ok) {
        const data = await pRes.json();
        if (Array.isArray(data.procurements) && data.procurements.length > 0) {
          setProcurements(data.procurements);
        }
      }
      if (fRes.ok) {
        const fData = await fRes.json();
        if (Array.isArray(fData.recent_disclosures)) {
          setFilings(fData);
        }
      }
      setLastRefreshed(new Date());
      setIsLive(true);
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
            {isLive ? "CONTINUOUS DISCLOSURE & PROCUREMENT RADAR ACTIVE" : "SNAPSHOT FEED"}
          </span>
          <span className="text-text-subtle text-xs">•</span>
          <span className="font-mono text-[11px] text-text-muted">
            {procurements.length} Tenders + {filings.total_count} Continuous Disclosures
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

      {/* Navigation Sub-Tabs */}
      <div className="flex flex-wrap items-center gap-2 border-b border-border/70 pb-3">
        <button
          type="button"
          onClick={() => setActiveTab("tenders")}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-bold transition-all ${
            activeTab === "tenders"
              ? "bg-primary/20 text-aurora border border-primary/40 shadow-sm"
              : "bg-surface/60 text-text-muted hover:text-text-main border border-borderSubtle"
          }`}
        >
          <Building2 className="h-4 w-4" />
          CanadaBuys & Defence Tenders ({procurements.length})
        </button>

        <button
          type="button"
          onClick={() => setActiveTab("sedar")}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-bold transition-all ${
            activeTab === "sedar"
              ? "bg-primary/20 text-aurora border border-primary/40 shadow-sm"
              : "bg-surface/60 text-text-muted hover:text-text-main border border-borderSubtle"
          }`}
        >
          <FileText className="h-4 w-4" />
          SEDAR+ Continuous Disclosures ({filings.recent_disclosures.length})
        </button>

        <button
          type="button"
          onClick={() => setActiveTab("ea")}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-bold transition-all ${
            activeTab === "ea"
              ? "bg-primary/20 text-aurora border border-primary/40 shadow-sm"
              : "bg-surface/60 text-text-muted hover:text-text-main border border-borderSubtle"
          }`}
        >
          <Scale className="h-4 w-4" />
          EA Registries & Decisions ({filings.recent_ea_notices.length})
        </button>

        <button
          type="button"
          onClick={() => setActiveTab("amendments")}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-bold transition-all ${
            activeTab === "amendments"
              ? "bg-primary/20 text-aurora border border-primary/40 shadow-sm"
              : "bg-surface/60 text-text-muted hover:text-text-main border border-borderSubtle"
          }`}
        >
          <Clock className="h-4 w-4" />
          Tender Amendments & Awards ({filings.recent_amendments.length})
        </button>
      </div>

      {/* TAB 1: CANADABUYS & DEFENCE TENDERS */}
      {activeTab === "tenders" && (
        <div className="space-y-6">
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
        </div>
      )}

      {/* TAB 2: SEDAR+ CONTINUOUS DISCLOSURE */}
      {activeTab === "sedar" && (
        <div className="space-y-4">
          <div className="p-4 rounded-xl border border-primary/30 bg-primary/5 text-xs text-text-muted flex items-start gap-3">
            <FileText className="h-5 w-5 text-aurora shrink-0 mt-0.5" />
            <div>
              <span className="font-bold text-text-main">Automated Continuous Disclosure Parser:</span>{" "}
              Streaming Canadian public company MD&A, AIF, and Material Change filings from SEDAR+. Disclosures are scanned for capex adjustments, FID declarations, and major EPC contract awards with deterministic SHA-256 cryptographic provenance.
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {filings.recent_disclosures.map((f: FilingRecord) => (
              <div key={f.id} className="glass-card p-5 rounded-xl border border-border/80 space-y-4 hover:border-aurora/40 transition-all">
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
                  <div className="space-y-1">
                    <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
                      <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-aurora font-bold">
                        {f.ticker}:{f.exchange}
                      </span>
                      <span className="px-2 py-0.5 rounded bg-primary/20 text-text-main font-semibold">
                        {f.filing_type.replace(/_/g, " ")}
                      </span>
                      {f.stage_change_detected && (
                        <span className="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 flex items-center gap-1 font-bold">
                          <CheckCircle2 className="h-3 w-3" /> FID / Stage Advance: {f.detected_stage}
                        </span>
                      )}
                    </div>

                    <h3 className="text-base font-bold text-text-main">
                      <a href={f.source_url} target="_blank" rel="noopener noreferrer" className="hover:text-aurora transition-colors inline-flex items-center gap-1.5">
                        {f.issuer_name} — {f.document_title}
                        <ExternalLink className="h-3.5 w-3.5 text-text-subtle" />
                      </a>
                    </h3>

                    <div className="flex items-center gap-2 text-xs text-text-muted font-mono">
                      <Calendar className="h-3.5 w-3.5 text-text-subtle" />
                      Filing Date: {new Date(f.filing_date).toLocaleDateString("en-CA")}
                    </div>
                  </div>

                  <div className="text-right shrink-0">
                    <div className="text-[10px] font-mono text-text-subtle uppercase">Capex Stated</div>
                    <div className="text-base font-black font-mono text-gold">
                      {f.capex_revision_cad ? formatCAD(f.capex_revision_cad) : "Undisclosed"}
                    </div>
                  </div>
                </div>

                {/* Contract Awards */}
                {f.contract_awards && f.contract_awards.length > 0 && (
                  <div className="p-3 rounded-lg border border-borderSubtle bg-surface/80 space-y-2">
                    <div className="text-xs font-bold text-text-main flex items-center gap-1.5">
                      <HardHat className="h-3.5 w-3.5 text-aurora" />
                      EPC & Engineering Subcontract Awards Disclosed:
                    </div>
                    {f.contract_awards.map((award, idx) => (
                      <div key={idx} className="flex flex-col sm:flex-row sm:items-center justify-between text-xs text-text-muted gap-1 pl-5 border-l-2 border-aurora/40">
                        <div>
                          <span className="font-semibold text-text-main">{award.contractor_name}</span> — {award.scope_of_work}
                        </div>
                        <div className="font-mono text-gold font-bold shrink-0">
                          {formatCAD(award.value_cad)}
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {/* Material Events */}
                {f.material_events && f.material_events.length > 0 && (
                  <div className="space-y-1">
                    <div className="text-[11px] font-mono uppercase text-text-subtle">Key Regulatory & Operational Highlights:</div>
                    <ul className="space-y-1 text-xs text-text-muted">
                      {f.material_events.map((evt, idx) => (
                        <li key={idx} className="flex items-start gap-2">
                          <span className="h-1.5 w-1.5 rounded-full bg-aurora mt-1.5 shrink-0" />
                          <span>{evt}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                <div className="flex flex-wrap items-center justify-between gap-2 pt-2 border-t border-borderSubtle text-[10px] font-mono text-text-subtle">
                  <span>SHA-256 Digest: {f.raw_content_sha256.slice(0, 20)}...</span>
                  <span className="text-aurora">Audit Hash: {f.audit_hash.slice(0, 20)}...</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* TAB 3: ENVIRONMENTAL ASSESSMENT REGISTRIES */}
      {activeTab === "ea" && (
        <div className="space-y-4">
          <div className="p-4 rounded-xl border border-primary/30 bg-primary/5 text-xs text-text-muted flex items-start gap-3">
            <Scale className="h-5 w-5 text-aurora shrink-0 mt-0.5" />
            <div>
              <span className="font-bold text-text-main">Provincial & Federal EA Registries:</span>{" "}
              Live milestone feeds from the BC Environmental Assessment Office (EAO), Ontario Environmental Registry (ERO), Alberta AER, and the Impact Assessment Agency of Canada (IAAC). Tracks active public comment deadlines and certificate approvals.
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {filings.recent_ea_notices.map((n: EARecord) => (
              <div key={n.id} className="glass-card p-5 rounded-xl border border-border/80 space-y-3 hover:border-aurora/40 transition-all">
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
                  <div className="space-y-1">
                    <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
                      <span className="px-2 py-0.5 rounded bg-primary/20 text-aurora font-bold">
                        {n.registry_source}
                      </span>
                      <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-text-main font-semibold">
                        {n.province}
                      </span>
                      <span className="px-2 py-0.5 rounded bg-surface text-text-subtle">
                        {n.registry_project_id}
                      </span>
                      {n.approved && (
                        <span className="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 flex items-center gap-1 font-bold">
                          <CheckCircle2 className="h-3 w-3" /> Certificate Issued / Approved
                        </span>
                      )}
                      {n.comment_deadline && (
                        <span className="px-2 py-0.5 rounded bg-amber-500/20 text-amber-300 border border-amber-500/30 flex items-center gap-1 font-bold">
                          <BadgeAlert className="h-3 w-3" /> Public Comments Close: {new Date(n.comment_deadline).toLocaleDateString("en-CA")}
                        </span>
                      )}
                    </div>

                    <h3 className="text-base font-bold text-text-main">
                      <a href={n.notice_url} target="_blank" rel="noopener noreferrer" className="hover:text-aurora transition-colors inline-flex items-center gap-1.5">
                        {n.project_name} — {n.notice_title}
                        <ExternalLink className="h-3.5 w-3.5 text-text-subtle" />
                      </a>
                    </h3>

                    <p className="text-xs text-text-muted leading-relaxed">
                      {n.summary}
                    </p>
                  </div>

                  <div className="text-right shrink-0">
                    <div className="text-[10px] font-mono text-text-subtle uppercase">Milestone</div>
                    <div className="text-xs font-bold font-mono text-text-main">
                      {n.milestone.replace(/_/g, " ")}
                    </div>
                    {n.conditions_count && (
                      <div className="text-[11px] font-mono text-aurora mt-1">
                        {n.conditions_count} Binding Conditions
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex flex-wrap items-center justify-between gap-2 pt-2 border-t border-borderSubtle text-[10px] font-mono text-text-subtle">
                  <span>Published: {new Date(n.published_date).toLocaleDateString("en-CA")}</span>
                  <span className="text-aurora">Registry Audit Hash: {n.audit_hash.slice(0, 20)}...</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* TAB 4: TENDER AMENDMENTS & AWARDS */}
      {activeTab === "amendments" && (
        <div className="space-y-4">
          <div className="p-4 rounded-xl border border-primary/30 bg-primary/5 text-xs text-text-muted flex items-start gap-3">
            <Clock className="h-5 w-5 text-aurora shrink-0 mt-0.5" />
            <div>
              <span className="font-bold text-text-main">CanadaBuys & DCC Solicitation Amendments:</span>{" "}
              Live tracking of closing date extensions, technical Q&A addenda, and final contract award notices issued to engineering and construction proponents.
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {filings.recent_amendments.map((a: TenderAmendment) => (
              <div key={a.id} className="glass-card p-5 rounded-xl border border-border/80 space-y-3 hover:border-aurora/40 transition-all">
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
                  <div className="space-y-1">
                    <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
                      <span className="px-2 py-0.5 rounded bg-primary/20 text-aurora font-bold">
                        {a.tender_reference}
                      </span>
                      <span className="px-2 py-0.5 rounded bg-surface border border-borderSubtle text-text-main font-semibold">
                        Amendment #{a.amendment_number}
                      </span>
                      <span className="px-2 py-0.5 rounded bg-blue-500/20 text-blue-400 border border-blue-500/30 font-semibold">
                        {a.type.replace(/_/g, " ")}
                      </span>
                    </div>

                    <h3 className="text-sm sm:text-base font-bold text-text-main">
                      <a href={a.source_url} target="_blank" rel="noopener noreferrer" className="hover:text-aurora transition-colors inline-flex items-center gap-1.5">
                        {a.summary}
                        <ExternalLink className="h-3.5 w-3.5 text-text-subtle" />
                      </a>
                    </h3>

                    {a.winning_bidder && (
                      <div className="flex items-center gap-2 text-xs text-text-muted font-mono">
                        <Landmark className="h-3.5 w-3.5 text-aurora" />
                        Winning Bidder: <span className="text-text-main font-bold">{a.winning_bidder}</span>
                        {a.winning_bidder_bn && (
                          <span className="text-text-subtle">(BN: {a.winning_bidder_bn})</span>
                        )}
                      </div>
                    )}
                  </div>

                  <div className="text-right shrink-0">
                    <div className="text-[10px] font-mono text-text-subtle uppercase">Contract Award Value</div>
                    <div className="text-base font-black font-mono text-gold">
                      {a.contract_value_cad ? formatCAD(a.contract_value_cad) : "In Evaluation / Scope Addendum"}
                    </div>
                    {a.revised_closing && (
                      <div className="text-[11px] font-mono text-amber-400 mt-1">
                        Extended: {new Date(a.revised_closing).toLocaleDateString("en-CA")}
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex flex-wrap items-center justify-between gap-2 pt-2 border-t border-borderSubtle text-[10px] font-mono text-text-subtle">
                  <span>Issued: {new Date(a.issued_date).toLocaleDateString("en-CA")}</span>
                  <span className="text-aurora">Amendment Audit Hash: {a.audit_hash.slice(0, 20)}...</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Provenance & Compliance Note */}
      <div className="p-4 rounded-xl border border-borderSubtle bg-surface text-xs leading-relaxed text-text-muted flex items-start gap-3">
        <ShieldCheck className="h-5 w-5 text-aurora shrink-0 mt-0.5" />
        <div>
          <span className="font-bold text-text-main">Official Continuous Disclosure & Government Procurement Provenance:</span>{" "}
          All tender notices, amendments, continuous corporate disclosures, and environmental assessment milestones are ingested from authoritative Open Data, SEDAR+, PSPC CanadaBuys, DCC, and provincial registry systems. Every record maintains SHA-256 cryptographic provenance and immutable change tracking.
        </div>
      </div>
    </div>
  );
}
