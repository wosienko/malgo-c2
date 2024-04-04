import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { createProject, getCountOfProjects, getProjects } from '$lib/services/project-service';
import { projectSchema } from '$lib/validationSchemas';
import { isUserAdmin } from '$lib/services/roles-service';

export const GET: RequestHandler = async ({ url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 10;

	if (page < 1 || pageSize < 1) {
		return json(
			{
				message: 'Invalid page or pageSize'
			},
			{
				status: 400
			}
		);
	}

	return json({
		projects: await getProjects(page, pageSize),
		count: await getCountOfProjects()
	});
};

// Admin only
export const POST: RequestHandler = async ({ request, locals }) => {
	const userId = locals.user!.id;
	if (!(await isUserAdmin(userId))) {
		return json(
			{
				message: 'You are not authorized to create a project'
			},
			{
				status: 403
			}
		);
	}

	const body = await request.json();

	const validation = projectSchema.safeParse(body);
	if (!validation.success) {
		return json(validation.error.issues, { status: 400 });
	}

	const newProjectId = await createProject(
		validation.data.name,
		validation.data.startDate,
		validation.data.endDate
	);
	if (!newProjectId) {
		return json(
			{
				message: 'Project with this name already exists'
			},
			{
				status: 400
			}
		);
	}

	return json({ id: newProjectId });
};
