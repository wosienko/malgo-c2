import { pgTable, text, timestamp, uuid } from 'drizzle-orm/pg-core';
import { Users } from './users';

export const Sessions = pgTable('sessions', {
	id: text('id').primaryKey(),
	userId: uuid('user_id')
		.notNull()
		.references(() => Users.id, { onDelete: 'cascade' }),
	expiresAt: timestamp('expires_at', {
		withTimezone: true,
		mode: 'date'
	}).notNull()
});
