import type { LayoutServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = ({ locals }) => {
	if (!locals.user || !locals.session) error(401, 'Unauthorized');
};
