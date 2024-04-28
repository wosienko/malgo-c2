import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { isUserOperator } from '$lib/services/roles-service';
import { getCountOfCommandsForSession } from '$lib/services/c2-commands-service';
import { uuidSchema } from '$lib/validationSchemas';

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

	return json(await getCountOfCommandsForSession(sessionId.data));
};
