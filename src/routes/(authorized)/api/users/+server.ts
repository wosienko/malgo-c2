import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ url }) => {
	const page = url.searchParams.get('page');
	const pageSize = url.searchParams.get('pageSize');

	return json({ page, pageSize });
};
