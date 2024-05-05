import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { uuidSchema } from '$lib/validationSchemas';
import { getLatestCommandForSession } from '$lib/services/c2-sessions-service';

export const GET: RequestHandler = async ({ params }) => {
	const { id, sessionid } = params;
	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	return json(await getLatestCommandForSession(sessionId.data));
};
