import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getCountOfUsers, getUsersWithRoles } from '$lib/services/user-service';

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
		users: await getUsersWithRoles(page, pageSize),
		count: await getCountOfUsers()
	});
};
