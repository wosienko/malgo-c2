import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { uuidSchema } from '$lib/validationSchemas';
import { getAllSessionsForProject } from '$lib/services/c2-sessions-service';

export const GET: RequestHandler = async ({ params }) => {
	const { id } = params;
	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	return json(await getAllSessionsForProject(projectId.data));
};
