import type { LayoutServerLoad } from './$types';
import { dev } from '$app/environment';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = () => {
	if (!dev) throw error(404);
};
