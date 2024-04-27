import type { RequestHandler } from './$types';
import { isUserOperator } from '$lib/services/roles-service';
import { json } from '@sveltejs/kit';
import { uuidSchema } from '$lib/validationSchemas';
import { isUserAllowedToQueryProject } from '$lib/services/project-service';
import { getLatestCommandForSession } from '$lib/services/c2-sessions-service';

export const GET: RequestHandler = async ({ params, locals }) => {
	const userId = locals.user!.id;
	if (!(await isUserOperator(userId))) {
		return json(
			{
				message: 'You are not authorized'
			},
			{
				status: 403
			}
		);
	}

	const { id, sessionid } = params;
	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	const isUserAuthorized = await isUserAllowedToQueryProject(userId, projectId.data);
	if (!isUserAuthorized) {
		return json(
			{
				message: 'You are not authorized to query this project'
			},
			{
				status: 403
			}
		);
	}

	return json(await getLatestCommandForSession(sessionId.data));
};
