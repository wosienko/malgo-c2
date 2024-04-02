import { z } from 'zod';

const fieldSchema = z.string().trim().min(1);
const emailSchema = fieldSchema.email();
const passwordSchema = z.string().trim().min(8);

const uuidSchema = z.string().uuid();

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

const updateUserSchema = z.object({
	email: emailSchema.optional(),
	name: fieldSchema.optional(),
	surname: fieldSchema.optional(),
	admin: z.boolean().optional(),
	operator: z.boolean().optional()
});

type UpdateUserSchema = z.infer<typeof updateUserSchema>;

export {
	fieldSchema,
	emailSchema,
	passwordSchema,
	uuidSchema,
	loginSchema,
	registerSchema,
	updateUserSchema
};
export type { LoginSchema, RegisterSchema, UpdateUserSchema };
