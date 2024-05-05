import { error, type Handle, type HandleServerError, redirect } from '@sveltejs/kit';
import { createDefaultAdminAndRoles } from '$lib/services/user-service';
import { deleteExpiredSessions, validateSession } from '$lib/services/session-service';
import { isUserAdmin, isUserOperator } from '$lib/services/roles-service';
import { isUserAllowedToQueryProject } from '$lib/services/project-service';

export const handle: Handle = async ({ event, resolve }) => {
	await deleteExpiredSessions();
	await validateSession(event);

	if (event.route.id?.startsWith('/(authorized)/')) {
		const nextRoute = event.route.id.replace('/(authorized)/', '');
		if (!event.locals.session && !event.locals.user) return redirect(303, '/login');

		if (nextRoute.startsWith('admin') && !(await isUserAdmin(event.locals.user!.id)))
			return redirect(303, '/home');

		// API routes
		if (nextRoute.includes('/(admin)') && !(await isUserAdmin(event.locals.user!.id)))
			return error(403, 'Forbidden');

		if (
			(nextRoute.startsWith('projects') || nextRoute.startsWith('api/projects')) &&
			!(await isUserOperator(event.locals.user!.id))
		) {
			return error(403, 'Forbidden');
		}

		if (nextRoute.includes('projects/[id]')) {
			const projectId = event.params.id;
			if (!projectId) return error(400, 'Bad Request');
			if (!(await isUserAllowedToQueryProject(event.locals.user!.id, projectId))) {
				return error(403, 'Forbidden');
			}
		}
	}

	return resolve(event);
};

export const handleError: HandleServerError = async ({ error, status }) => {
	if (status === 401 || status === 403 || status === 404) {
		return;
	}
	console.error(error);
};

// on startup code - insert default roles and admin user
await createDefaultAdminAndRoles();
