import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getCountOfProjectsForOperator } from '$lib/services/project-service';
import { isUserOperator } from '$lib/services/roles-service';

export const GET: RequestHandler = async ({ locals }) => {
	const userId = locals.user!.id;
	if (!(await isUserOperator(userId))) {
		return json(
			{
				message: 'You are not authorized to create a project'
			},
			{
				status: 403
			}
		);
	}

	return json(await getCountOfProjectsForOperator(userId));
};
