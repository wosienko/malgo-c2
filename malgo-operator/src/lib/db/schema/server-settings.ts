import { pgTable, text } from 'drizzle-orm/pg-core';

export const ServerSettings = pgTable('server_settings', {
	key: text('key').primaryKey(),
	value: text('value').notNull()
});
