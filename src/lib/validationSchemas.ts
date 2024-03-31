import { z } from 'zod';

const fieldSchema = z.string().trim().min(1);
const emailSchema = fieldSchema.email();
const passwordSchema = z.string().trim().min(8);

const loginSchema = z.object({
	email: emailSchema,
	password: passwordSchema
});

type LoginSchema = z.infer<typeof loginSchema>;

const registerSchema = z.object({
	email: emailSchema,
	password: passwordSchema,
	passwordConfirmation: passwordSchema,
	name: fieldSchema,
	surname: fieldSchema
});

type RegisterSchema = z.infer<typeof registerSchema>;

export { fieldSchema, emailSchema, passwordSchema, loginSchema, registerSchema };
export type { LoginSchema, RegisterSchema };
