import { pgTable, text, timestamp, uuid, jsonb } from 'drizzle-orm/pg-core';
import { sql } from 'drizzle-orm';
import { Projects } from './projects';

export const C2Sessions = pgTable('c2_sessions', {
	id: uuid('id')
		.primaryKey()
		.default(sql`gen_random_uuid()`),
	project_id: uuid('project_id')
		.notNull()
		.references(() => Projects.id, { onDelete: 'cascade' }),
	name: text('name')
		.notNull()
		.default(sql`gen_random_uuid()::text`),
	createdAt: timestamp('created_at', {
		withTimezone: true,
		mode: 'date'
	})
		.notNull()
		.default(sql`now()`),
	heartbeatAt: timestamp('heartbeat_at', {
		withTimezone: true,
		mode: 'date'
	})
		.notNull()
		.default(sql`now()`),
	data: jsonb('data')
		.notNull()
		.default(sql`'{}'`)
});
