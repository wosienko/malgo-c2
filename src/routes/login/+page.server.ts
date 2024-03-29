import type { Actions } from './$types';
import { loginSchema } from '$lib/validationSchemas';
import { fail } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const userLogin = Object.fromEntries(data.entries());

		const result = loginSchema.safeParse(userLogin);
		if (!result.success) {
			return fail(400, {
				issues: result.error.issues
			});
		}

		const issue: ZodIssue = {
			code: 'custom',
			message: 'Invalid email or password',
			path: ['Login']
		};
		return fail(400, {
			issues: [issue]
		});
	}
} satisfies Actions;
