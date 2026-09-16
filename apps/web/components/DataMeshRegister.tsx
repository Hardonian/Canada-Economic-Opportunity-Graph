import sourceSnapshot from "@/data/sources.snapshot.json";

interface SnapshotSource {
  id?: string;
  name?: string;
  publisher_name?: string;
  jurisdiction?: string;
  authority_tier?: number;
  source_family?: string;
  access_method?: string;
  languages?: string[];
  update_frequency?: string;
  lifecycle?: string;
  health?: string;
  integration_status?: string;
  evidence_record_count?: number;
  coverage_class?: string;
}

const SOURCES = sourceSnapshot as SnapshotSource[];

const ACTIVE = SOURCES.filter((s) => s.lifecycle === "ACTIVE");
const REGISTERED = SOURCES.filter((s) => s.lifecycle !== "ACTIVE");

function familyLabel(family?: string): string {
  if (!family) return "—";
  return family.replace(/_/g, " ").toLowerCase();
}

function languageLabel(languages?: string[]): string {
  if (!languages || languages.length === 0) return "—";
  if (languages.some((l) => l.startsWith("fr"))) return "EN / FR";
  return "EN";
}

/**
 * A formal, government-register style census of every canonical source in the
 * graph: which are linked to evidence, which are registered but not yet
 * ingested, and the authority tier of each. Deliberately table-shaped rather
 * than card-shaped so the page reads as an official data register.
 */
export default function DataMeshRegister() {
  return (
    <section aria-labelledby="data-register-title" className="space-y-4">
      <div className="flex flex-col gap-1 border-b-2 border-border pb-2 sm:flex-row sm:items-end sm:justify-between">
        <h2
          id="data-register-title"
          className="font-mono text-sm font-bold uppercase tracking-wider text-text-main"
        >
          Canonical source register
        </h2>
        <p className="font-mono text-[10px] text-text-subtle">
          {ACTIVE.length} evidence-linked · {REGISTERED.length} registered not yet ingested
        </p>
      </div>

      <div className="overflow-x-auto rounded-lg border-2 border-border bg-card">
        <table className="min-w-full border-collapse text-left text-xs">
          <caption className="sr-only">
            Canonical data source register with authority tier and ingestion status
          </caption>
          <thead>
            <tr className="border-b-2 border-border bg-surface">
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Publisher / Source
              </th>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Family
              </th>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Tier
              </th>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Lang
              </th>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Cadence
              </th>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wider text-text-subtle">
                Status
              </th>
            </tr>
          </thead>
          <tbody>
            {SOURCES.map((source, index) => {
              const linked = source.integration_status === "EVIDENCE_LINKED";
              return (
                <tr
                  key={source.id ?? `${source.publisher_name}-${index}`}
                  className="border-b border-borderSubtle/60 last:border-b-0 hover:bg-cardHover/40"
                >
                  <th scope="row" className="px-3 py-2 text-left align-top font-normal">
                    <span className="block font-bold text-text-main">
                      {source.name ?? "—"}
                    </span>
                    <span className="block text-[10px] text-text-subtle">
                      {source.publisher_name ?? "—"}
                    </span>
                  </th>
                  <td className="px-3 py-2 align-top font-mono text-[10px] text-text-muted">
                    {familyLabel(source.source_family)}
                  </td>
                  <td className="px-3 py-2 align-top font-mono text-[10px] text-text-muted">
                    {source.authority_tier ?? "—"}
                  </td>
                  <td className="px-3 py-2 align-top font-mono text-[10px] text-text-muted">
                    {languageLabel(source.languages)}
                  </td>
                  <td className="px-3 py-2 align-top font-mono text-[10px] text-text-muted">
                    {(source.update_frequency ?? "—").replace(/_/g, " ").toLowerCase()}
                  </td>
                  <td className="px-3 py-2 align-top">
                    <span
                      className={
                        linked
                          ? "inline-flex items-center gap-1 rounded-full border border-primary/50 bg-primary/10 px-2 py-0.5 font-mono text-[10px] font-bold text-aurora"
                          : "inline-flex items-center gap-1 rounded-full border border-border bg-surface px-2 py-0.5 font-mono text-[10px] font-bold text-text-subtle"
                      }
                    >
                      {linked ? "EVIDENCE-LINKED" : "REGISTERED"}
                    </span>
                    {linked && (source.evidence_record_count ?? 0) > 0 && (
                      <span className="ml-2 font-mono text-[10px] text-text-subtle">
                        {source.evidence_record_count} rec
                      </span>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </section>
  );
}
