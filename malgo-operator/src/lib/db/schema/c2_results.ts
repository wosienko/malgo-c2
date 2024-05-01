import { pgTable, uuid, integer, timestamp, primaryKey, customType } from 'drizzle-orm/pg-core';
import { C2Commands } from './c2_commands';
import { relations, sql } from 'drizzle-orm';

const bytea = customType<{ data: string; notNull: true; default: false }>({
	dataType() {
		return 'bytea';
	},
});

export const C2ResultChunks = pgTable(
	'c2_result_chunks',
	{
		command_id: uuid('command_id')
			.notNull()
			.references(() => C2Commands.id, { onDelete: 'cascade' }),
		resultChunk: bytea('result_chunk').notNull(),
		chunkOffset: integer('chunk_offset').notNull(),
		createdAt: timestamp('created_at', {
			withTimezone: true,
			mode: 'date'
		})
			.notNull()
			.default(sql`now()`)
	},
	(table) => ({
		pk: primaryKey({ columns: [table.command_id, table.chunkOffset] })
	})
);

export const c2ResultChunkRelations = relations(C2ResultChunks, ({ one }) => ({
	Command: one(C2Commands, {
		fields: [C2ResultChunks.command_id],
		references: [C2Commands.id]
	})
}));

export function hexToBytes(hex: string) {
	const bytes: number[] = [];
	for (let c = 0; c < hex.length; c += 2)
		// bytes.push(parseInt(hex.substr(c, 2), 16));
		bytes.push(parseInt(String.prototype.substring.call(hex, c, c + 2), 16));
	return bytes;
}