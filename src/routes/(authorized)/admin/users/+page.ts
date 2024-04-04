import type { PageLoad } from './$types';
import type { UsersWithRoles } from '$lib';

export const load: PageLoad = async ({ fetch, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 8;
	const response = await fetch(`/api/user?page=${page}&pageSize=${pageSize}`);
	const users: UsersWithRoles = await response.json();

	return {
		users,
		page,
		pageSize
	};
};
