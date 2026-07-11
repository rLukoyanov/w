import { pb } from './pocketbase';
import type { RecordModel } from 'pocketbase';

export interface User extends RecordModel {
	email: string;
	name: string;
	avatar?: string;
	verified: boolean;
}

export const user = {
	async login(email: string, password: string) {
		return await pb.collection('users').authWithPassword<User>(email, password);
	},

	async register(email: string, password: string, passwordConfirm: string, name?: string) {
		return await pb.collection('users').create<User>({
			email,
			password,
			passwordConfirm,
			name
		});
	},

	async logout() {
		pb.authStore.clear();
	},

	async refresh() {
		return await pb.collection('users').authRefresh<User>();
	},

	get current() {
		return pb.authStore.record as User | null;
	},

	get isAuthenticated() {
		return pb.authStore.isValid;
	}
};
