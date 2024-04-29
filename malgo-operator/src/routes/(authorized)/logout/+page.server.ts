import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { invalidateSession } from '$lib/services/session-service';

export const load: PageServerLoad = async () => {
	// kept to enable the server route
};

export const actions: Actions = {
	default: async ({ locals, cookies }) => {
		await invalidateSession(locals.session!.id, cookies);
		redirect(302, '/login');
	}
};
