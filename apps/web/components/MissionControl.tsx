"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import {
  ArrowUpRight,
  BadgeCheck,
  BookOpenCheck,
  Boxes,
  BriefcaseBusiness,
  CheckCircle2,
  ChevronRight,
  CircleDollarSign,
  ClipboardCheck,
  Download,
  FileLock2,
  FileSearch,
  Fingerprint,
  Gauge,
  Landmark,
  Network,
  RadioTower,
  ShieldCheck,
  Sparkles,
  Target,
  UsersRound,
  Workflow,
} from "lucide-react";
import type { Project, RadarStats, Signal } from "@/lib/types";

type WorkspaceWireState = {
  version: 1;
  selected_project_id: string;
  saved_project_ids: string[];
  scenario: {
    catalytic_capital_cad: number;
    permit_months_saved: number;
    domestic_procurement_pct: number;
    indigenous_equity_pct: number;
  };
  decisions: unknown[];
  updated_at: string;
};

type WorkspaceResponse = {
  workspace: WorkspaceWireState;
  security_mode?: "CONFIGURED" | "DEMO";
  error?: string;
};

type Mission = Project & {
  priority: number;
  evidenceScore: number;
  readinessScore: number;
};

const money = new Intl.NumberFormat("en-CA", {
  style: "currency",
  currency: "CAD",
  notation: "compact",
  maximumFractionDigits: 1,
});

const dateTime = new Intl.DateTimeFormat("en-CA", {
  month: "short",
  day: "numeric",
  timeZone: "UTC",
});

const activeStages = new Set(["FINANCING", "PROCUREMENT", "FID_LIKELY", "FID", "CONSTRUCTION"]);

function projectMission(project: Project): Mission {
  const evidenceScore = Math.min(30, (project.evidence?.length ?? 0) * 6);
  const scoreCoverage = Math.min(15, (project.score_details?.length ?? 0) * 5);
  const readinessScore = activeStages.has(project.current_stage) ? 30 : project.current_stage === "PERMITTING" ? 22 : 14;
  const buildability = Math.round((project.scores?.buildability ?? 0) * 0.25);
  const scale = project.capex_cad > 0 ? Math.min(15, Math.round(Math.log10(project.capex_cad / 1_000_000 + 1) * 4)) : 0;
  return { ...project, priority: Math.min(100, evidenceScore + scoreCoverage + readinessScore + buildability + scale), evidenceScore, readinessScore };
}

function statusTone(value: string): string {
  if (["VERIFIED", "SUPPORTED", "HEALTHY", "LIVE_UPSTREAM_API"].includes(value)) return "border-primary/40 bg-primary/10 text-aurora";
  if (["REPORTED", "BUNDLED_REVIEWED_SNAPSHOT"].includes(value)) return "border-gold/40 bg-gold/10 text-gold";
  return "border-border bg-surface text-text-muted";
}

function stageLabel(value: string): string {
  return value.replaceAll("_", " ").toLowerCase().replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function formatDate(value: string): string {
  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? "Timestamp unavailable" : dateTime.format(parsed);
}

function signalConfidence(value: number): string {
  const normalized = value <= 1 ? value * 100 : value;
  return `${Math.max(0, Math.min(100, Math.round(normalized)))}%`;
}

function MetricCard({ label, value, detail, tone = "aurora" }: { label: string; value: string; detail: string; tone?: "aurora" | "gold" | "neutral" }) {
  const valueClass = tone === "gold" ? "text-gold" : tone === "neutral" ? "text-text-main" : "text-aurora";
  return (
    <div className="rounded-2xl border border-border/80 bg-background/60 p-4 shadow-inner">
      <p className="font-mono text-[10px] font-bold uppercase tracking-[0.13em] text-text-subtle">{label}</p>
      <p className={`mt-2 text-2xl font-black tracking-tight tabular-nums ${valueClass}`}>{value}</p>
      <p className="mt-1 text-[11px] leading-relaxed text-text-muted">{detail}</p>
    </div>
  );
}

function EvidenceRail({ mission }: { mission: Mission }) {
  const evidence = mission.evidence ?? [];
  const scores = mission.score_details ?? [];
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <span className="font-mono text-[10px] uppercase tracking-wider text-text-subtle">Public evidence rail</span>
        <span className="text-xs font-bold text-aurora">{evidence.length} linked record{evidence.length === 1 ? "" : "s"}</span>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-background">
        <div className="h-full rounded-full bg-aurora transition-all" style={{ width: `${Math.min(100, 18 + evidence.length * 16 + scores.length * 8)}%` }} />
      </div>
      <div className="grid gap-2 sm:grid-cols-2">
        <div className="rounded-xl border border-borderSubtle bg-background/70 p-3">
          <p className="font-mono text-[10px] text-text-subtle">EVIDENCE RECORDS</p>
          <p className="mt-1 text-lg font-black text-text-main">{evidence.length}</p>
          <p className="mt-1 text-[11px] text-text-muted">Source URLs and content hashes travel with the project record.</p>
        </div>
        <div className="rounded-xl border border-borderSubtle bg-background/70 p-3">
          <p className="font-mono text-[10px] text-text-subtle">SCORING COVERAGE</p>
          <p className="mt-1 text-lg font-black text-text-main">{scores.length}</p>
          <p className="mt-1 text-[11px] text-text-muted">Deterministic score records are inspectable, not black-box recommendations.</p>
        </div>
      </div>
      {evidence.length > 0 && (
        <div className="rounded-xl border border-borderSubtle bg-surface/70 p-3">
          <p className="font-mono text-[10px] uppercase tracking-wider text-text-subtle">Latest linked publisher</p>
          <a href={evidence[0].source_url} target="_blank" rel="noreferrer" className="mt-1 inline-flex max-w-full items-center gap-1 text-xs font-semibold text-aurora hover:underline">
            <span className="truncate">{evidence[0].publisher}</span><ArrowUpRight className="h-3.5 w-3.5 shrink-0" />
          </a>
          <p className="mt-1 font-mono text-[10px] text-text-subtle">SHA-256 {evidence[0].content_hash.slice(0, 16)}…</p>
        </div>
      )}
    </div>
  );
}

export default function MissionControl({
  projects,
  signals,
  stats,
  sourceMode,
  generatedAt,
}: {
  projects: Project[];
  signals: Signal[];
  stats: RadarStats;
  sourceMode: "LIVE_UPSTREAM_API" | "BUNDLED_REVIEWED_SNAPSHOT";
  generatedAt?: string;
}) {
  const sectors = useMemo(() => ["All", ...Array.from(new Set(projects.map((project) => project.sector))).sort()], [projects]);
  const [sector, setSector] = useState("All");
  const [selectedID, setSelectedID] = useState("");
  const [queueState, setQueueState] = useState<"idle" | "saving" | "saved" | "error">("idle");
  const [queueMessage, setQueueMessage] = useState("");
  const [briefState, setBriefState] = useState("idle");

  const missions = useMemo(() => projects.map(projectMission).sort((a, b) => b.priority - a.priority), [projects]);
  const filteredMissions = useMemo(() => sector === "All" ? missions : missions.filter((mission) => mission.sector === sector), [missions, sector]);
  const selected = filteredMissions.find((mission) => mission.id === selectedID) ?? filteredMissions[0] ?? missions[0];
  const linkedSignals = useMemo(() => signals.filter((signal) => signal.project_id === selected?.id || signal.project_name === selected?.name).slice(0, 4), [selected, signals]);
  const displayedSignals = useMemo(() => (linkedSignals.length ? linkedSignals : signals).slice(0, 5), [linkedSignals, signals]);

  async function queueForExecution() {
    if (!selected) return;
    setQueueState("saving");
    setQueueMessage("");
    try {
      const currentResponse = await fetch("/api/workspace", { cache: "no-store" });
      const current = await currentResponse.json() as WorkspaceResponse;
      if (!currentResponse.ok || !current.workspace) throw new Error(current.error || "Execution workspace is unavailable.");
      const saved = [...new Set([...current.workspace.saved_project_ids, selected.id])].slice(-12);
      const next: WorkspaceWireState = { ...current.workspace, selected_project_id: selected.id, saved_project_ids: saved };
      const saveResponse = await fetch("/api/workspace", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(next),
      });
      const savedPayload = await saveResponse.json() as WorkspaceResponse;
      if (!saveResponse.ok) throw new Error(savedPayload.error || "Unable to queue the mission.");
      setQueueState("saved");
      setQueueMessage(savedPayload.security_mode === "DEMO" ? "Queued in a demo workspace. Configure COG_WORKSPACE_SECRET before production use." : "Queued in the encrypted execution workspace.");
    } catch (error) {
      setQueueState("error");
      setQueueMessage(error instanceof Error ? error.message : "Unable to queue the mission.");
    }
  }

  function downloadBrief() {
    if (!selected) return;
    const brief = {
      format: "CEGS Mission Brief v1",
      generated_at: new Date().toISOString(),
      caveat: "A transparent public-evidence briefing aid; it is not investment, legal, procurement, Indigenous consultation, or policy advice.",
      mission: {
        id: selected.id,
        name: selected.name,
        sector: selected.sector,
        location: `${selected.location_name}, ${selected.province}`,
        stage: selected.current_stage,
        reported_capex_cad: selected.capex_cad,
        confidence: selected.confidence,
        priority_index: selected.priority,
        evidence: (selected.evidence ?? []).map((record) => ({ publisher: record.publisher, url: record.source_url, content_hash: record.content_hash, retrieved_at: record.retrieval_timestamp })),
      },
    };
    const url = URL.createObjectURL(new Blob([JSON.stringify(brief, null, 2)], { type: "application/json" }));
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `${selected.slug || selected.id}-mission-brief.json`;
    anchor.click();
    URL.revokeObjectURL(url);
    setBriefState("downloaded");
    window.setTimeout(() => setBriefState("idle"), 2400);
  }

  if (!selected) {
    return <div className="mx-auto max-w-5xl px-4 py-16 text-text-muted">No public mission records are available in this release.</div>;
  }

  return (
    <div className="mx-auto max-w-7xl space-y-7 px-4 py-8 sm:px-6 lg:px-8">
      <section className="relative overflow-hidden rounded-3xl border border-primary/35 bg-gradient-to-br from-[#0f2b20] via-card to-[#07110d] p-6 shadow-2xl shadow-black/50 sm:p-9">
        <div aria-hidden="true" className="pointer-events-none absolute -right-24 -top-24 h-80 w-80 rounded-full bg-aurora/15 blur-3xl" />
        <div aria-hidden="true" className="pointer-events-none absolute -bottom-36 left-1/3 h-80 w-80 rounded-full bg-gold/10 blur-3xl" />
        <div className="relative grid gap-8 xl:grid-cols-[1.18fr_0.82fr] xl:items-end">
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/40 bg-background/70 px-3 py-1 font-mono text-[10px] font-bold tracking-[0.15em] text-aurora">
              <RadioTower className="h-3.5 w-3.5" /> NATIONAL MISSION CONTROL // CEGS
            </div>
            <h1 className="mt-4 max-w-4xl text-4xl font-black tracking-tight text-text-main sm:text-6xl">
              See the economic system.<br /><span className="text-aurora">Move the next decision.</span>
            </h1>
            <p className="mt-5 max-w-3xl text-sm leading-relaxed text-text-muted sm:text-base">
              CanadaOpportunityGraph connects public evidence, transparent prioritization, capital workflows, and secure collaboration into a single Canadian operating loop. No opaque recommendation engine. Every move starts with a traceable record.
            </p>
            <div className="mt-6 flex flex-wrap gap-3">
              <a href="#mission-queue" className="inline-flex min-h-11 items-center gap-2 rounded-xl bg-primary px-4 text-sm font-black text-[#050b08] shadow-lg shadow-emerald-950/40 transition hover:bg-aurora-mint">
                <Target className="h-4 w-4" /> Explore mission queue
              </a>
              <Link href="/execution" className="inline-flex min-h-11 items-center gap-2 rounded-xl border border-border bg-background/65 px-4 text-sm font-bold text-text-main transition hover:border-primary hover:text-aurora">
                Open execution workspace <ArrowUpRight className="h-4 w-4" />
              </Link>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <MetricCard label="Tracked capital" value={money.format(stats.total_capex_cad)} detail="Reported project capital in the current graph." />
            <MetricCard label="Decision surface" value={String(stats.total_projects)} detail="Projects visible to the public evidence layer." tone="neutral" />
            <MetricCard label="Live procurement" value={String(stats.active_procurements_count)} detail="Open tender records surfaced for operator review." tone="gold" />
            <div className="rounded-2xl border border-border/80 bg-background/60 p-4 shadow-inner">
              <p className="font-mono text-[10px] font-bold uppercase tracking-[0.13em] text-text-subtle">Data posture</p>
              <span className={`mt-2 inline-flex rounded-full border px-2.5 py-1 font-mono text-[10px] font-bold ${statusTone(sourceMode)}`}>{sourceMode === "LIVE_UPSTREAM_API" ? "LIVE UPSTREAM" : "REVIEWED SNAPSHOT"}</span>
              <p className="mt-2 text-[11px] leading-relaxed text-text-muted">{generatedAt ? `Release timestamp: ${formatDate(generatedAt)} UTC.` : "Data lineage is exposed with every public record."}</p>
            </div>
          </div>
        </div>
      </section>

      <section className="grid gap-5 lg:grid-cols-[0.8fr_1.2fr]">
        <article className="glass-panel rounded-2xl border border-border/80 p-5">
          <div className="flex items-center gap-2">
            <Fingerprint className="h-4 w-4 text-aurora" />
            <h2 className="text-sm font-black text-text-main">The Canadian differentiation</h2>
          </div>
          <p className="mt-3 text-xs leading-relaxed text-text-muted">A national graph wins trust only when its evidence is inspectable, its models are bounded, its collaboration is sovereign, and its interfaces are designed for the institutions that must act together.</p>
          <div className="mt-5 space-y-3">
            {[
              ["Evidence before narrative", "Source URL, retrieval time, content hash, and confidence travel with public records.", BookOpenCheck],
              ["Federated, not extractive", "Public graph intelligence is distinct from tenant-isolated private deal rooms.", Network],
              ["Action with accountability", "Execution choices carry a rationale rather than disappearing into slides and inboxes.", ClipboardCheck],
            ].map(([title, copy, Icon]) => {
              const FeatureIcon = Icon as typeof BookOpenCheck;
              return <div key={title as string} className="flex gap-3 rounded-xl border border-borderSubtle bg-surface/65 p-3"><FeatureIcon className="mt-0.5 h-4 w-4 shrink-0 text-aurora" /><div><h3 className="text-xs font-bold text-text-main">{title as string}</h3><p className="mt-1 text-[11px] leading-relaxed text-text-muted">{copy as string}</p></div></div>;
            })}
          </div>
        </article>
        <article className="rounded-2xl border border-gold/30 bg-gradient-to-br from-gold/10 via-surface to-card p-5">
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-start">
            <div>
              <div className="inline-flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-gold"><Workflow className="h-3.5 w-3.5" /> THE OPERATING LOOP</div>
              <h2 className="mt-2 text-xl font-black text-text-main">From the first signal to a defensible move.</h2>
            </div>
            <Link href="/operations" className="inline-flex items-center gap-1 text-xs font-bold text-gold hover:text-text-main">System health <ArrowUpRight className="h-3.5 w-3.5" /></Link>
          </div>
          <ol className="mt-5 grid gap-3 sm:grid-cols-4">
            {[
              ["01", "Observe", "Public source and signal ingestion", RadioTower],
              ["02", "Verify", "Lineage, evidence, and score coverage", BadgeCheck],
              ["03", "Decide", "Scenario and rationale in an encrypted workspace", Target],
              ["04", "Mobilize", "Tenant-controlled documents and collaborators", FileLock2],
            ].map(([number, title, copy, Icon]) => {
              const StepIcon = Icon as typeof RadioTower;
              return <li key={number as string} className="rounded-xl border border-border/80 bg-background/60 p-3"><StepIcon className="h-4 w-4 text-gold" /><p className="mt-3 font-mono text-[10px] font-bold text-gold">{number as string}</p><h3 className="mt-1 text-sm font-black text-text-main">{title as string}</h3><p className="mt-1 text-[11px] leading-relaxed text-text-muted">{copy as string}</p></li>;
            })}
          </ol>
        </article>
      </section>

      <section id="mission-queue" className="scroll-mt-36 rounded-3xl border border-border/80 bg-card/70 p-5 shadow-2xl shadow-black/20 sm:p-6">
        <div className="flex flex-col gap-4 border-b border-borderSubtle pb-5 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <div className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-aurora"><Gauge className="h-3.5 w-3.5" /> Transparent mission queue</div>
            <h2 className="mt-2 text-2xl font-black text-text-main">Where operators can start, with the reasoning in view.</h2>
            <p className="mt-2 max-w-3xl text-xs leading-relaxed text-text-muted">Priority Index is a visible routing heuristic: public evidence coverage, score coverage, lifecycle readiness, reported scale, and buildability. It does not predict outcomes or make an investment, procurement, or policy decision.</p>
          </div>
          <div className="flex max-w-full flex-wrap gap-2" role="group" aria-label="Filter mission queue by sector">
            {sectors.map((item) => <button key={item} type="button" onClick={() => setSector(item)} aria-pressed={sector === item} className={`min-h-9 rounded-lg border px-3 text-[11px] font-semibold transition ${sector === item ? "border-primary/60 bg-primary/10 text-aurora" : "border-borderSubtle bg-surface text-text-muted hover:border-primary/50 hover:text-text-main"}`}>{item}</button>)}
          </div>
        </div>

        <div className="mt-5 grid gap-5 xl:grid-cols-[0.76fr_1.24fr]">
          <div className="max-h-[660px] space-y-2 overflow-y-auto pr-1" aria-label="Mission queue">
            {filteredMissions.slice(0, 18).map((mission, index) => {
              const active = selected.id === mission.id;
              return <button key={mission.id} type="button" onClick={() => { setSelectedID(mission.id); setQueueState("idle"); setQueueMessage(""); }} className={`w-full rounded-xl border p-3 text-left transition ${active ? "border-primary/70 bg-primary/10 shadow-[inset_3px_0_0_#00D488]" : "border-borderSubtle bg-surface/60 hover:border-primary/45 hover:bg-surface"}`}>
                <div className="flex items-start justify-between gap-3"><div className="min-w-0"><p className="font-mono text-[10px] font-bold text-text-subtle">QUEUE {String(index + 1).padStart(2, "0")} · {mission.sector.toUpperCase()}</p><h3 className="mt-1 truncate text-sm font-bold text-text-main">{mission.name}</h3><p className="mt-1 truncate text-[11px] text-text-muted">{mission.location_name}, {mission.province} · {stageLabel(mission.current_stage)}</p></div><span className="shrink-0 rounded-lg border border-primary/35 bg-background px-2 py-1 text-right"><span className="block font-mono text-[9px] text-text-subtle">INDEX</span><strong className="text-sm text-aurora">{mission.priority}</strong></span></div>
                <div className="mt-3 flex items-center justify-between gap-3 text-[11px]"><span className={`rounded border px-1.5 py-0.5 font-mono ${statusTone(mission.confidence)}`}>{mission.confidence}</span><span className="font-mono text-text-muted">{money.format(mission.capex_cad)}</span></div>
              </button>;
            })}
          </div>

          <article className="min-w-0 rounded-2xl border border-primary/30 bg-gradient-to-b from-[#10261c] to-surface p-5 sm:p-6">
            <div className="flex flex-col justify-between gap-4 border-b border-borderSubtle pb-5 sm:flex-row sm:items-start">
              <div className="min-w-0"><p className="font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-aurora">Mission dossier</p><h3 className="mt-2 text-2xl font-black tracking-tight text-text-main">{selected.name}</h3><p className="mt-2 text-xs leading-relaxed text-text-muted">{selected.summary}</p></div><span className={`w-fit shrink-0 rounded-full border px-2.5 py-1 font-mono text-[10px] font-bold ${statusTone(selected.confidence)}`}>{selected.confidence}</span>
            </div>
            <div className="mt-5 grid gap-5 lg:grid-cols-[1.02fr_0.98fr]">
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-3"><MetricCard label="Reported CAPEX" value={money.format(selected.capex_cad)} detail={selected.capex_status ? `${selected.capex_status.toLowerCase()} capital status` : "Source-reported capital"} /><MetricCard label="Priority index" value={`${selected.priority}/100`} detail="Transparent routing heuristic" tone="gold" /></div>
                <div className="rounded-xl border border-borderSubtle bg-background/55 p-4"><div className="flex items-center justify-between"><h4 className="text-xs font-bold text-text-main">Why it is on deck</h4><Sparkles className="h-4 w-4 text-gold" /></div><ul className="mt-3 space-y-2 text-xs text-text-muted"><li className="flex gap-2"><CheckCircle2 className="mt-0.5 h-3.5 w-3.5 shrink-0 text-aurora" /><span>{selected.evidenceScore} points from linked public evidence.</span></li><li className="flex gap-2"><CheckCircle2 className="mt-0.5 h-3.5 w-3.5 shrink-0 text-aurora" /><span>{selected.readinessScore} points from the recorded {stageLabel(selected.current_stage)} lifecycle stage.</span></li><li className="flex gap-2"><CheckCircle2 className="mt-0.5 h-3.5 w-3.5 shrink-0 text-aurora" /><span>{selected.score_details?.length ?? 0} inspectable score record{(selected.score_details?.length ?? 0) === 1 ? "" : "s"} associated with the project.</span></li></ul></div>
                <div className="flex flex-wrap gap-2"><Link href={`/projects/${encodeURIComponent(selected.slug)}`} className="inline-flex min-h-10 items-center gap-2 rounded-xl border border-primary/40 bg-primary/10 px-3 text-xs font-bold text-aurora hover:bg-primary/20"><FileSearch className="h-4 w-4" /> Open dossier</Link><button type="button" onClick={downloadBrief} className="inline-flex min-h-10 items-center gap-2 rounded-xl border border-border bg-background px-3 text-xs font-bold text-text-main hover:border-primary hover:text-aurora"><Download className="h-4 w-4" /> {briefState === "downloaded" ? "Brief exported" : "Export JSON brief"}</button></div>
              </div>
              <EvidenceRail mission={selected} />
            </div>
            <div className="mt-5 rounded-2xl border border-gold/35 bg-gold/10 p-4"><div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"><div><div className="flex items-center gap-2 text-xs font-black text-text-main"><BriefcaseBusiness className="h-4 w-4 text-gold" /> Move from intelligence to accountable execution</div><p aria-live="polite" className="mt-1 max-w-xl text-[11px] leading-relaxed text-text-muted">{queueMessage || "Queue this public record into the encrypted execution workspace, then escalate sensitive materials only inside the private deal room."}</p></div><div className="flex shrink-0 flex-wrap gap-2"><button type="button" disabled={queueState === "saving"} onClick={queueForExecution} className="inline-flex min-h-10 items-center gap-2 rounded-xl bg-gold px-3 text-xs font-black text-[#1f1600] transition hover:bg-[#ffe09a] disabled:opacity-60"><Boxes className="h-4 w-4" /> {queueState === "saving" ? "Queuing…" : queueState === "saved" ? "Queued" : "Queue to execution"}</button><Link href="/deal-room" className="inline-flex min-h-10 items-center gap-2 rounded-xl border border-gold/50 bg-background/70 px-3 text-xs font-bold text-gold hover:bg-background"><FileLock2 className="h-4 w-4" /> Private deal room</Link></div></div></div>
          </article>
        </div>
      </section>

      <section className="grid gap-5 lg:grid-cols-[1fr_0.96fr]">
        <article className="glass-panel rounded-2xl border border-border/80 p-5 sm:p-6"><div className="flex items-center justify-between gap-3"><div><div className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-aurora"><RadioTower className="h-3.5 w-3.5" /> Signal ledger</div><h2 className="mt-2 text-xl font-black text-text-main">Signals with their context intact.</h2></div><Link href="/procurement" className="inline-flex items-center gap-1 text-xs font-bold text-aurora hover:text-text-main">View activity <ChevronRight className="h-3.5 w-3.5" /></Link></div><div className="mt-5 space-y-2">{displayedSignals.length ? displayedSignals.map((signal) => <div key={signal.id} className="rounded-xl border border-borderSubtle bg-surface/65 p-3"><div className="flex items-start justify-between gap-3"><div className="min-w-0"><p className="font-mono text-[10px] font-bold uppercase text-text-subtle">{signal.type.replaceAll("_", " ")} · {formatDate(signal.timestamp)} UTC</p><p className="mt-1 text-xs font-semibold text-text-main">{signal.project_name}</p><p className="mt-1 text-[11px] leading-relaxed text-text-muted">{signal.description}</p></div><span className="shrink-0 rounded-md border border-border bg-background px-2 py-1 font-mono text-[10px] text-aurora">{signalConfidence(signal.confidence)}</span></div></div>) : <p className="rounded-xl border border-borderSubtle bg-surface p-4 text-xs text-text-muted">No signals are available in the current public release.</p>}</div></article>
        <article className="rounded-2xl border border-primary/30 bg-gradient-to-br from-primary/10 via-card to-surface p-5 sm:p-6"><div className="flex items-center gap-2 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-aurora"><ShieldCheck className="h-3.5 w-3.5" /> Institutional-grade by default</div><h2 className="mt-2 text-xl font-black text-text-main">An adoption model that respects the boundary.</h2><div className="mt-5 grid gap-3 sm:grid-cols-2"><div className="rounded-xl border border-borderSubtle bg-background/65 p-4"><Landmark className="h-4 w-4 text-aurora" /><h3 className="mt-3 text-sm font-bold text-text-main">Public truth layer</h3><p className="mt-1 text-[11px] leading-relaxed text-text-muted">Open CEGS records, APIs, source lineage, and exportable evidence turn the graph into a shared reference point.</p><Link href="/cegs" className="mt-3 inline-flex items-center gap-1 text-xs font-bold text-aurora hover:underline">Explore the standard <ArrowUpRight className="h-3.5 w-3.5" /></Link></div><div className="rounded-xl border border-borderSubtle bg-background/65 p-4"><UsersRound className="h-4 w-4 text-gold" /><h3 className="mt-3 text-sm font-bold text-text-main">Private action layer</h3><p className="mt-1 text-[11px] leading-relaxed text-text-muted">OIDC/SAML identity, tenant membership, clearance-aware documents, audit events, and retention controls keep collaboration deliberate.</p><Link href="/deal-room" className="mt-3 inline-flex items-center gap-1 text-xs font-bold text-gold hover:underline">Open deal room <ArrowUpRight className="h-3.5 w-3.5" /></Link></div></div><div className="mt-4 rounded-xl border border-primary/25 bg-background/55 p-3 text-[11px] leading-relaxed text-text-muted"><strong className="text-text-main">The compounding advantage:</strong> a public evidence graph earns trust and integrations; controlled action spaces turn that trust into repeatable institutional workflows without exposing confidential data.</div></article>
      </section>
    </div>
  );
}
