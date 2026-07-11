import { createInstance } from '$lib/api/pocketbase';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const pb = createInstance();
	pb.authStore.loadFromCookie(event.request.headers.get('cookie') || '');

	if (pb.authStore.isValid) {
		try {
			await pb.collection('users').authRefresh();
		} catch {
			pb.authStore.clear();
		}
	}

	event.locals.user = pb.authStore.record;

	const response = await resolve(event);
	response.headers.append('set-cookie', pb.authStore.exportToCookie({ httpOnly: false }));
	return response;
};
