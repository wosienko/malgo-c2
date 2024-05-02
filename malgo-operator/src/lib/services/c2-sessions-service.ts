import { db } from '$lib/db/db.server';
import { desc, eq } from 'drizzle-orm';
import { C2Sessions } from '$lib/db/schema/c2_sessions';
import { C2Commands } from '$lib/db/schema/c2_commands';
import type { Command } from '$lib/components/custom/command/command';
import { hexToBytes } from '$lib/db/schema/c2_results';

export const getAllSessionsForProject = async (projectId: string) => {
	return db.query.C2Sessions.findMany({
		where: eq(C2Sessions.project_id, projectId),
		orderBy: [desc(C2Sessions.createdAt)]
	});
};

export const getLatestCommandForSession = async (
	sessionId: string
): Promise<Command | undefined> => {
	// TODO: redo as a single query with a join and order by
	const result = await db.query.C2Commands.findFirst({
		where: eq(C2Commands.session_id, sessionId),
		orderBy: [desc(C2Commands.createdAt)],
		with: {
			Operator: true,
			ResultChunks: true
		}
	});

	if (!result) {
		return undefined;
	}

	let resultChunks: string = '';

	if (result.status === 'completed') {
		let resultBytesBuffer: Buffer;
		for (const chunk of result.ResultChunks) {
			resultBytesBuffer = Buffer.from(hexToBytes(chunk.resultChunk));
			resultChunks += resultBytesBuffer.toString('utf8');
		}
	}

	return {
		id: result.id,
		type: result.type,
		status: result.status,
		operator: `${result.Operator?.name} ${result.Operator?.surname}`,
		command: result.command,
		created_at: formatDateAndTime(result.createdAt),
		result: resultChunks
	};
};

export const getKeyValuesForSession = async (sessionId: string) => {
	const result = await db.query.C2Sessions.findFirst({
		columns: {
			data: true
		},
		where: eq(C2Sessions.id, sessionId)
	});

	return result ? result.data : {};
};

export const getSessionName = async (sessionId: string) => {
	const result = await db.query.C2Sessions.findFirst({
		columns: {
			name: true
		},
		where: eq(C2Sessions.id, sessionId)
	});

	return result ? result.name : '';
};

const formatDateAndTime = (d: Date): string => {
	// format DD.MM.YYYY HH:MM:SS
	return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
};
