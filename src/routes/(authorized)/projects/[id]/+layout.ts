import type { LayoutLoad } from './$types';
import { uuidSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';

export const load: LayoutLoad = async ({ params, data, fetch }) => {
	const { id } = params;

	const projectId = uuidSchema.safeParse(id);
	if (!projectId.success) {
		return json(projectId.error.issues, { status: 400 });
	}

	return {
		project: data.project,
		sessions: fetch(`/api/projects/${projectId.data}/sessions`).then((r) => r.json())
	};
};
