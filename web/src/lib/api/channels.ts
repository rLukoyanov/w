import { baseClient } from './client';
import type { Channel } from './types';

class ChannelsClient {
	async create(serverId: string, name: string, type: string = 'text'): Promise<Channel> {
		return baseClient['request']<Channel>(`/servers/${serverId}/channels`, {
			method: 'POST',
			body: JSON.stringify({ name, type })
		});
	}

	async get(id: string): Promise<Channel> {
		return baseClient['request']<Channel>(`/channels/${id}`);
	}
}

export const channelsClient = new ChannelsClient();
