import { lucia } from '$lib/auth.server';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	// ...
};

export const actions: Actions = {
	default: async ({ locals, cookies }) => {
		await lucia.invalidateSession(locals.session!.id);
		const sessionCookie = lucia.createBlankSessionCookie();
		cookies.set(sessionCookie.name, sessionCookie.value, {
			path: '.',
			...sessionCookie.attributes
		});
		redirect(302, '/login');
	}
};
