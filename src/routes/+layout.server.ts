import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, url }) => {
	if (locals.user) {
		if (url.pathname === '/' || url.pathname.startsWith('/login') || url.pathname.startsWith('/register')) {
			redirect(303, '/dashboard');
		}
		return;
	}

	if (!url.pathname.startsWith('/login') && !url.pathname.startsWith('/register')) {
		redirect(303, '/login');
	}
};
