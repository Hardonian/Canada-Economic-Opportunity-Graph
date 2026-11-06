import type { Metadata } from "next";
import ExecutionSuite from "@/components/ExecutionSuite";
import { getProjects } from "@/lib/data";

export const metadata: Metadata = {
  title: "Economic Execution Suite",
  description: "A project decision workspace for evidence, scenarios, and economically accountable execution.",
};

export default async function ExecutionPage() {
  const projects = await getProjects();
  return <ExecutionSuite projects={projects} />;
}
