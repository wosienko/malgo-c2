import type { RequestHandler } from './$types';
import {
	type AdminPasswordChangeSchema,
	uuidSchema,
	adminPasswordChangeSchema
} from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import { changePassword } from '$lib/services/user-service';

export const PATCH: RequestHandler = async ({ params, request }) => {
	const { id } = params;

	const userId = uuidSchema.safeParse(id);
	if (!userId.success) {
		return json(userId.error.issues, { status: 400 });
	}

	const body: AdminPasswordChangeSchema = await request.json();

	const passwordChange = adminPasswordChangeSchema.safeParse(body);
	if (!passwordChange.success) {
		return json(passwordChange.error.issues, { status: 400 });
	}

	const result = await changePassword(userId.data, passwordChange.data.password);
	if (result) {
		return json({ message: result }, { status: 400 });
	}

	return json({ id });
};
