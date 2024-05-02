import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return {
		user: await fetch('/api/self').then((r) => r.json())
	};
};
