import { pgTable, uuid, text, integer, timestamp } from 'drizzle-orm/pg-core';
import { C2Commands } from './c2_commands';
import { relations, sql } from 'drizzle-orm';

export const C2ResultChunks = pgTable('c2_result_chunks', {
	command_id: uuid('command_id')
		.notNull()
		.references(() => C2Commands.id, { onDelete: 'cascade' }),
	resultChunk: text('result_chunk').notNull().default(''),
	chunkSize: integer('chunk_size')
		.notNull()
		.default(sql`length(result_chunk)`),
	chunkOffset: integer('chunk_offset').notNull(),
	createdAt: timestamp('created_at', {
		withTimezone: true,
		mode: 'date'
	})
		.notNull()
		.default(sql`now()`)
});

export const c2ResultChunkRelations = relations(C2ResultChunks, ({ one }) => ({
	Command: one(C2Commands, {
		fields: [C2ResultChunks.command_id],
		references: [C2Commands.id]
	})
}));
