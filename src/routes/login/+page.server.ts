import type { Actions } from './$types';
import { loginSchema } from '$lib/validationSchemas';
import { fail, redirect } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
import { db } from '$lib/db/db.server';
import { Users } from '$lib/db/schema/users';
import { Argon2id } from 'oslo/password';
import { eq } from 'drizzle-orm';
import { lucia } from '$lib/auth.server';
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

		let existingUser = await db.query.Users.findFirst({
			columns: {
				id: true,
				email: true,
				password: true
			},
			where: eq(Users.email, result.data.email)
		});

		// if user does not exist, create a blank user object to hamper timing attacks
		if (!existingUser)
			existingUser = { id: '', email: '', password: await new Argon2id().hash('') };

		const passwordMatch = await new Argon2id().verify(existingUser.password, result.data.password);
		if (!passwordMatch) {
			const issue: ZodIssue = {
				code: 'custom',
				message: 'Incorrect email or password',
				path: ['Login']
			};

			return fail(400, {
				issues: [issue]
			});
		}

		const session = await lucia.createSession(existingUser.id, {});
		const sessionCookie = lucia.createSessionCookie(session.id);
		cookies.set(sessionCookie.name, sessionCookie.value, {
			path: '/',
			...sessionCookie.attributes
		});

		redirect(302, '/home');
	}
} satisfies Actions;
