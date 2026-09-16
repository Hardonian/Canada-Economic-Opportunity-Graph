"use client";

import dynamic from "next/dynamic";
import { Radio, Loader2 } from "lucide-react";
import type { Project } from "@/lib/types";

const DynamicGeospatialMap = dynamic(
  () => import("@/components/GeospatialMap"),
  {
    ssr: false,
    loading: () => (
      <div className="w-full h-[680px] rounded-2xl border border-border/80 bg-[#040806] flex flex-col items-center justify-center p-6 shadow-2xl relative overflow-hidden">
        <div className="absolute inset-0 bg-[radial-gradient(#00F5A0_1px,transparent_1px)] [background-size:24px_24px] opacity-10 pointer-events-none"></div>
        <div className="flex flex-col items-center gap-3 z-10">
          <div className="p-3 rounded-full bg-surface border border-primary/40 text-aurora animate-pulse shadow-[0_0_15px_rgba(0,245,160,0.2)]">
            <Radio className="h-6 w-6" />
          </div>
          <div className="text-sm font-mono font-bold text-text-main flex items-center gap-2">
            <Loader2 className="h-4 w-4 animate-spin text-aurora" />
            INITIALIZING HIGH-RESOLUTION SATELLITE ENGINE
          </div>
          <div className="text-xs font-mono text-text-subtle text-center max-w-sm">
            Calibrating global trade routes, conflict markers, opportunity zones, and multi-spectral raster layers...
          </div>
        </div>
      </div>
    ),
  }
);

interface GeospatialMapWrapperProps {
  projects: Project[];
  activeProject: Project;
  onSelectProject: (p: Project) => void;
  selectedSector: string;
  selectedStage: string;
}

export default function GeospatialMapWrapper(props: GeospatialMapWrapperProps) {
  return <DynamicGeospatialMap {...props} />;
}
