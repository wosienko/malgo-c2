import type { RequestHandler } from './$types';
import { isUserOperator } from '$lib/services/roles-service';
import { json } from '@sveltejs/kit';
import { uuidSchema } from '$lib/validationSchemas';
import { getSessionName } from '$lib/services/c2-sessions-service';

export const GET: RequestHandler = async ({ locals, params }) => {
	const userId = locals.user!.id;
	if (!(await isUserOperator(userId))) {
		return json(
			{
				message: 'You are not authorized query this endpoint'
			},
			{
				status: 403
			}
		);
	}

	const { sessionid } = params;
	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	return json(await getSessionName(sessionId.data));
};
