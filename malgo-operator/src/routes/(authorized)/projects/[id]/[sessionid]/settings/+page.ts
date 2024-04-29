import { uuidSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	const { id, sessionid } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	const sessionId = uuidSchema.safeParse(sessionid);
	if (!sessionId.success) {
		return json(sessionId.error.issues, { status: 400 });
	}

	return {
		session_name: await fetch(
			`/api/projects/${projectId.data}/sessions/${sessionId.data}/name`
		).then((r) => r.json())
	};
};
