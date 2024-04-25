import * as usersSchema from './schema/users';
import * as sessionsSchema from './schema/sessions';
import * as configSchema from './schema/server-settings';
import * as projectsSchema from './schema/projects';
import * as c2SessionsSchema from './schema/c2_sessions';
import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import { DATABASE_URL } from '$env/static/private';

const client = postgres(DATABASE_URL);
export const db = drizzle(client, {
	schema: {
		...usersSchema,
		...sessionsSchema,
		...configSchema,
		...projectsSchema,
		...c2SessionsSchema
	}
});

// Create default account and roles
