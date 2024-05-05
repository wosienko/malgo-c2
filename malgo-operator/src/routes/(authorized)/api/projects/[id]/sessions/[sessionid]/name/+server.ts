import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { uuidSchema } from '$lib/validationSchemas';
import { getSessionName } from '$lib/services/c2-sessions-service';

export const GET: RequestHandler = async ({ params }) => {
	const { sessionid } = params;
	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	return json(await getSessionName(sessionId.data));
};
