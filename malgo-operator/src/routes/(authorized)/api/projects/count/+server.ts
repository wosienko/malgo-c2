import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getCountOfProjectsForOperator } from '$lib/services/project-service';

export const GET: RequestHandler = async ({ locals }) => {
	const userId = locals.user!.id;

	return json(await getCountOfProjectsForOperator(userId));
};
