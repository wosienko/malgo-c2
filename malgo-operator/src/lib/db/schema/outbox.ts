import {
	pgTable,
	timestamp,
	varchar,
	json,
	serial,
	customType,
	primaryKey,
	bigint
} from 'drizzle-orm/pg-core';
import { sql } from 'drizzle-orm';

const xid8 = customType<{ data: string; notNull: true; default: false }>({
	dataType() {
		return 'xid8';
	}
});

export const watermill_events_to_forward = pgTable(
	'watermill_events_to_forward',
	{
		offset: serial('offset'),
		uuid: varchar('uuid', { length: 36 }).notNull(),
		created_at: timestamp('created_at')
			.notNull()
			.default(sql`CURRENT_TIMESTAMP`),
		payload: json('payload').default(null),
		metadata: json('metadata').default(null),
		transaction_id: xid8('transaction_id').notNull()
	},
	(table) => ({
		pk: primaryKey({ columns: [table.transaction_id, table.offset] })
	})
);

export const watermill_offsets_events_to_forward = pgTable(
	'watermill_offsets_events_to_forward',
	{
		consumer_group: varchar('consumer_group', { length: 255 }).notNull(),
		offset_acked: bigint('offset_acked', { mode: 'bigint' }),
		last_processed_transaction_id: xid8('last_processed_transaction_id').notNull()
	},
	(table) => ({
		pk: primaryKey({ columns: [table.consumer_group] })
	})
);
