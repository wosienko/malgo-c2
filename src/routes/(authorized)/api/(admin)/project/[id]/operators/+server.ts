import type { RequestHandler } from './$types';
import { userProjectSchema, uuidSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import { assignOperatorsToProject, getOperatorsForProject } from '$lib/services/project-service';

export const GET: RequestHandler = async ({ params }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	return json(await getOperatorsForProject(projectId.data));
};

export const PUT: RequestHandler = async ({ params, request }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const body = await request.json();
	const operatorIds = userProjectSchema.safeParse(body);
	if (!operatorIds.success) {
		return json(operatorIds.error.issues, { status: 400 });
	}

	const errorMessage = await assignOperatorsToProject(projectId.data, operatorIds.data.users);
	if (errorMessage) {
		return json({ message: errorMessage }, { status: 400 });
	}

	return json({ id, operatorIds });
};
