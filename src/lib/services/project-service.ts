import { Projects, UserProjects } from '$lib/db/schema/projects';
import { db } from '$lib/db/db.server';
import { and, desc, count, eq, sql } from 'drizzle-orm';
import type { UuidSchema } from '$lib/validationSchemas';
import { Roles, UserRoles } from '$lib/db/schema/users';

export const getProjects = async (page: number, pageSize: number) => {
	return db.query.Projects.findMany({
		columns: {
			id: true,
			name: true,
			startDate: true,
			endDate: true,
			description: true
		},
		orderBy: [
			sql`CASE WHEN ${Projects.endDate} > CURRENT_DATE THEN 0 ELSE 1 END`,
			desc(Projects.startDate)
		],
		limit: pageSize,
		offset: (page - 1) * pageSize
	});
};

export const getCountOfProjects = async (): Promise<number> => {
	return db
		.select({ count: count() })
		.from(Projects)
		.then((result) => result[0].count);
};

export const getCountOfProjectsForOperator = async (operatorId: string): Promise<number> => {
	const projectsWithOperator = db
		.select({ project_id: UserProjects.project_id })
		.from(UserProjects)
		.where(eq(UserProjects.user_id, operatorId))
		.as('operator_projects');

	return db
		.select({ count: count() })
		.from(projectsWithOperator)
		.innerJoin(Projects, eq(Projects.id, projectsWithOperator.project_id))
		.then((result) => result[0].count);
};

export const getProjectsForOperator = async (
	operatorId: string,
	page: number,
	pageSize: number
) => {
	const projectsWithOperator = db
		.select({ project_id: UserProjects.project_id })
		.from(UserProjects)
		.where(eq(UserProjects.user_id, operatorId))
		.as('operator_projects');

	return db
		.select({
			id: Projects.id,
			name: Projects.name,
			startDate: Projects.startDate,
			endDate: Projects.endDate,
			description: Projects.description
		})
		.from(projectsWithOperator)
		.innerJoin(Projects, eq(Projects.id, projectsWithOperator.project_id))
		.where(
			sql`(${Projects.startDate} <= CURRENT_DATE AND ${Projects.endDate} >= CURRENT_DATE) OR (${Projects.endDate} + INTERVAL '14 days') >= CURRENT_DATE`
		)
		.orderBy(
			sql`CASE WHEN ${Projects.endDate} >= CURRENT_DATE THEN 0 ELSE 1 END`,
			desc(Projects.startDate)
		)
		.limit(pageSize)
		.offset((page - 1) * pageSize);
};

export const getAllProjectsForOperator = async (
	operatorId: string,
	page: number,
	pageSize: number
) => {
	const projectsWithOperator = db
		.select({ project_id: UserProjects.project_id })
		.from(UserProjects)
		.where(eq(UserProjects.user_id, operatorId))
		.as('operator_projects');

	return db
		.select({
			id: Projects.id,
			name: Projects.name,
			startDate: Projects.startDate,
			endDate: Projects.endDate,
			description: Projects.description
		})
		.from(projectsWithOperator)
		.innerJoin(Projects, eq(Projects.id, projectsWithOperator.project_id))
		.orderBy(
			sql`CASE WHEN ${Projects.endDate} >= CURRENT_DATE THEN 0 ELSE 1 END`,
			desc(Projects.startDate)
		)
		.limit(pageSize)
		.offset((page - 1) * pageSize);
};

export const createProject = async (
	name: string,
	startDate: Date,
	endDate: Date,
	description: string
) => {
	let result: UuidSchema | undefined = undefined;
	try {
		result = await db
			.insert(Projects)
			.values({
				name: name,
				startDate: startDate.toDateString(),
				endDate: endDate.toDateString(),
				description: description
			})
			.returning({ id: Projects.id })
			.then((result) => {
				if (result.length > 0) {
					return result[0].id;
				}
			});
	} catch (e) {
		console.error(e);
		return;
	}
	return result;
};

export const updateProjectData = async (
	projectId: string,
	name: string,
	startDate: Date,
	endDate: Date,
	description: string
) => {
	try {
		const updatedProjectId = await db
			.update(Projects)
			.set({
				name: name,
				startDate: startDate.toDateString(),
				endDate: endDate.toDateString(),
				description: description
			})
			.where(eq(Projects.id, projectId))
			.returning({ id: Projects.id });
		return updatedProjectId.length > 0 ? '' : 'Project not found';
	} catch (e) {
		return 'Error updating project data. Verify the data and try again';
	}
};

export const deleteProject = async (projectId: string) => {
	try {
		const deletedProjectId = await db
			.delete(Projects)
			.where(eq(Projects.id, projectId))
			.returning({ id: Projects.id });
		return deletedProjectId.length > 0 ? '' : 'Project not found';
	} catch (e) {
		return 'Error deleting project. Verify the data and try again';
	}
};

export const getOperatorsForProject = async (projectId: string) => {
	try {
		const operatorRoleId = db
			.select({ id: Roles.id })
			.from(Roles)
			.where(eq(Roles.name, 'Operator'))
			.as('operator_role_id');

		return await db
			.select({ id: UserRoles.user_id })
			.from(UserRoles)
			.innerJoin(UserProjects, eq(UserProjects.user_id, UserRoles.user_id))
			.innerJoin(operatorRoleId, eq(operatorRoleId.id, UserRoles.role_id))
			.where(eq(UserProjects.project_id, projectId))
			.then((result) => {
				return result.map((row) => row.id);
			});
	} catch (e) {
		console.error(e);
		return 'Error getting operators for project. Verify the data and try again';
	}
};

export const assignOperatorsToProject = async (projectId: string, operatorIds: string[]) => {
	try {
		const currentOperators = await getOperatorsForProject(projectId);
		if (currentOperators instanceof Array) {
			const operatorsToRemove = currentOperators.filter(
				(operatorId) => !operatorIds.includes(operatorId)
			);
			for (const operatorId of operatorsToRemove) {
				await db
					.delete(UserProjects)
					.where(and(eq(UserProjects.project_id, projectId), eq(UserProjects.user_id, operatorId)));
			}
		}

		const operatorRoleId = db
			.select({ id: Roles.id })
			.from(Roles)
			.where(eq(Roles.name, 'Operator'))
			.as('operator_role_id');
		const operators = await db
			.select({ id: UserRoles.user_id })
			.from(UserRoles)
			.innerJoin(operatorRoleId, eq(operatorRoleId.id, UserRoles.role_id))
			.then((result) => {
				return result.map((row) => row.id);
			})
			.then((ids) => {
				return ids.filter((operatorId) => operatorIds.includes(operatorId));
			});

		if (operators.length === 0) {
			return;
		}

		await db
			.insert(UserProjects)
			.values(
				operators.map((operatorId) => ({
					user_id: operatorId,
					project_id: projectId
				}))
			)
			.onConflictDoNothing();
	} catch (e) {
		console.error(e);
		return 'Error assigning operators to project. Verify the data and try again';
	}
};
