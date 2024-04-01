import type { LayoutServerLoad } from './$types';
import { error } from '@sveltejs/kit';
import { isUserAdmin } from '$lib/services/roles-service';

export const load: LayoutServerLoad = async ({ locals }) => {
	if (!(await isUserAdmin(locals.user!.id))) error(403, 'Forbidden');
};
