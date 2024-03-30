import { DrizzlePostgreSQLAdapter } from '@lucia-auth/adapter-drizzle';
import { Sessions } from '$lib/db/schema/sessions';
import { Users } from '$lib/db/schema/users';
import { db } from '$lib/db/db.server';
import { Lucia, TimeSpan } from 'lucia';
import { dev } from '$app/environment';

const adapter = new DrizzlePostgreSQLAdapter(db, Sessions, Users);

export const lucia = new Lucia(adapter, {
	sessionCookie: {
		attributes: {
			secure: !dev
		}
	},
	// expire session after 1 hour
	sessionExpiresIn: new TimeSpan(60, 'm')
});

declare module 'lucia' {
	interface Register {
		Lucia: typeof lucia;
	}
}
