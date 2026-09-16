import type { Metadata } from "next";
import PlanningPageClient from "@/components/PlanningPageClient";
import { getProjects } from "@/lib/data";

export const metadata: Metadata = {
  title: "National Planning & Sovereign War Game | CanadaOpportunityGraph",
  description:
    "Multi-variable sovereign capital allocation optimization, macroeconomic geopolitical stress-testing, Bayesian megaproject overrun hazard modeling, and craft labor collision radar across Canadian strategic projects.",
};

export default async function PlanningPage() {
  const projects = await getProjects();
  return <PlanningPageClient projects={projects} />;
}
