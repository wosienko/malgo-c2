import { db } from '$lib/db/db.server';
import { eq } from 'drizzle-orm';
import { ServerSettings as ServerSettingsTable } from '$lib/db/schema/server-settings';

export enum ServerSettings {
	AdminConfigCompleted = 'adminConfigCompleted'
}

export const getServerSetting = async (key: ServerSettings): Promise<string | null> => {
	const setting = await db.query.ServerSettings.findFirst({
		where: eq(ServerSettingsTable.key, key)
	});
	return setting?.value ?? null;
};

export const setServerSetting = async (key: ServerSettings, value: string): Promise<void> => {
	await db.insert(ServerSettingsTable).values({ key, value });
};
