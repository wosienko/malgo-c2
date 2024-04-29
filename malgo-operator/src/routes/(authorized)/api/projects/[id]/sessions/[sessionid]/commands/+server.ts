import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { isUserOperator } from '$lib/services/roles-service';
import { getCommandsForSession } from '$lib/services/c2-commands-service';
import { uuidSchema } from '$lib/validationSchemas';

export const GET: RequestHandler = async ({ url, locals, params }) => {
	const userId = locals.user!.id;
	if (!(await isUserOperator(userId))) {
		return json(
			{
				message: 'You are not do this query'
			},
			{
				status: 403
			}
		);
	}

	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 5;

	if (page < 1 || pageSize < 1) {
		return json(
			{
				message: 'Invalid page or pageSize'
			},
			{
				status: 400
			}
		);
	}

	const { sessionid } = params;
	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	return json(await getCommandsForSession(sessionId.data, page, pageSize));
};
