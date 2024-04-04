import { z } from 'zod';

const fieldSchema = z.string().trim().min(1);
const emailSchema = fieldSchema.email();
const passwordSchema = z.string().trim().min(8);

const uuidSchema = z.string().uuid();

type UuidSchema = z.infer<typeof uuidSchema>;

const dateSchema = z.coerce.date();

type DateSchema = z.infer<typeof dateSchema>;

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

const passwordChangeSchema = z
	.object({
		currentPassword: passwordSchema,
		newPassword: passwordSchema,
		newPasswordConfirmation: passwordSchema
	})
	.refine((data) => data.newPassword === data.newPasswordConfirmation, {
		message: "New passwords don't match",
		path: ['New Password Confirmation']
	});

type PasswordChangeSchema = z.infer<typeof passwordChangeSchema>;

const projectSchema = z
	.object({
		name: fieldSchema,
		startDate: dateSchema,
		endDate: dateSchema
	})
	.refine((data) => data.endDate > data.startDate, {
		message: 'End date must be after start date',
		path: ['End Date']
	});

type ProjectSchema = z.infer<typeof projectSchema>;

const userProjectSchema = z.object({
	users: z.array(uuidSchema)
});

type UserProjectSchema = z.infer<typeof userProjectSchema>;

export {
	fieldSchema,
	emailSchema,
	passwordSchema,
	uuidSchema,
	dateSchema,
	loginSchema,
	registerSchema,
	updateUserSchema,
	adminPasswordChangeSchema,
	adminRegisterSchema,
	passwordChangeSchema,
	projectSchema,
	userProjectSchema
};
export type {
	UuidSchema,
	DateSchema,
	LoginSchema,
	RegisterSchema,
	UpdateUserSchema,
	AdminPasswordChangeSchema,
	AdminRegisterSchema,
	PasswordChangeSchema,
	ProjectSchema,
	UserProjectSchema
};
