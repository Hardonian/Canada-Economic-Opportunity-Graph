import type { Metadata } from "next";
import Link from "next/link";
import {
  Activity,
  AlertTriangle,
  ArrowUpRight,
  Braces,
  CheckCircle2,
  Clock3,
  Database,
  HardDrive,
  Network,
  RefreshCw,
  Server,
  ShieldCheck,
  Workflow,
} from "lucide-react";
import { getOperationsStatus } from "@/lib/operations";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Operations & Product Control Plane",
  description: "Runtime readiness, ingestion health, data freshness, and product capability status for CanadaOpportunityGraph.",
};

const money = new Intl.NumberFormat("en-CA", {
  style: "currency",
  currency: "CAD",
  notation: "compact",
  maximumFractionDigits: 1,
});

const dateTime = new Intl.DateTimeFormat("en-CA", {
  dateStyle: "medium",
  timeStyle: "short",
  timeZone: "UTC",
});

function statusTone(status: string): string {
  if (["OPERATIONAL", "HEALTHY", "CURRENT", "ACTIVE", "LIVE"].includes(status)) {
    return "border-primary/50 bg-primary/10 text-aurora";
  }
  if (["RESILIENT_SNAPSHOT", "DEGRADED", "CURATED_SNAPSHOT"].includes(status)) {
    return "border-gold/50 bg-gold/10 text-gold";
  }
  return "border-crimson/50 bg-crimson/10 text-crimson-light";
}

export default async function OperationsPage() {
  const operations = await getOperationsStatus();
  const readiness = [
    {
      label: "Serving layer",
      detail: operations.upstream.ready ? `Go API ready in ${operations.upstream.latency_ms} ms` : "Reviewed snapshot remains available",
      passed: operations.upstream.ready || operations.data.freshness === "CURRENT",
    },
    {
      label: "Durable storage",
      detail: operations.runtime.storage_mode === "persistent" ? "WAL-backed persistent store" : "Immutable snapshot profile",
      passed: operations.runtime.storage_mode === "persistent",
    },
    {
      label: "Scheduled ingestion",
      detail: operations.runtime.scheduler_enabled ? `Refresh every ${operations.runtime.scheduler_interval}` : "Disabled in this runtime",
      passed: operations.runtime.scheduler_enabled,
    },
    {
      label: "Connector mesh",
      detail: operations.ingestion.connectors.length
        ? `${operations.ingestion.healthy} healthy · ${operations.ingestion.degraded} degraded · ${operations.ingestion.broken} broken`
        : "Available when the Go API is connected",
      passed: operations.ingestion.connectors.length > 0 && operations.ingestion.broken === 0,
    },
    {
      label: "Release freshness",
      detail: operations.data.age_hours >= 0 ? `${Math.round(operations.data.age_hours / 24)} days since reviewed release` : "Release timestamp invalid",
      passed: operations.data.freshness === "CURRENT",
    },
    {
      label: "Same-origin API gateway",
      detail: "Read and model routes are exposed through /api/v1/*",
      passed: true,
    },
  ];

  return (
    <div className="mx-auto max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="relative overflow-hidden rounded-2xl border border-border/80 bg-gradient-to-br from-card via-surface to-background p-6 shadow-2xl sm:p-8">
        <div aria-hidden="true" className="pointer-events-none absolute -right-20 -top-24 h-72 w-72 rounded-full bg-aurora/10 blur-3xl" />
        <div className="relative z-10 grid gap-6 lg:grid-cols-[1fr_auto] lg:items-end">
          <div>
            <div className={`inline-flex items-center gap-2 rounded-full border px-3 py-1 font-mono text-[11px] font-bold ${statusTone(operations.status)}`}>
              <span className="h-2 w-2 rounded-full bg-current" />
              {operations.status.replaceAll("_", " ")}
            </div>
            <h1 className="mt-4 text-3xl font-black tracking-tight text-text-main sm:text-5xl">
              Product <span className="text-aurora">Control Plane</span>
            </h1>
            <p className="mt-3 max-w-3xl text-sm leading-relaxed text-text-muted">
              One operational view across the serving API, persistent data layer, ingestion connectors, reviewed releases, and product modules. Every status below is derived from a runtime probe or the checksummed bundled snapshot.
            </p>
          </div>
          <div className="rounded-xl border border-border bg-background/80 p-4 font-mono text-[11px] text-text-muted">
            <div className="flex items-center gap-2 text-text-main">
              <Clock3 aria-hidden="true" className="h-4 w-4 text-aurora" /> Last probe
            </div>
            <div className="mt-1">{dateTime.format(new Date(operations.checked_at))} UTC</div>
            <div className="mt-2 text-text-subtle">Mode: {operations.mode.replaceAll("_", " ")}</div>
          </div>
        </div>
      </header>

      <section aria-label="Operational metrics" className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <MetricCard icon={Database} label="Tracked projects" value={operations.data.total_projects.toLocaleString("en-CA")} detail={operations.data.dataset_version} />
        <MetricCard icon={Activity} label="Capital represented" value={money.format(operations.data.total_capex_cad)} detail={`${operations.data.evidence_records.toLocaleString("en-CA")} evidence records`} />
        <MetricCard icon={Network} label="Canonical sources" value={operations.data.canonical_sources.toLocaleString("en-CA")} detail={`${operations.ingestion.connectors.length} runtime connectors`} />
        <MetricCard icon={RefreshCw} label="Data freshness" value={operations.data.freshness} detail={`${Math.max(0, Math.round(operations.data.age_hours / 24))} days old`} tone={operations.data.freshness === "CURRENT" ? "green" : "gold"} />
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <div className="glass-panel rounded-2xl p-5 sm:p-6">
          <div className="flex items-start justify-between gap-4 border-b border-borderSubtle pb-4">
            <div>
              <div className="flex items-center gap-2">
                <Workflow aria-hidden="true" className="h-5 w-5 text-aurora" />
                <h2 className="text-lg font-black text-text-main">Runtime topology</h2>
              </div>
              <p className="mt-1 text-xs text-text-muted">The active request and refresh path for this deployment.</p>
            </div>
            <span className={`rounded border px-2 py-1 font-mono text-[10px] font-bold ${statusTone(operations.upstream.ready ? "HEALTHY" : "RESILIENT_SNAPSHOT")}`}>
              {operations.upstream.ready ? "CONNECTED" : "SNAPSHOT"}
            </span>
          </div>
          <div className="mt-5 grid gap-3 sm:grid-cols-4">
            <TopologyNode icon={Server} title="Next.js web" detail="Same-origin BFF" active />
            <TopologyNode icon={Braces} title="API gateway" detail="/api/v1/*" active />
            <TopologyNode icon={Activity} title="Go engine" detail={operations.upstream.ready ? "Ready" : "Not connected"} active={operations.upstream.ready} />
            <TopologyNode icon={HardDrive} title="Data store" detail={operations.runtime.storage_mode} active={operations.runtime.storage_mode === "persistent"} />
          </div>
          <dl className="mt-5 grid gap-3 rounded-xl border border-borderSubtle bg-background/60 p-4 text-xs sm:grid-cols-2">
            <RuntimeFact label="Environment" value={operations.runtime.environment} />
            <RuntimeFact label="API version" value={operations.upstream.version || "snapshot-only"} />
            <RuntimeFact label="Storage profile" value={operations.runtime.storage_mode} />
            <RuntimeFact label="Refresh cadence" value={operations.runtime.scheduler_enabled ? operations.runtime.scheduler_interval || "enabled" : "disabled"} />
          </dl>
        </div>

        <div className="glass-panel rounded-2xl p-5 sm:p-6">
          <div className="flex items-center gap-2">
            <ShieldCheck aria-hidden="true" className="h-5 w-5 text-gold" />
            <h2 className="text-lg font-black text-text-main">Production readiness</h2>
          </div>
          <div className="mt-4 space-y-2.5">
            {readiness.map((item) => (
              <div key={item.label} className="flex items-start gap-3 rounded-lg border border-borderSubtle bg-background/50 p-3">
                {item.passed
                  ? <CheckCircle2 aria-hidden="true" className="mt-0.5 h-4 w-4 shrink-0 text-aurora" />
                  : <AlertTriangle aria-hidden="true" className="mt-0.5 h-4 w-4 shrink-0 text-gold" />}
                <div>
                  <div className="text-xs font-bold text-text-main">{item.label}</div>
                  <div className="mt-0.5 text-[11px] text-text-muted">{item.detail}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="glass-panel rounded-2xl p-5 sm:p-6">
        <div className="flex flex-col justify-between gap-3 border-b border-borderSubtle pb-4 sm:flex-row sm:items-end">
          <div>
            <div className="flex items-center gap-2">
              <Network aria-hidden="true" className="h-5 w-5 text-aurora" />
              <h2 className="text-lg font-black text-text-main">Ingestion connector mesh</h2>
            </div>
            <p className="mt-1 text-xs text-text-muted">Health, operating mode, and observed document volume from the exact adapters serving the API.</p>
          </div>
          <Link href="/api/v1/integrations/status" className="inline-flex items-center gap-1 text-xs font-semibold text-aurora hover:underline">
            Integration JSON <ArrowUpRight aria-hidden="true" className="h-3.5 w-3.5" />
          </Link>
        </div>
        {operations.ingestion.connectors.length ? (
          <div className="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
            {operations.ingestion.connectors.map((connector) => (
              <article key={connector.name} className="rounded-xl border border-borderSubtle bg-background/55 p-4">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <h3 className="truncate text-xs font-bold text-text-main">{connector.name.replaceAll("_", " ")}</h3>
                    <p className="mt-1 font-mono text-[10px] text-text-subtle">Tier {connector.tier} · {connector.mode.replaceAll("_", " ")}</p>
                  </div>
                  <span className={`rounded border px-2 py-0.5 font-mono text-[9px] font-bold ${statusTone(connector.status)}`}>{connector.status}</span>
                </div>
                <div className="mt-3 grid grid-cols-3 gap-2 border-t border-borderSubtle pt-3 text-center">
                  <SmallStat label="Seen" value={connector.documents_seen} />
                  <SmallStat label="Changed" value={connector.documents_changed} />
                  <SmallStat label="Failures" value={connector.parse_failures} />
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="mt-4 rounded-xl border border-gold/40 bg-gold/10 p-4 text-sm text-text-muted">
            Connector telemetry is unavailable because this web deployment is serving its reviewed snapshot without a configured Go API. Product pages remain usable; live model execution and event streaming require <code className="text-text-main">COG_API_BASE</code>.
          </div>
        )}
      </section>

      <section>
        <div className="mb-4 flex flex-col justify-between gap-2 sm:flex-row sm:items-end">
          <div>
            <div className="font-mono text-[10px] font-bold uppercase tracking-wider text-aurora">Connected product suite</div>
            <h2 className="mt-1 text-2xl font-black text-text-main">Launch every operational capability</h2>
          </div>
          <Link href="/api/v1" className="inline-flex items-center gap-1 text-xs font-semibold text-text-muted hover:text-aurora">
            API index <ArrowUpRight aria-hidden="true" className="h-3.5 w-3.5" />
          </Link>
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {operations.capabilities.map((capability) => (
            <Link key={capability.id} href={capability.href} className="glass-card group rounded-xl p-5">
              <div className="flex items-start justify-between gap-4">
                <Braces aria-hidden="true" className="h-5 w-5 text-gold" />
                <span className={`rounded border px-2 py-0.5 font-mono text-[9px] font-bold ${statusTone(capability.status)}`}>{capability.status}</span>
              </div>
              <h3 className="mt-5 text-sm font-black text-text-main group-hover:text-aurora">{capability.label}</h3>
              <div className="mt-2 flex items-center justify-between gap-3 font-mono text-[10px] text-text-subtle">
                <span className="truncate">{capability.endpoint}</span>
                <ArrowUpRight aria-hidden="true" className="h-3.5 w-3.5 shrink-0 text-aurora" />
              </div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}

function MetricCard({ icon: Icon, label, value, detail, tone = "green" }: { icon: typeof Database; label: string; value: string; detail: string; tone?: "green" | "gold" }) {
  return (
    <div className="glass-card rounded-xl p-4">
      <div className="flex items-center justify-between text-[11px] font-semibold uppercase tracking-wide text-text-subtle">
        <span>{label}</span>
        <Icon aria-hidden="true" className={`h-4 w-4 ${tone === "green" ? "text-aurora" : "text-gold"}`} />
      </div>
      <div className={`mt-2 text-2xl font-black font-tabular ${tone === "green" ? "text-text-main" : "text-gold"}`}>{value}</div>
      <div className="mt-1 text-[10px] text-text-muted">{detail}</div>
    </div>
  );
}

function TopologyNode({ icon: Icon, title, detail, active }: { icon: typeof Server; title: string; detail: string; active: boolean }) {
  return (
    <div className={`relative rounded-xl border p-3 text-center ${active ? "border-primary/40 bg-primary/10" : "border-gold/40 bg-gold/10"}`}>
      <Icon aria-hidden="true" className={`mx-auto h-5 w-5 ${active ? "text-aurora" : "text-gold"}`} />
      <div className="mt-2 text-xs font-bold text-text-main">{title}</div>
      <div className="mt-0.5 truncate font-mono text-[9px] uppercase text-text-subtle">{detail}</div>
    </div>
  );
}

function RuntimeFact({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <dt className="text-text-subtle">{label}</dt>
      <dd className="font-mono font-bold text-text-main">{value}</dd>
    </div>
  );
}

function SmallStat({ label, value }: { label: string; value: number }) {
  return (
    <div>
      <div className="font-tabular text-sm font-black text-text-main">{value.toLocaleString("en-CA")}</div>
      <div className="mt-0.5 text-[9px] uppercase text-text-subtle">{label}</div>
    </div>
  );
}
