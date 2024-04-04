import { pgTable, uuid, varchar, date, primaryKey } from 'drizzle-orm/pg-core';
import { sql, relations } from 'drizzle-orm';
import { Users } from './users';

export const Projects = pgTable('projects', {
	id: uuid('id')
		.primaryKey()
		.default(sql`gen_random_uuid()`),
	name: varchar('name').notNull().unique(),
	startDate: date('start_date').notNull(),
	endDate: date('end_date').notNull()
});

export const UserProjects = pgTable(
	'user_projects',
	{
		user_id: uuid('user_id')
			.notNull()
			.references(() => Users.id, { onDelete: 'cascade' }),
		project_id: uuid('project_id')
			.notNull()
			.references(() => Projects.id, { onDelete: 'cascade' })
	},
	(table) => ({
		pk: primaryKey({ columns: [table.user_id, table.project_id] })
	})
);

export const projectRelations = relations(Projects, ({ many }) => ({
	UserProjects: many(UserProjects)
}));

export const userProjectsRelations = relations(UserProjects, ({ one }) => ({
	Project: one(Projects, {
		fields: [UserProjects.project_id],
		references: [Projects.id]
	}),
	User: one(Users, {
		fields: [UserProjects.user_id],
		references: [Users.id]
	})
}));
