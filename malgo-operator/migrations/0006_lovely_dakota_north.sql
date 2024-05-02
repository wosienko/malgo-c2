CREATE TABLE IF NOT EXISTS "watermill_events_to_forward" (
	"offset" serial NOT NULL,
	"uuid" varchar(36) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"payload" json DEFAULT 'null'::json,
	"metadata" json DEFAULT 'null'::json,
	"transaction_id" "xid8" NOT NULL,
	CONSTRAINT "watermill_events_to_forward_transaction_id_offset_pk" PRIMARY KEY("transaction_id","offset")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "watermill_offsets_events_to_forward" (
	"consumer_group" varchar(255) NOT NULL,
	"offset_acked" bigint,
	"last_processed_transaction_id" "xid8" NOT NULL,
	CONSTRAINT "watermill_offsets_events_to_forward_consumer_group_pk" PRIMARY KEY("consumer_group")
);
