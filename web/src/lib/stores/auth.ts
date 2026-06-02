import { writable } from 'svelte/store';
import type { User } from '$lib/api';

interface AuthState {
	user: User | null;
	loading: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: false
	});

	return {
		subscribe,
		setUser: (user: User | null) => update((state) => ({ ...state, user })),
		setLoading: (loading: boolean) => update((state) => ({ ...state, loading })),
		logout: () => set({ user: null, loading: false })
	};
}

export const auth = createAuthStore();
