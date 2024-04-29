import type { RequestHandler } from './$types';
import { type PasswordChangeSchema, passwordChangeSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import { changePassword, validatePasswordForId } from '$lib/services/user-service';

export const PATCH: RequestHandler = async ({ locals, request }) => {
	const userId = locals.user!.id;

	const body: PasswordChangeSchema = await request.json();

	const passwordChange = passwordChangeSchema.safeParse(body);
	if (!passwordChange.success) {
		return json(passwordChange.error.issues, { status: 400 });
	}

	const isPasswordValid = await validatePasswordForId(userId, passwordChange.data.currentPassword);
	if (!isPasswordValid) {
		return json({ message: 'Invalid password' }, { status: 400 });
	}

	const result = await changePassword(userId, passwordChange.data.newPassword);
	if (result) {
		return json({ message: result }, { status: 400 });
	}

	return json({ userId });
};
