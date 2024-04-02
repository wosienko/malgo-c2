import type { PageLoad } from './$types';
import type { UsersWithRoles } from '$lib/db/schema/users';

export const load: PageLoad = async ({ fetch }) => {
	const page = 1;
	const pageSize = 6;
	const response = await fetch(`/api/user?page=${page}&pageSize=${pageSize}`);
	const users: UsersWithRoles = await response.json();

	return {
		users,
		page,
		pageSize
	};
};
