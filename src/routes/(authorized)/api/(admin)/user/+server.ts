import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getCountOfUsers, getUsersWithRoles, registerNewUser } from '$lib/services/user-service';
import { adminRegisterSchema } from '$lib/validationSchemas';
import { setAdmin, setOperator } from '$lib/services/roles-service';

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

export const POST: RequestHandler = async ({ request }) => {
	const body = await request.json();

	const validation = adminRegisterSchema.safeParse(body);
	if (!validation.success) {
		return json(validation.error.issues, { status: 400 });
	}

	const newUserId = await registerNewUser(validation.data);
	if (!newUserId) {
		return json(
			{
				message: 'User with this email already exists'
			},
			{
				status: 400
			}
		);
	}

	await setAdmin(newUserId, validation.data.admin);
	await setOperator(newUserId, validation.data.operator);

	return json({ id: newUserId });
};
