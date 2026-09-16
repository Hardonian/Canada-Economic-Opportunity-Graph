import { getProjectBySlug } from "@/lib/data";
import { fetchProjectFlyvbjerg } from "@/lib/intelligence";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET(_request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const project = await getProjectBySlug(id);
  if (!project) return publicJSON({ error: "Project not found" }, { status: 404 });

  const flyvbjerg = await fetchProjectFlyvbjerg(project.id, project);
  return publicJSON({ flyvbjerg, project_id: project.id, project_name: project.name });
}

export const OPTIONS = publicOptions;
