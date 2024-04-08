import { db } from '$lib/db/db.server';
import { eq } from 'drizzle-orm';
import { C2Sessions } from '$lib/db/schema/c2_sessions';

export const getAllSessionsForProject = async (projectId: string) => {
	return db.query.C2Sessions.findMany({
		where: eq(C2Sessions.project_id, projectId)
	});
};
