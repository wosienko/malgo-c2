import type { PageLoad } from './$types';
import type { Projects, UsersWithRoles, UserWithRoles } from '$lib';

export const load: PageLoad = async ({ fetch, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 8;
	const projects: Promise<Projects> = fetch(`/api/project?page=${page}&pageSize=${pageSize}`).then(
		(r) => r.json()
	);
	const users: Promise<UsersWithRoles> = fetch(`/api/user?page=1&pageSize=9999`)
		.then((r) => r.json())
		.then((users) => {
			return users.users.filter((user: UserWithRoles) => user.operator);
		});

	return {
		projects,
		users,
		page,
		pageSize
	};
};
