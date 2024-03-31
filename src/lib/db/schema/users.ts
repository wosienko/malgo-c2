import { pgTable, serial, uuid, varchar } from 'drizzle-orm/pg-core';
import { sql, relations } from 'drizzle-orm';

export const Users = pgTable('users', {
	id: uuid('id')
		.primaryKey()
		.default(sql`gen_random_uuid()`),
	name: varchar('name').notNull(),
	surname: varchar('surname').notNull(),
	email: varchar('email').unique().notNull(),
	password: varchar('password').notNull()
});

export const Roles = pgTable('roles', {
	id: serial('id').primaryKey(),
	name: varchar('name').unique().notNull()
});

export const UserRoles = pgTable('user_roles', {
	id: uuid('id')
		.primaryKey()
		.default(sql`gen_random_uuid()`),
	user_id: uuid('user_id')
		.notNull()
		.references(() => Users.id),
	role_id: serial('role_id').notNull()
});

export const userRelations = relations(Users, ({ many }) => ({
	UserRoles: many(UserRoles)
}));

export const roleRelations = relations(Roles, ({ many }) => ({
	UserRoles: many(UserRoles)
}));

export const userRolesRelations = relations(UserRoles, ({ one }) => ({
	User: one(Users, {
		fields: [UserRoles.user_id],
		references: [Users.id]
	}),
	Role: one(Roles, {
		fields: [UserRoles.role_id],
		references: [Roles.id]
	})
}));
