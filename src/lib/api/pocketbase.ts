import PocketBase from 'pocketbase';
import { PUBLIC_POCKETBASE_URL } from '$env/static/public';

export const pb = new PocketBase(PUBLIC_POCKETBASE_URL);

export function createInstance() {
	return new PocketBase(PUBLIC_POCKETBASE_URL);
}

if (typeof document !== 'undefined') {
	pb.authStore.onChange(() => {
		document.cookie = pb.authStore.exportToCookie({ httpOnly: false });
	});
}
