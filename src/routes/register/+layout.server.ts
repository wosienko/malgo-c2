import type { LayoutServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = () => {
	if (!import.meta.env.PROD) throw error(404);
};
