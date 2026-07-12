import { pb } from './pocketbase';
import type { RecordModel } from 'pocketbase';

export interface Server extends RecordModel {
	name: string;
	description?: string;
	avatar?: string;
	owner: string;
}

export const servers = {
	async list() {
		return await pb.collection('servers').getFullList<Server>();
	},

	async get(id: string) {
		return await pb.collection('servers').getOne<Server>(id);
	},

	async create(data: { name: string; description?: string }) {
		return await pb.collection('servers').create<Server>(data);
	},

	async update(id: string, data: Partial<Server>) {
		return await pb.collection('servers').update<Server>(id, data);
	},

	async remove(id: string) {
		return await pb.collection('servers').delete(id);
	}
};
