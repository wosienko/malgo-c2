import { db } from '$lib/db/db.server';
import { UserRoles, Roles } from '$lib/db/schema/users';
import { and, eq } from 'drizzle-orm';

export const isUserAdmin = async (userId: string) => {
	return await isUserSomething(userId, 'Admin');
};

export const isUserOperator = async (userId: string) => {
	return await isUserSomething(userId, 'Operator');
};

export const setAdmin = async (userId: string, isAdmin: boolean) => {
	return await setRole(userId, 'Admin', isAdmin);
};

export const setOperator = async (userId: string, isOperator: boolean) => {
	return await setRole(userId, 'Operator', isOperator);
};

const isUserSomething = async (userId: string, roleName: string) => {
	return await db
		.select({ role: Roles.name })
		.from(UserRoles)
		.where(eq(UserRoles.user_id, userId))
		.innerJoin(Roles, eq(UserRoles.role_id, Roles.id))
		.then((roles) => {
			return roles.some((role) => role.role === roleName);
		});
};

const setRole = async (userId: string, roleName: string, isRole: boolean) => {
	await db.transaction(async (tx) => {
		const role = await tx
			.select({ id: Roles.id })
			.from(Roles)
			.where(eq(Roles.name, roleName))
			.limit(1);

		if (role.length === 0) {
			return;
		}

		const userRole = await tx
			.select()
			.from(UserRoles)
			.where(and(eq(UserRoles.user_id, userId), eq(UserRoles.role_id, role[0].id)))
			.limit(1);

		const isUserRole = userRole.length > 0;

		if (isRole && !isUserRole) {
			await tx.insert(UserRoles).values({
				user_id: userId,
				role_id: role[0].id
			});
		} else if (!isRole && isUserRole) {
			await tx
				.delete(UserRoles)
				.where(and(eq(UserRoles.user_id, userId), eq(UserRoles.role_id, role[0].id)));
		}
	});
};
