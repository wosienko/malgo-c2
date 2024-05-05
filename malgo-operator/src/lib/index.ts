// place files you want to import through the `$lib` alias in this folder.
export const version = '1.0.0';

export type ApiError = {
	message: string;
};

export interface User {
	id: string;
	name: string;
	surname: string;
	email: string;
}

export interface UserWithRoles extends User {
	admin: boolean;
	operator: boolean;
}

export type UsersWithRoles = {
	users: UserWithRoles[];
	count: number;
};
