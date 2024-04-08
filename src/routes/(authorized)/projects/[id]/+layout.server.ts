import type { LayoutServerLoad } from './$types';
import { uuidSchema } from '$lib/validationSchemas';
import { error, json } from '@sveltejs/kit';
import { getProjectById } from '$lib/services/project-service';

export const load: LayoutServerLoad = async ({ params }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const project = await getProjectById(projectId.data);

	if (!project) {
		return error(404, 'Project not found');
	}

	return {
		project
	};
};
