import { db } from '$lib/db/db.server';
import { desc, eq } from 'drizzle-orm';
import { C2Sessions } from '$lib/db/schema/c2_sessions';
import { C2Commands } from '$lib/db/schema/c2_commands';
import type { Command } from '$lib/components/custom/command/command';

export const getAllSessionsForProject = async (projectId: string) => {
	return db.query.C2Sessions.findMany({
		where: eq(C2Sessions.project_id, projectId),
		orderBy: [desc(C2Sessions.createdAt)]
	});
};

export const getLatestCommandForSession = async (
	sessionId: string
): Promise<Command | undefined> => {
	const result = await db.query.C2Commands.findFirst({
		where: eq(C2Commands.session_id, sessionId),
		orderBy: [desc(C2Commands.createdAt)],
		with: {
			Operator: true
		}
	});

	if (!result) {
		return undefined;
	}

	return {
		id: result.id,
		type: result.type,
		status: result.status,
		operator: `${result.Operator?.name} ${result.Operator?.surname}`,
		command: result.command,
		created_at: formatDateAndTime(result.createdAt),
		result: '' // TODO: add result
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

const formatDateAndTime = (d: Date): string => {
	// format DD.MM.YYYY HH:MM:SS
	return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
};
