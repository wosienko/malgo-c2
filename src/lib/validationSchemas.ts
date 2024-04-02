import { z } from 'zod';

const fieldSchema = z.string().trim().min(1);
const emailSchema = fieldSchema.email();
const passwordSchema = z.string().trim().min(8);

const uuidSchema = z.string().uuid();

type UuidSchema = z.infer<typeof uuidSchema>;

const loginSchema = z.object({
	email: emailSchema,
	password: passwordSchema
});

type LoginSchema = z.infer<typeof loginSchema>;

const registerSchema = z
	.object({
		email: emailSchema,
		password: passwordSchema,
		passwordConfirmation: passwordSchema,
		name: fieldSchema,
		surname: fieldSchema
	})
	.refine((data) => data.passwordConfirmation === data.password, {
		message: "Passwords don't match",
		path: ['Password Confirmation']
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

const adminPasswordChangeSchema = z
	.object({
		password: passwordSchema,
		passwordConfirmation: passwordSchema
	})
	.refine((data) => data.password === data.passwordConfirmation, {
		message: "New passwords don't match",
		path: ['Password Confirmation']
	});

type AdminPasswordChangeSchema = z.infer<typeof adminPasswordChangeSchema>;

const adminRegisterSchema = z
	.object({
		name: fieldSchema,
		surname: fieldSchema,
		email: emailSchema,
		password: passwordSchema,
		passwordConfirmation: passwordSchema,
		admin: z.boolean(),
		operator: z.boolean()
	})
	.refine((data) => data.password === data.passwordConfirmation, {
		message: "Passwords don't match",
		path: ['Password Confirmation']
	});

type AdminRegisterSchema = z.infer<typeof adminRegisterSchema>;

export {
	fieldSchema,
	emailSchema,
	passwordSchema,
	uuidSchema,
	loginSchema,
	registerSchema,
	updateUserSchema,
	adminPasswordChangeSchema,
	adminRegisterSchema
};
export type {
	UuidSchema,
	LoginSchema,
	RegisterSchema,
	UpdateUserSchema,
	AdminPasswordChangeSchema,
	AdminRegisterSchema
};
