import type { Handle, HandleServerError } from '@sveltejs/kit';
import { createDefaultAdminAndRoles } from '$lib/services/user-service';
import { deleteExpiredSessions, validateSession } from '$lib/services/session-service';

export const handle: Handle = async ({ event, resolve }) => {
	await deleteExpiredSessions();
	await validateSession(event);
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
