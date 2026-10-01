"use client";

import Link from "next/link";
import { ArrowLeft, FileText, Plane, Leaf, ExternalLink, Printer } from "lucide-react";
import CarneyBriefingDossier from "@/components/finance/CarneyBriefingDossier";

export default function CarneyBriefingPage() {
  return (
    <div className="mx-auto max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8">
      {/* Navigation Breadcrumb */}
      <div className="flex flex-wrap items-center justify-between gap-4 border-b border-borderSubtle pb-4">
        <div className="flex items-center gap-3">
          <Link
            href="/briefing"
            className="inline-flex items-center gap-1.5 text-xs font-mono text-text-subtle hover:text-gold transition-colors"
          >
            <ArrowLeft className="h-3.5 w-3.5" />
            <span>Back to Decision Briefs</span>
          </Link>
          <span className="text-text-subtle">/</span>
          <span className="text-xs font-mono font-bold text-gold uppercase">Carney Initiative Dossier</span>
        </div>

        <div className="flex items-center gap-2">
          <Link
            href="/briefing/memo"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-borderSubtle bg-surface text-xs font-medium text-text-muted hover:text-text-main transition-colors"
          >
            <FileText className="h-3.5 w-3.5" />
            <span>Cabinet Memo Studio</span>
          </Link>
          <Link
            href="/finance"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-borderSubtle bg-surface text-xs font-medium text-text-muted hover:text-text-main transition-colors"
          >
            <span>Project Finance Lab</span>
          </Link>
        </div>
      </div>

      {/* Hero Title */}
      <div className="space-y-3">
        <div className="inline-flex items-center gap-2 rounded-full border border-gold/40 bg-gold/10 px-3 py-1 text-xs font-mono font-bold text-gold">
          <span>EXECUTIVE STRATEGY BRIEFING</span>
          <span>·</span>
          <span>THE RIGHT HONOURABLE MARK CARNEY</span>
        </div>
        <h1 className="text-3xl sm:text-5xl font-black text-text-main tracking-tight">
          National Airports System & <span className="text-gold">Transition Finance</span>
        </h1>
        <p className="text-sm sm:text-base text-text-muted max-w-4xl leading-relaxed">
          The institutional quantitative briefing environment for Mark Carney’s flagship growth initiatives: Restructuring Transport Canada’s National Airports System (NAS) commercial ground leases into an $18B CAD private pension concession program, and operationalizing the GFANZ / SFAC Canadian Transition Finance Taxonomy.
        </p>
      </div>

      {/* Interactive Dossier Component */}
      <CarneyBriefingDossier />
    </div>
  );
}
