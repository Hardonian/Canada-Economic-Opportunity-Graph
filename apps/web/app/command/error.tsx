"use client";

import { AlertTriangle, RefreshCw } from "lucide-react";

export default function CommandError({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  return (
    <div className="mx-auto max-w-3xl px-4 py-16 sm:px-6">
      <section className="rounded-3xl border border-gold/40 bg-gold/10 p-7 text-center shadow-2xl shadow-black/20">
        <AlertTriangle className="mx-auto h-8 w-8 text-gold" />
        <h1 className="mt-4 text-2xl font-black text-text-main">Mission Control could not load</h1>
        <p className="mx-auto mt-3 max-w-xl text-sm leading-relaxed text-text-muted">The command view is intentionally fail-soft: public intelligence should not be replaced with invented values when a service is unavailable. Try again or use the reviewed snapshot surfaces.</p>
        <button type="button" onClick={reset} className="mt-6 inline-flex min-h-11 items-center gap-2 rounded-xl bg-gold px-4 text-sm font-black text-[#251b03] hover:bg-[#ffe09a]">
          <RefreshCw className="h-4 w-4" /> Try again
        </button>
      </section>
    </div>
  );
}
