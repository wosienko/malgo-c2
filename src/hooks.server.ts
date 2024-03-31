import type { Handle, HandleServerError } from '@sveltejs/kit';
import { lucia } from '$lib/auth.server';

export const handle: Handle = async ({ event, resolve }) => {
	// invalidate expired sessions
	await lucia.deleteExpiredSessions();

	const sessionId = event.cookies.get(lucia.sessionCookieName);
	if (!sessionId) {
		event.locals.user = null;
		event.locals.session = null;
		return resolve(event);
	}

	const { session, user } = await lucia.validateSession(sessionId);
	if (session && session.fresh) {
		const sessionCookie = lucia.createSessionCookie(session.id);
		event.cookies.set(sessionCookie.name, sessionCookie.value, {
			path: '/',
			...sessionCookie.attributes
		});
	}
	if (!session) {
		const sessionCookie = lucia.createBlankSessionCookie();
		event.cookies.set(sessionCookie.name, sessionCookie.value, {
			path: '/',
			...sessionCookie.attributes
		});
	}

	event.locals.user = user;
	event.locals.session = session;
	return resolve(event);
};

export const handleError: HandleServerError = async ({ error, status }) => {
	if (status === 401 || status === 403 || status === 404) {
		return;
	}
	console.error(error);
};

// on startup code - insert default roles and admin user
import { db } from '$lib/db/db.server';
import { Users } from '$lib/db/schema/users';
import { Roles } from '$lib/db/schema/users';
import { UserRoles } from '$lib/db/schema/users';
import { Argon2id } from 'oslo/password';
import { DEFAULT_ADMIN_EMAIL, DEFAULT_ADMIN_PASSWORD } from '$env/static/private';

await db.transaction(async (tx) => {
	const adminRoleId = await tx
		.insert(Roles)
		.values({ name: 'Admin' })
		.returning({ id: Roles.id })
		.onConflictDoNothing();

	const operatorRoleId = await tx
		.insert(Roles)
		.values({ name: 'Operator' })
		.returning({ id: Roles.id })
		.onConflictDoNothing();

	const userRoleId = await tx
		.insert(Roles)
		.values({ name: 'User' })
		.returning({ id: Roles.id })
		.onConflictDoNothing();

	const adminId = await tx
		.insert(Users)
		.values({
			name: 'Admin',
			surname: 'Admin',
			email: DEFAULT_ADMIN_EMAIL,
			password: await new Argon2id().hash(DEFAULT_ADMIN_PASSWORD)
		})
		.returning({ id: Users.id })
		.onConflictDoNothing();

	if (!adminId.length || !adminRoleId.length || !operatorRoleId.length || !userRoleId.length) {
		console.log('Default roles and admin user already exist');
		return;
	}

	await tx
		.insert(UserRoles)
		.values([
			{
				user_id: adminId[0].id,
				role_id: adminRoleId[0].id
			},
			{
				user_id: adminId[0].id,
				role_id: operatorRoleId[0].id
			},
			{
				user_id: adminId[0].id,
				role_id: userRoleId[0].id
			}
		])
		.onConflictDoNothing();
});
