"use client";

import { useState } from "react";
import NationalPlanningWorkbench from "@/components/NationalPlanningWorkbench";
import DecisionPlanningWorkbench from "@/components/DecisionPlanningWorkbench";
import type { Project } from "@/lib/types";
import { SlidersHorizontal, ShieldAlert } from "lucide-react";

export default function PlanningPageClient({ projects }: { projects: Project[] }) {
  const [view, setView] = useState<"national" | "classic">("national");

  return (
    <div className="min-h-screen bg-slate-950">
      <div className="bg-slate-900/80 border-b border-slate-800 px-4 md:px-8 py-2.5 flex items-center justify-between sticky top-0 z-40 backdrop-blur">
        <div className="flex items-center gap-3">
          <div className="flex bg-slate-950 rounded-lg p-1 border border-slate-800">
            <button
              onClick={() => setView("national")}
              className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-bold transition ${
                view === "national"
                  ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              <ShieldAlert className="w-3.5 h-3.5 text-cyan-400" />
              Sovereign War Game & Optimizer (v2.0)
            </button>
            <button
              onClick={() => setView("classic")}
              className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-bold transition ${
                view === "classic"
                  ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              <SlidersHorizontal className="w-3.5 h-3.5 text-slate-400" />
              Classic Delivery Workbench (v1.0)
            </button>
          </div>
        </div>
        <div className="hidden sm:flex items-center gap-2 text-[11px] text-slate-400 font-mono">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          CEGS-SEC-V2.0 • AUDIT-GRADE PROVENANCE
        </div>
      </div>

      {view === "national" ? (
        <NationalPlanningWorkbench projects={projects} />
      ) : (
        <DecisionPlanningWorkbench projects={projects} />
      )}
    </div>
  );
}
