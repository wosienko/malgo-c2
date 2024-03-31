import { db } from '$lib/db/db.server';
import { UserRoles, Roles } from '$lib/db/schema/users';
import { eq } from 'drizzle-orm';

export const isUserAdmin = async (userId: string) => {
	return await isUserSomething(userId, 'Admin');
};

export const isUserOperator = async (userId: string) => {
	return await isUserSomething(userId, 'Operator');
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
