import { baseClient } from './client';
import type { Message } from './types';

class MessagesClient {
	async get(channelId: string, limit: number = 50): Promise<Message[]> {
		const res = await baseClient['request'](`/channels/${channelId}/messages?limit=${limit}`);
		return res.json();
	}

	async create(channelId: string, content: string): Promise<Message> {
		const res = await baseClient['request'](`/channels/${channelId}/messages`, {
			method: 'POST',
			body: JSON.stringify({ content })
		});
		return res.json();
	}
}

export const messagesClient = new MessagesClient();
