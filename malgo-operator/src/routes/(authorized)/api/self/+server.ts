import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getBasicInformationForId, updateUserData } from '$lib/services/user-service';
import { updateUserSchema } from '$lib/validationSchemas';

export const GET: RequestHandler = async ({ locals }) => {
	const userId = locals.user!.id;

	return json(await getBasicInformationForId(userId));
};

export const PATCH: RequestHandler = async ({ locals, request }) => {
	const userId = locals.user!.id;

	const body = await request.json();

	const result = updateUserSchema.safeParse(body);
	if (!result.success) {
		return json(result.error.issues, { status: 400 });
	}

	const message = await updateUserData(
		userId,
		result.data.name,
		result.data.surname,
		result.data.email
	);
	if (message) {
		return json({ message }, { status: 400 });
	}

	return json({ id: userId });
};
