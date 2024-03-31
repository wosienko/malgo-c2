import type { Actions } from './$types';
import { loginSchema } from '$lib/validationSchemas';
import { fail, redirect } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
import { findIdForLoginAttempt } from '$lib/services/user-service';
import { createSession } from '$lib/services/session-service';

export const actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const userLogin = Object.fromEntries(data.entries());

		const result = loginSchema.safeParse(userLogin);
		if (!result.success) {
			return fail(400, {
				issues: result.error.issues
			});
		}

		const userId = await findIdForLoginAttempt(result.data);
		if (userId === '') {
			const issue: ZodIssue = {
				code: 'custom',
				message: 'Incorrect email or password',
				path: ['Login']
			};

			return fail(400, {
				issues: [issue]
			});
		}

		await createSession(userId, cookies);

		redirect(302, '/home');
	}
} satisfies Actions;
