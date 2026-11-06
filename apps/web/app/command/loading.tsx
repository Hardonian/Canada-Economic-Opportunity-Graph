export default function CommandLoading() {
  return (
    <div className="mx-auto max-w-7xl space-y-7 px-4 py-8 sm:px-6 lg:px-8" aria-live="polite" aria-busy="true">
      <div className="h-80 animate-pulse rounded-3xl border border-border/80 bg-card" />
      <div className="grid gap-5 lg:grid-cols-[0.8fr_1.2fr]">
        <div className="h-64 animate-pulse rounded-2xl border border-border/80 bg-card" />
        <div className="h-64 animate-pulse rounded-2xl border border-border/80 bg-card" />
      </div>
      <div className="h-[560px] animate-pulse rounded-3xl border border-border/80 bg-card" />
      <span className="sr-only">Loading mission control data</span>
    </div>
  );
}
