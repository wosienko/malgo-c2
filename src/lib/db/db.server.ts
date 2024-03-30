import * as usersSchema from './schema/users';
import * as sessionsSchema from './schema/sessions';
import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import { DATABASE_URL } from '$env/static/private';

const client = postgres(DATABASE_URL);
export const db = drizzle(client, { schema: { ...usersSchema, ...sessionsSchema } });

// Create default account and roles
