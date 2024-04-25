DO $$ BEGIN
 CREATE TYPE "c2_command_status" AS ENUM('created', 'sending', 'sent', 'retrieving', 'success', 'error', 'canceled');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 CREATE TYPE "c2_command_type" AS ENUM('command', 'download', 'upload', 'settings');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "c2_commands" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"session_id" uuid NOT NULL,
	"type" "c2_command_type" DEFAULT 'command' NOT NULL,
	"status" "c2_command_status" DEFAULT 'created' NOT NULL,
	"command" text NOT NULL,
	"result_size" integer DEFAULT 0 NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"operator_id" uuid
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "c2_result_chunks" (
	"command_id" uuid NOT NULL,
	"result_chunk" text DEFAULT '' NOT NULL,
	"chunk_offset" integer NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	CONSTRAINT "c2_result_chunks_command_id_chunk_offset_pk" PRIMARY KEY("command_id","chunk_offset")
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "c2_commands" ADD CONSTRAINT "c2_commands_session_id_c2_sessions_id_fk" FOREIGN KEY ("session_id") REFERENCES "c2_sessions"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "c2_commands" ADD CONSTRAINT "c2_commands_operator_id_users_id_fk" FOREIGN KEY ("operator_id") REFERENCES "users"("id") ON DELETE set null ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "c2_result_chunks" ADD CONSTRAINT "c2_result_chunks_command_id_c2_commands_id_fk" FOREIGN KEY ("command_id") REFERENCES "c2_commands"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
