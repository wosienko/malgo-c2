import { db } from '$lib/db/db.server';
import { count, desc, eq } from 'drizzle-orm';
import { C2Commands } from '$lib/db/schema/c2_commands';
import { hexToBytes } from '$lib/db/schema/c2_results';

export const getCommandsForSession = async (sessionId: string, page: number, pageSize: number) => {
	return db.query.C2Commands.findMany({
		where: eq(C2Commands.session_id, sessionId),
		with: {
			Operator: true,
			ResultChunks: true
		},
		orderBy: [desc(C2Commands.createdAt)],
		limit: pageSize,
		offset: (page - 1) * pageSize
	}).then((commands) =>
		commands.map((command) => {
			let resultChunks: string = '';

			if (command.status === 'completed') {
				let resultBytesBuffer: Buffer;
				for (const chunk of command.ResultChunks) {
					resultBytesBuffer = Buffer.from(hexToBytes(chunk.resultChunk));
					resultChunks += resultBytesBuffer.toString('utf8');
				}
			}

			return {
				id: command.id,
				type: command.type,
				status: command.status,
				operator: `${command.Operator?.name} ${command.Operator?.surname}`,
				command: command.command,
				created_at: formatDateAndTime(command.createdAt),
				result: resultChunks
			};
		})
	);
};

export const getCountOfCommandsForSession = async (sessionId: string): Promise<number> => {
	return db
		.select({ count: count() })
		.from(C2Commands)
		.where(eq(C2Commands.session_id, sessionId))
		.then((result) => result[0].count);
};

const formatDateAndTime = (d: Date): string => {
	// format DD.MM.YYYY HH:MM:SS
	return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
};
