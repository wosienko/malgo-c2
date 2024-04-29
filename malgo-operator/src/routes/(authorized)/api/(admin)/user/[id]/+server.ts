import type { RequestHandler } from './$types';
import { updateUserSchema, uuidSchema, type UpdateUserSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';
import { setAdmin, setOperator } from '$lib/services/roles-service';
import { deleteUser, updateUserData } from '$lib/services/user-service';

export const PATCH: RequestHandler = async ({ params, request }) => {
	const { id } = params;

	const userId = uuidSchema.safeParse(id);
	if (!userId.success) {
		return json(userId.error.issues, { status: 400 });
	}

	const body: UpdateUserSchema = await request.json();

	const result = updateUserSchema.safeParse(body);
	if (!result.success) {
		return json(result.error.issues, { status: 400 });
	}

	if (result.data.operator !== undefined) await setOperator(id, result.data.operator);
	if (result.data.admin !== undefined) await setAdmin(id, result.data.admin);

	const errorMessage = await updateUserData(
		userId.data,
		result.data.name,
		result.data.surname,
		result.data.email
	);
	if (errorMessage) {
		return json({ message: errorMessage }, { status: 400 });
	}

	return json({ id, ...result.data });
};

export const DELETE: RequestHandler = async ({ params }) => {
	const { id } = params;

	const userId = uuidSchema.safeParse(id);
	if (!userId.success) {
		return json(userId.error.issues, { status: 400 });
	}

	const errorMessage = await deleteUser(userId.data);
	if (errorMessage) {
		return json({ message: errorMessage }, { status: 400 });
	}

	return json({ id });
};
