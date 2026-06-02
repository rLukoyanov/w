import { baseClient } from './client';
import type { Message } from './types';

class MessagesClient {
	async get(channelId: string, limit: number = 50): Promise<Message[]> {
		return baseClient['request']<Message[]>(`/channels/${channelId}/messages?limit=${limit}`);
	}

	async create(channelId: string, content: string): Promise<Message> {
		return baseClient['request']<Message>(`/channels/${channelId}/messages`, {
			method: 'POST',
			body: JSON.stringify({ content })
		});
	}
}

export const messagesClient = new MessagesClient();
