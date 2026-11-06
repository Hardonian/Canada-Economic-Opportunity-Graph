import type { Metadata } from "next";
import MissionControl from "@/components/MissionControl";
import { getProjects, getRadarData, getSignals } from "@/lib/data";

export const metadata: Metadata = {
  title: "National Mission Control",
  description: "A transparent Canadian decision loop connecting public evidence, execution workspaces, and secure deal rooms.",
};

export default async function CommandPage() {
  const [radar, projects, signals] = await Promise.all([getRadarData(), getProjects(), getSignals()]);
  return <MissionControl projects={projects} signals={signals} stats={radar.stats} sourceMode={radar.source_mode} generatedAt={radar.generated_at} />;
}
