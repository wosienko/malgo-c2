import type { LayoutServerLoad } from './$types';
import { isUserAdmin, isUserOperator } from '$lib/services/roles-service';
import { getNameForId } from '$lib/services/user-service';

export const load: LayoutServerLoad = async ({ locals }) => {
	const userId = locals.user?.id ?? '00000000-0000-0000-0000-000000000000';
	return {
		loggedIn: !!locals.user && !!locals.session,
		name: await getNameForId(userId),
		isAdmin: await isUserAdmin(userId),
		isOperator: await isUserOperator(userId)
	};
};
