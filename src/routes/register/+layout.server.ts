import type { LayoutServerLoad } from './$types';
import { PUBLIC_ENVIRONMENT } from '$env/static/public';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = () => {
	if (PUBLIC_ENVIRONMENT !== 'DEV') throw error(404);
};
