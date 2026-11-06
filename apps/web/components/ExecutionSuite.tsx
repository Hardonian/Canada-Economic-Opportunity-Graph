"use client";

import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  ArrowUpRight,
  BadgeCheck,
  Bookmark,
  CheckCircle2,
  ChevronRight,
  CircleDollarSign,
  Clock3,
  FileJson2,
  Globe2,
  Landmark,
  LockKeyhole,
  Network,
  Plus,
  Save,
  ShieldCheck,
  Sparkles,
  Target,
  UsersRound,
} from "lucide-react";
import type { Project } from "@/lib/types";
import type { DecisionStatus, WorkspaceState } from "@/lib/workspace-session";

type SecurityMode = "CONFIGURED" | "DEMO";

const INITIAL_WORKSPACE: WorkspaceState = {
  version: 1,
  selected_project_id: "",
  saved_project_ids: [],
  scenario: {
    catalytic_capital_cad: 100_000_000,
    permit_months_saved: 6,
    domestic_procurement_pct: 55,
    indigenous_equity_pct: 10,
  },
  decisions: [],
  updated_at: "",
};

const cad = new Intl.NumberFormat("en-CA", { style: "currency", currency: "CAD", notation: "compact", maximumFractionDigits: 1 });
const preciseCad = new Intl.NumberFormat("en-CA", { style: "currency", currency: "CAD", maximumFractionDigits: 0 });

function statusTheme(status: DecisionStatus) {
  if (status === "ADVANCE") return "border-emerald-400/40 bg-emerald-500/10 text-emerald-300";
  if (status === "INVESTIGATE") return "border-amber-400/40 bg-amber-500/10 text-amber-200";
  return "border-sky-400/40 bg-sky-500/10 text-sky-200";
}

export default function ExecutionSuite({ projects }: { projects: Project[] }) {
  const [workspace, setWorkspace] = useState<WorkspaceState>(INITIAL_WORKSPACE);
  const [mode, setMode] = useState<SecurityMode>("DEMO");
  const [ready, setReady] = useState(false);
  const [saving, setSaving] = useState(false);
  const [notice, setNotice] = useState("");
  const [decisionStatus, setDecisionStatus] = useState<DecisionStatus>("INVESTIGATE");
  const [rationale, setRationale] = useState("");

  const selectedProject = useMemo(
    () => projects.find((project) => project.id === workspace.selected_project_id || project.slug === workspace.selected_project_id) ?? projects[0],
    [projects, workspace.selected_project_id],
  );

  const persist = useCallback(async (next: WorkspaceState, message = "Workspace saved") => {
    setWorkspace(next);
    setSaving(true);
    try {
      const response = await fetch("/api/workspace", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(next),
      });
      const payload = await response.json();
      if (!response.ok) throw new Error(payload.error || "Unable to save workspace.");
      setWorkspace(payload.workspace);
      setMode(payload.security_mode);
      setNotice(message);
    } catch (error) {
      setNotice(error instanceof Error ? error.message : "Unable to save workspace.");
    } finally {
      setSaving(false);
    }
  }, []);

  useEffect(() => {
    let active = true;
    fetch("/api/workspace", { cache: "no-store" })
      .then(async (response) => {
        const payload = await response.json();
        if (!response.ok) throw new Error(payload.error || "Unable to load workspace.");
        return payload;
      })
      .then((payload) => {
        if (!active) return;
        setWorkspace(payload.workspace);
        setMode(payload.security_mode);
      })
      .catch(() => setNotice("Workspace session is unavailable; your scenario has not been persisted."))
      .finally(() => active && setReady(true));
    return () => { active = false; };
  }, []);

  useEffect(() => {
    if (!ready || workspace.selected_project_id || !projects[0]) return;
    const next = { ...workspace, selected_project_id: projects[0].id };
    void persist(next, "Default project loaded");
  }, [projects, persist, ready, workspace]);

  const scenario = workspace.scenario;
  const buildability = selectedProject?.scores?.buildability ?? 45;
  const capex = selectedProject?.capex_cad ?? 0;
  const capitalMultiplier = 2.2 + buildability / 100;
  const capitalMobilized = scenario.catalytic_capital_cad * capitalMultiplier;
  const domesticSpend = capex * (scenario.domestic_procurement_pct / 100);
  const indigenousEquity = capex * (scenario.indigenous_equity_pct / 100);
  const personYears = Math.round((capex / 1_000_000) * (selectedProject?.sector === "Nuclear & Clean Power" ? 7.5 : 5.9));
  const isSaved = Boolean(selectedProject && workspace.saved_project_ids.includes(selectedProject.id));

  function updateScenario<K extends keyof WorkspaceState["scenario"]>(key: K, value: number) {
    const next = { ...workspace, scenario: { ...workspace.scenario, [key]: value } };
    setWorkspace(next);
  }

  function saveScenario() {
    void persist({ ...workspace, selected_project_id: selectedProject?.id ?? "" }, "Scenario encrypted in this browser session");
  }

  function toggleSavedProject() {
    if (!selectedProject) return;
    const ids = isSaved
      ? workspace.saved_project_ids.filter((id) => id !== selectedProject.id)
      : [...workspace.saved_project_ids, selectedProject.id].slice(-12);
    void persist({ ...workspace, selected_project_id: selectedProject.id, saved_project_ids: ids }, isSaved ? "Project removed from execution list" : "Project added to execution list");
  }

  function recordDecision() {
    if (!selectedProject || rationale.trim().length < 12) {
      setNotice("Add a concise rationale (at least 12 characters) to create an accountable decision record.");
      return;
    }
    const next = {
      ...workspace,
      selected_project_id: selectedProject.id,
      decisions: [{
        id: `decision-${Date.now().toString(36)}`,
        project_id: selectedProject.id,
        status: decisionStatus,
        rationale: rationale.trim(),
        created_at: new Date().toISOString(),
      }, ...workspace.decisions].slice(0, 12),
    };
    setRationale("");
    void persist(next, "Decision recorded with project context");
  }

  return (
    <div className="mx-auto max-w-7xl space-y-7 px-4 py-8 sm:px-6 lg:px-8">
      <section className="relative overflow-hidden rounded-3xl border border-primary/35 bg-gradient-to-br from-card via-[#10261c] to-[#07110d] p-6 shadow-2xl shadow-black/40 sm:p-9">
        <div className="pointer-events-none absolute -right-28 -top-28 h-80 w-80 rounded-full bg-aurora/15 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-32 left-1/3 h-72 w-72 rounded-full bg-gold/10 blur-3xl" />
        <div className="relative grid gap-7 lg:grid-cols-[1.4fr_0.8fr] lg:items-end">
          <div>
            <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-primary/35 bg-background/70 px-3 py-1 text-[10px] font-bold tracking-[0.16em] text-aurora">
              <Sparkles className="h-3.5 w-3.5" /> ECONOMIC EXECUTION SUITE
            </div>
            <h1 className="max-w-3xl text-3xl font-black tracking-tight text-text-main sm:text-5xl">
              Turn opportunity intelligence into <span className="text-aurora">accountable execution.</span>
            </h1>
            <p className="mt-4 max-w-2xl text-sm leading-relaxed text-text-muted">
              A decision workspace for capital, policy, procurement, and community partners. Model the intervention, preserve the rationale, and keep the underlying public evidence one click away.
            </p>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div className="rounded-2xl border border-border/80 bg-background/60 p-4">
              <div className="font-mono text-[10px] uppercase tracking-wide text-text-subtle">Projects in session</div>
              <div className="mt-1 text-3xl font-black text-text-main">{workspace.saved_project_ids.length}</div>
              <div className="mt-1 text-[11px] text-text-muted">Curated execution list</div>
            </div>
            <div className="rounded-2xl border border-border/80 bg-background/60 p-4">
              <div className="font-mono text-[10px] uppercase tracking-wide text-text-subtle">Decision records</div>
              <div className="mt-1 text-3xl font-black text-aurora">{workspace.decisions.length}</div>
              <div className="mt-1 text-[11px] text-text-muted">Rationale required</div>
            </div>
          </div>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[0.82fr_1.18fr]">
        <section className="glass-panel rounded-2xl border border-border/80 p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3 border-b border-borderSubtle pb-4">
            <div className="flex items-center gap-2">
              <Target className="h-4 w-4 text-aurora" />
              <h2 className="text-sm font-bold text-text-main">1. Select an execution target</h2>
            </div>
            <span className="font-mono text-[10px] text-text-subtle">PUBLIC GRAPH</span>
          </div>
          <label className="mt-5 block text-[11px] font-semibold text-text-muted" htmlFor="execution-project">Tracked project</label>
          <select
            id="execution-project"
            value={selectedProject?.id ?? ""}
            onChange={(event) => void persist({ ...workspace, selected_project_id: event.target.value }, "Execution target updated")}
            className="mt-2 w-full rounded-xl border border-border bg-background px-3 py-3 text-sm text-text-main"
          >
            {projects.map((project) => <option key={project.id} value={project.id}>{project.name} — {project.province}</option>)}
          </select>

          {selectedProject && <div className="mt-5 space-y-4">
            <div className="rounded-xl border border-borderSubtle bg-surface/70 p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-[10px] font-mono uppercase tracking-wider text-aurora">{selectedProject.sector}</p>
                  <h3 className="mt-1 text-base font-bold text-text-main">{selectedProject.name}</h3>
                  <p className="mt-1 text-xs text-text-muted">{selectedProject.location_name}, {selectedProject.province} · {selectedProject.current_stage.replaceAll("_", " ")}</p>
                </div>
                <span className="rounded-lg border border-primary/35 bg-primary/10 px-2 py-1 text-xs font-bold text-aurora">{buildability}/100</span>
              </div>
              <div className="mt-4 grid grid-cols-2 gap-2 text-xs">
                <div className="rounded-lg bg-background/70 p-2.5"><span className="block text-[10px] text-text-subtle">Reported CAPEX</span><strong className="font-mono text-text-main">{cad.format(capex)}</strong></div>
                <div className="rounded-lg bg-background/70 p-2.5"><span className="block text-[10px] text-text-subtle">Evidence records</span><strong className="font-mono text-text-main">{selectedProject.evidence?.length ?? 0}</strong></div>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <button onClick={toggleSavedProject} className="inline-flex items-center justify-center gap-2 rounded-xl border border-primary/40 bg-primary/10 px-3 py-2.5 text-xs font-bold text-aurora transition hover:bg-primary/20">
                <Bookmark className="h-3.5 w-3.5" /> {isSaved ? "Saved" : "Add to list"}
              </button>
              <Link href={`/projects/${selectedProject.slug}`} className="inline-flex items-center justify-center gap-2 rounded-xl border border-border bg-surface px-3 py-2.5 text-xs font-bold text-text-main transition hover:border-aurora">
                Dossier <ArrowUpRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>}
        </section>

        <section className="glass-panel rounded-2xl border border-primary/30 p-5 sm:p-6">
          <div className="flex flex-col gap-3 border-b border-borderSubtle pb-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-2"><CircleDollarSign className="h-4 w-4 text-gold" /><h2 className="text-sm font-bold text-text-main">2. Economic Flight Simulator</h2></div>
            <span className="rounded-full border border-gold/35 bg-gold/10 px-2.5 py-1 font-mono text-[10px] font-bold text-gold">SCENARIO, NOT FORECAST</span>
          </div>
          <p className="mt-3 text-xs leading-relaxed text-text-muted">Test an intervention against the selected project. Outputs are transparent planning estimates based on your assumptions and the recorded project baseline—not a financing commitment or public forecast.</p>
          <div className="mt-6 grid gap-5 md:grid-cols-2">
            <div className="space-y-4">
              <RangeField label="Catalytic public / partner capital" value={scenario.catalytic_capital_cad} min={0} max={500_000_000} step={25_000_000} display={preciseCad.format(scenario.catalytic_capital_cad)} onChange={(value) => updateScenario("catalytic_capital_cad", value)} />
              <RangeField label="Approval or permit months removed" value={scenario.permit_months_saved} min={0} max={36} step={1} display={`${scenario.permit_months_saved} months`} onChange={(value) => updateScenario("permit_months_saved", value)} />
              <RangeField label="Target domestic procurement" value={scenario.domestic_procurement_pct} min={0} max={100} step={5} display={`${scenario.domestic_procurement_pct}%`} onChange={(value) => updateScenario("domestic_procurement_pct", value)} />
              <RangeField label="Target Indigenous equity participation" value={scenario.indigenous_equity_pct} min={0} max={50} step={1} display={`${scenario.indigenous_equity_pct}%`} onChange={(value) => updateScenario("indigenous_equity_pct", value)} />
            </div>
            <div className="grid grid-cols-2 gap-3 self-start">
              <Metric icon={Landmark} label="Illustrative capital mobilized" value={cad.format(capitalMobilized)} detail={`${capitalMultiplier.toFixed(2)}× scenario multiplier`} tone="text-aurora" />
              <Metric icon={Clock3} label="Potential FID acceleration" value={`${scenario.permit_months_saved} mo`} detail="Assumption-led" tone="text-gold" />
              <Metric icon={Globe2} label="Domestic procurement target" value={cad.format(domesticSpend)} detail={`${scenario.domestic_procurement_pct}% of reported CAPEX`} tone="text-sky-300" />
              <Metric icon={UsersRound} label="Indigenous equity target" value={cad.format(indigenousEquity)} detail={`${scenario.indigenous_equity_pct}% of reported CAPEX`} tone="text-amber-300" />
              <Metric icon={Network} label="Construction person-years" value={personYears.toLocaleString("en-CA")} detail="Sector model estimate" tone="text-text-main" />
              <Metric icon={BadgeCheck} label="Evidence available" value={`${selectedProject?.evidence?.length ?? 0} records`} detail={`Updated ${selectedProject ? new Date(selectedProject.last_meaningful_update).toLocaleDateString("en-CA") : "—"}`} tone="text-text-main" />
            </div>
          </div>
          <div className="mt-5 flex flex-wrap items-center gap-3 border-t border-borderSubtle pt-5">
            <button onClick={saveScenario} disabled={saving} className="inline-flex items-center gap-2 rounded-xl bg-primary px-4 py-2.5 text-xs font-black text-[#04110b] transition hover:bg-aurora-mint disabled:opacity-60"><Save className="h-3.5 w-3.5" />{saving ? "Saving…" : "Save scenario"}</button>
            {selectedProject && <a href={`/api/v1/cegs/export?project_id=${encodeURIComponent(selectedProject.id)}`} target="_blank" rel="noreferrer" className="inline-flex items-center gap-2 rounded-xl border border-border bg-surface px-4 py-2.5 text-xs font-bold text-text-main transition hover:border-primary"><FileJson2 className="h-3.5 w-3.5 text-aurora" />Open CEGS export</a>}
            {notice && <span role="status" className="text-xs text-text-muted">{notice}</span>}
          </div>
        </section>
      </div>

      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <section className="glass-panel rounded-2xl border border-border/80 p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3 border-b border-borderSubtle pb-4"><div className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-aurora" /><h2 className="text-sm font-bold text-text-main">3. Accountable decision log</h2></div><span className="font-mono text-[10px] text-text-subtle">RATIONALE REQUIRED</span></div>
          <div className="mt-5 grid gap-3 sm:grid-cols-[180px_1fr_auto]">
            <select value={decisionStatus} onChange={(event) => setDecisionStatus(event.target.value as DecisionStatus)} className="rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main"><option value="ADVANCE">Advance</option><option value="INVESTIGATE">Investigate</option><option value="MONITOR">Monitor</option></select>
            <input value={rationale} maxLength={280} onChange={(event) => setRationale(event.target.value)} placeholder="Decision rationale and next accountable action" className="rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main placeholder:text-text-subtle" />
            <button onClick={recordDecision} className="inline-flex items-center justify-center gap-2 rounded-xl border border-primary/45 bg-primary/10 px-4 py-2.5 text-xs font-bold text-aurora hover:bg-primary/20"><Plus className="h-3.5 w-3.5" />Record</button>
          </div>
          <div className="mt-5 space-y-2">
            {workspace.decisions.length ? workspace.decisions.map((decision) => {
              const project = projects.find((candidate) => candidate.id === decision.project_id);
              return <div key={decision.id} className="flex flex-col gap-3 rounded-xl border border-borderSubtle bg-surface/55 p-3 text-xs sm:flex-row sm:items-center"><span className={`w-fit rounded-full border px-2 py-1 font-mono text-[10px] font-bold ${statusTheme(decision.status)}`}>{decision.status}</span><div className="min-w-0 flex-1"><p className="font-bold text-text-main">{project?.name ?? decision.project_id}</p><p className="mt-0.5 text-text-muted">{decision.rationale}</p></div><time className="font-mono text-[10px] text-text-subtle">{new Date(decision.created_at).toLocaleDateString("en-CA")}</time></div>;
            }) : <p className="rounded-xl border border-dashed border-borderSubtle p-4 text-xs text-text-muted">No decision records yet. A good demo has a clear owner, rationale, and next action—not a mystery score.</p>}
          </div>
        </section>

        <aside className="rounded-2xl border border-amber-400/25 bg-gradient-to-br from-[#1a160d] to-card p-5 sm:p-6">
          <div className="flex items-center gap-2"><LockKeyhole className="h-4 w-4 text-gold" /><h2 className="text-sm font-bold text-text-main">Deployment security posture</h2></div>
          <div className="mt-4 space-y-3 text-xs leading-relaxed text-text-muted">
            <p className="rounded-lg border border-borderSubtle bg-background/40 p-3"><strong className="text-text-main">This workspace uses an encrypted, HttpOnly, same-site session cookie.</strong> It is intentionally limited to scenario and decision metadata; never store deal documents, personal data, or credentials here.</p>
            <p className="rounded-lg border border-borderSubtle bg-background/40 p-3">{mode === "CONFIGURED" ? "A deployment-specific workspace secret is configured." : "Demo mode is active. Set COG_WORKSPACE_SECRET (32+ characters) before production use."}</p>
            <p className="text-[11px]">Production multi-party rooms require managed identity (OIDC/SAML + SCIM), server-side tenancy, customer-managed retention, ABAC enforcement, and independent security validation.</p>
          </div>
          <Link href="/operations" className="mt-5 inline-flex items-center gap-1.5 text-xs font-bold text-gold hover:text-gold-light">Open operational control plane <ChevronRight className="h-3.5 w-3.5" /></Link>
        </aside>
      </div>
      {!ready && <div className="sr-only" aria-live="polite">Loading workspace</div>}
    </div>
  );
}

function RangeField({ label, value, min, max, step, display, onChange }: { label: string; value: number; min: number; max: number; step: number; display: string; onChange: (value: number) => void }) {
  return <label className="block rounded-xl border border-borderSubtle bg-surface/45 p-3"><span className="flex items-center justify-between gap-3 text-xs text-text-muted"><span>{label}</span><strong className="font-mono text-text-main">{display}</strong></span><input className="mt-3 w-full accent-emerald-400" type="range" value={value} min={min} max={max} step={step} onChange={(event) => onChange(Number(event.target.value))} /></label>;
}

function Metric({ icon: Icon, label, value, detail, tone }: { icon: typeof Landmark; label: string; value: string; detail: string; tone: string }) {
  return <div className="min-h-32 rounded-xl border border-borderSubtle bg-background/55 p-3"><Icon className={`h-4 w-4 ${tone}`} /><p className="mt-3 text-[10px] leading-tight text-text-subtle">{label}</p><p className={`mt-1 font-mono text-base font-black ${tone}`}>{value}</p><p className="mt-1 text-[10px] text-text-muted">{detail}</p></div>;
}
