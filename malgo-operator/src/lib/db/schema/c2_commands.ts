import { pgEnum, pgTable, uuid, text, integer, timestamp } from 'drizzle-orm/pg-core';
import { C2Sessions } from './c2_sessions';
import { relations, sql } from 'drizzle-orm';
import { Users } from './users';
import { C2ResultChunks } from './c2_results';

export const C2CommandTypeEnum = pgEnum('c2_command_type', [
	'command',
	'download',
	'upload',
	'settings'
]);

export const C2CommandStatusEnum = pgEnum('c2_command_status', [
	'created',
	'queried',
	'sending',
	'sent',
	'retrieving',
	'completed',
	'canceled'
]);

export const C2Commands = pgTable('c2_commands', {
	id: uuid('id')
		.primaryKey()
		.default(sql`gen_random_uuid()`),
	session_id: uuid('session_id')
		.notNull()
		.references(() => C2Sessions.id, { onDelete: 'cascade' }),
	type: C2CommandTypeEnum('type').notNull().default('command'),
	status: C2CommandStatusEnum('status').notNull().default('created'),
	command: text('command').notNull(),
	resultSize: integer('result_size').notNull().default(0),
	createdAt: timestamp('created_at', {
		withTimezone: true,
		mode: 'date'
	})
		.notNull()
		.default(sql`now()`),
	operator_id: uuid('operator_id').references(() => Users.id, { onDelete: 'set null' })
});

export const c2CommandRelations = relations(C2Commands, ({ one, many }) => ({
	Session: one(C2Sessions, {
		fields: [C2Commands.session_id],
		references: [C2Sessions.id]
	}),
	Operator: one(Users, {
		fields: [C2Commands.operator_id],
		references: [Users.id]
	}),
	ResultChunks: many(C2ResultChunks)
}));
