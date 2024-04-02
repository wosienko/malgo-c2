import { db } from '$lib/db/db.server';
import { eq, count } from 'drizzle-orm';
import { Roles, UserRoles, Users, type UserWithRoles } from '$lib/db/schema/users';
import { Argon2id } from 'oslo/password';
import type { LoginSchema, RegisterSchema } from '$lib/validationSchemas';
import { DEFAULT_ADMIN_EMAIL, DEFAULT_ADMIN_PASSWORD } from '$env/static/private';
import { ServerSettings } from '$lib/db/schema/server-settings';

export const findIdForLoginAttempt = async (userLogin: LoginSchema): Promise<string> => {
	let existingUser = await db.query.Users.findFirst({
		columns: {
			id: true,
			email: true,
			password: true
		},
		where: eq(Users.email, userLogin.email)
	});

	// if user does not exist, create a blank user object to hamper timing attacks
	if (!existingUser) existingUser = { id: '', email: '', password: await new Argon2id().hash('') };

	const validPassword = await new Argon2id().verify(existingUser.password, userLogin.password);
	return validPassword ? existingUser.id : '';
};

export const registerNewUser = async (userRegister: RegisterSchema): Promise<boolean> => {
	try {
		await db.transaction(async (tx) => {
			const newUser = await tx
				.insert(Users)
				.values({
					name: userRegister.name,
					surname: userRegister.surname,
					email: userRegister.email,
					password: await new Argon2id().hash(userRegister.password)
				})
				.returning({ id: Users.id });

			const rolesAvailable = await tx.query.Roles.findMany({
				columns: {
					id: true,
					name: true
				}
			});
			for (const role of rolesAvailable) {
				if (role.name !== 'User') continue;
				await tx.insert(UserRoles).values({
					user_id: newUser[0].id,
					role_id: role.id
				});
			}
		});
	} catch (e) {
		return false;
	}
	return true;
};

export const getUsersWithRoles = async (page: number, pageSize: number) => {
	return db.query.Users.findMany({
		columns: {
			id: true,
			name: true,
			surname: true,
			email: true
		},
		with: {
			UserRoles: {
				columns: {},
				with: {
					Role: {
						columns: {
							name: true
						}
					}
				}
			}
		},
		orderBy: Users.surname,
		limit: pageSize,
		offset: (page - 1) * pageSize
	}).then((users) => {
		return users.map((user) => {
			const userWithRoles: UserWithRoles = {
				id: user.id,
				name: user.name,
				surname: user.surname,
				email: user.email,
				admin: user.UserRoles.some((userRole) => userRole.Role.name === 'Admin'),
				operator: user.UserRoles.some((userRole) => userRole.Role.name === 'Operator')
			};
			return userWithRoles;
		});
	});
};

export const getCountOfUsers = async (): Promise<number> => {
	return db
		.select({ count: count() })
		.from(Users)
		.then((result) => result[0].count);
};

export const updateUserData = async (
	userId: string,
	name?: string,
	surname?: string,
	email?: string
) => {
	try {
		const updatedUser = await db
			.update(Users)
			.set({
				name,
				surname,
				email
			})
			.where(eq(Users.id, userId))
			.returning({ id: Users.id });
		return updatedUser.length > 0 ? '' : 'User not found';
	} catch (e) {
		return 'Error updating user. Verify the data and try again';
	}
};

export const deleteUser = async (userId: string) => {
	try {
		const deletedUser = await db
			.delete(Users)
			.where(eq(Users.id, userId))
			.returning({ id: Users.id });
		return deletedUser.length > 0 ? '' : 'User not found';
	} catch (e) {
		return 'Error deleting user. Verify the data and try again';
	}
};

export const createDefaultAdminAndRoles = async (): Promise<void> => {
	await db.transaction(async (tx) => {
		const isAdminConfigCompleted = await tx.query.ServerSettings.findFirst({
			columns: {
				value: true
			},
			where: eq(ServerSettings.key, 'adminConfigCompleted')
		});
		if (isAdminConfigCompleted) {
			console.log('Default admin and roles already created');
			return;
		}

		const adminRoleId = await tx
			.insert(Roles)
			.values({ name: 'Admin' })
			.returning({ id: Roles.id })
			.onConflictDoNothing();

		const operatorRoleId = await tx
			.insert(Roles)
			.values({ name: 'Operator' })
			.returning({ id: Roles.id })
			.onConflictDoNothing();

		const adminId = await tx
			.insert(Users)
			.values({
				name: 'Admin',
				surname: 'Admin',
				email: DEFAULT_ADMIN_EMAIL,
				password: await new Argon2id().hash(DEFAULT_ADMIN_PASSWORD)
			})
			.returning({ id: Users.id })
			.onConflictDoNothing();

		if (!adminId.length || !adminRoleId.length || !operatorRoleId.length) {
			console.log('Default roles and admin user already exist');
			return;
		}

		await tx
			.insert(UserRoles)
			.values([
				{
					user_id: adminId[0].id,
					role_id: adminRoleId[0].id
				},
				{
					user_id: adminId[0].id,
					role_id: operatorRoleId[0].id
				}
			])
			.onConflictDoNothing();

		await tx.insert(ServerSettings).values({
			key: 'adminConfigCompleted',
			value: 'true'
		});
	});
};
