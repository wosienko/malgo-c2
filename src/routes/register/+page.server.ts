import type { Actions } from './$types';
import { registerSchema } from '$lib/validationSchemas';
import { fail } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
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

			let issue: ZodIssue = {
				code: 'custom',
				message: 'Passwords do not match',
				path: ['Confirm Password']
			};

			return fail(400, {
				issues: [issue]
			});
		}
	}
} satisfies Actions;
