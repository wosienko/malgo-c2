import type { Actions } from './$types';
import { registerSchema } from '$lib/validationSchemas';
import { fail } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';

import { db } from '$lib/db/db.server';
import { Users, UserRoles } from '$lib/db/schema/users';
import { Argon2id } from 'oslo/password';

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const userRegister = Object.fromEntries(data.entries());

		const result = registerSchema.safeParse(userRegister);
		if (!result.success) {
			return fail(400, {
				issues: result.error.issues
			});
		}

		if (result.data.password !== result.data.passwordConfirmation) {
			const issue: ZodIssue = {
				code: 'custom',
				message: 'Passwords do not match',
				path: ['Confirm Password']
			};

			return fail(400, {
				issues: [issue]
			});
		}

		try {
			await db.transaction(async (tx) => {
				const newUser = await tx
					.insert(Users)
					.values({
						name: result.data.name,
						surname: result.data.surname,
						email: result.data.email,
						password: await new Argon2id().hash(result.data.password)
					})
					.returning({ id: Users.id });

				const rolesAvailable = await tx.query.Roles.findMany({
					columns: {
						id: true
					}
				});
				for (const role of rolesAvailable) {
					await tx.insert(UserRoles).values({
						user_id: newUser[0].id,
						role_id: role.id
					});
				}
			});
		} catch (e) {
			return fail(500, {
				issues: [
					{
						code: 'custom',
						message: 'User with this email already exists',
						path: ['Email']
					} satisfies ZodIssue
				]
			});
		}

		return {
			message: 'User registered successfully'
		};
	}
} satisfies Actions;
