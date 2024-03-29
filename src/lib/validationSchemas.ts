import { z } from 'zod';

const fieldSchema = z.string().trim().min(1);
const emailSchema = fieldSchema.email();
const passwordSchema = z.string().trim().min(8);

const loginSchema = z.object({
	email: emailSchema,
	password: passwordSchema
});

const registerSchema = z.object({
	email: emailSchema,
	password: passwordSchema,
	passwordConfirmation: passwordSchema,
	name: fieldSchema,
	surname: fieldSchema
});

export { fieldSchema, emailSchema, passwordSchema, loginSchema, registerSchema };
