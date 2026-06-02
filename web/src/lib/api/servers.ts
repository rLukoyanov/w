import { baseClient } from './client';
import type { Server } from './types';

class ServersClient {
	async create(name: string): Promise<Server> {
		return baseClient['request']<Server>('/servers', {
			method: 'POST',
			body: JSON.stringify({ name })
		});
	}

	async get(id: string): Promise<Server> {
		return baseClient['request']<Server>(`/servers/${id}`);
	}
}

export const serversClient = new ServersClient();
