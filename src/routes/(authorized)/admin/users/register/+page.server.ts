import type { Actions } from './$types';
import { registerSchema } from '$lib/validationSchemas';
import { fail } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
import { registerNewUser } from '$lib/services/user-service';

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

		if (!(await registerNewUser(result.data))) {
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
