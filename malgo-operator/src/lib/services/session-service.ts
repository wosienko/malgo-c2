import type { Cookies, RequestEvent } from '@sveltejs/kit';
import { lucia } from '$lib/auth.server';

export const deleteExpiredSessions = async () => {
	await lucia.deleteExpiredSessions();
};

export const validateSession = async (event: RequestEvent) => {
	const sessionId = event.cookies.get(lucia.sessionCookieName);
	if (!sessionId) {
		event.locals.user = null;
		event.locals.session = null;
		return;
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
};

export const invalidateSession = async (sessionId: string, cookies: Cookies) => {
	await lucia.invalidateSession(sessionId);
	createEmptySessionCookie(cookies);
};

export const createSession = async (userId: string, cookies: Cookies) => {
	if (userId === '') {
		createEmptySessionCookie(cookies);
		return;
	}

	const session = await lucia.createSession(userId, {});
	const sessionCookie = lucia.createSessionCookie(session.id);
	cookies.set(sessionCookie.name, sessionCookie.value, {
		path: '/',
		...sessionCookie.attributes
	});
};

const createEmptySessionCookie = (cookies: Cookies) => {
	const sessionCookie = lucia.createBlankSessionCookie();
	cookies.set(sessionCookie.name, sessionCookie.value, {
		path: '.',
		...sessionCookie.attributes
	});
};
