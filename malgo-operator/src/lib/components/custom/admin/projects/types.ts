import type { UserWithRoles } from '$lib';

export type ProjectTableEntryType = {
	id: string;
	name: string;
	startDate: string;
	endDate: string;
	description: string;
};

export type CheckedUserType = {
	user: UserWithRoles;
	checked: boolean;
	changed: boolean;
};
