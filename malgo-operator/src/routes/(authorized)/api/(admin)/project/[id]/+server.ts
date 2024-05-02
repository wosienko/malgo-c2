import type { RequestHandler } from './$types';
import { uuidSchema, type ProjectSchema, projectSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import { deleteProject, updateProjectData } from '$lib/services/project-service';

export const PATCH: RequestHandler = async ({ params, request }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const body: ProjectSchema = await request.json();

	const result = projectSchema.safeParse(body);
	if (!result.success) {
		return json(result.error.issues, { status: 400 });
	}

	const errorMessage = await updateProjectData(
		projectId.data,
		result.data.name,
		result.data.startDate,
		result.data.endDate,
		result.data.description
	);
	if (errorMessage) {
		return json({ message: errorMessage }, { status: 400 });
	}

	return json({ id, ...result.data });
};

export const DELETE: RequestHandler = async ({ params }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const errorMessage = await deleteProject(projectId.data);
	if (errorMessage) {
		return json({ message: errorMessage }, { status: 400 });
	}

	return json({ id });
};
