import type { PageLoad } from './$types';
import type { Projects } from '$lib';

export const load: PageLoad = async ({ fetch, url }) => {
	const legacy = url.searchParams.get('legacy');
	const projects: Promise<Projects> = fetch(
		`/api/projects?page=1&pageSize=999&legacy=${legacy}`
	).then((r) => r.json());

	return {
		projects,
		count: fetch(`/api/projects/count`).then((r) => r.json())
	};
};
