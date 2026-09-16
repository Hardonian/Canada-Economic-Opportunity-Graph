import { getProjectBySlug } from "@/lib/data";
import { fetchProjectMRIO } from "@/lib/intelligence";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET(_request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const project = await getProjectBySlug(id);
  if (!project) return publicJSON({ error: "Project not found" }, { status: 404 });

  const mrio = await fetchProjectMRIO(project.id, project);
  return publicJSON({ mrio, project_id: project.id, project_name: project.name });
}

export const OPTIONS = publicOptions;
