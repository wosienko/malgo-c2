import type { PageLoad } from './$types';
import type { Projects } from '$lib';

export const load: PageLoad = async ({ fetch, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 8;
	const projects: Promise<Projects> = fetch(`/api/projects?page=${page}&pageSize=${pageSize}`).then(
		(r) => r.json()
	);

	return {
		projects,
		page,
		pageSize
	};
};
