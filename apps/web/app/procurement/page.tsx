import { getProcurements } from "@/lib/data";
import ProcurementExplorer from "@/components/ProcurementExplorer";
import { Activity, ShieldCheck } from "lucide-react";

export default async function ProcurementPage() {
  const procurements = await getProcurements();

  return (
    <div className="mx-auto max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="border-b border-border/80 pb-6 flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <div className="mb-3 inline-flex items-center gap-1.5 rounded-full border border-primary/40 bg-card px-3 py-1 font-mono text-[11px] text-aurora shadow-sm">
            <Activity className="h-3.5 w-3.5 text-aurora animate-pulse" />
            CANADIAN PUBLIC PROCUREMENT & DEFENCE RADAR
          </div>
          <h1 className="text-2xl font-black tracking-tight text-text-main sm:text-4xl">
            Canadian <span className="text-aurora">Procurement Radar</span>
          </h1>
          <p className="mt-2 max-w-3xl text-sm leading-relaxed text-text-muted">
            Tracking live federal, crown, and defence tender opportunities across CanadaBuys, Defence Construction Canada (DCC), and major capital project procurement pipelines.
          </p>
        </div>

        <div className="flex items-center gap-2 text-xs font-mono text-text-subtle">
          <span className="px-2.5 py-1 rounded-lg border border-borderSubtle bg-surface flex items-center gap-1.5 text-aurora">
            <ShieldCheck className="h-3.5 w-3.5" /> Tier 1 CanadaBuys Verified
          </span>
        </div>
      </header>

      <ProcurementExplorer initialProcurements={procurements} />
    </div>
  );
}

