import type { Actions } from './$types';
import { registerSchema } from '$lib/validationSchemas';
import { fail, redirect } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';

import { db } from '$lib/db/db.server';
import { Users } from '$lib/db/schema/users';
import { Argon2id } from 'oslo/password';

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const userRegister = Object.fromEntries(data.entries());
		console.log(userRegister);

		const result = registerSchema.safeParse(userRegister);
		if (!result.success) {
			return fail(400, {
				issues: result.error.issues
			});
		}

		if (result.data.password !== result.data.passwordConfirmation) {
			console.log('bbb');

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
			await db.insert(Users).values({
				name: result.data.name,
				surname: result.data.surname,
				email: result.data.email,
				password: await new Argon2id().hash(result.data.password)
			});
		} catch (e) {
			return fail(500, {
				issues: [
					{
						code: 'custom',
						message: 'User with this email already exists',
						path: ['Email']
					}
				]
			});
		}

		// successful registration
		redirect(303, '/login');
	}
} satisfies Actions;
