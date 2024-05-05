import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getAllProjectsForOperator, getProjectsForOperator } from '$lib/services/project-service';

export const GET: RequestHandler = async ({ url, locals }) => {
	const userId = locals.user!.id;

	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 10;

	const all = url.searchParams.get('all') === 'true';

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
	if (all) {
		return json(await getAllProjectsForOperator(userId, page, pageSize));
	}
	return json(await getProjectsForOperator(userId, page, pageSize));
};
