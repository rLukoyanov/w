import { baseClient } from './client';
import type { Message } from './types';

class MessagesClient {
	async get(channelId: string, limit: number = 50, before?: string): Promise<Message[]> {
		let url = `/channels/${channelId}/messages?limit=${limit}`;
		if (before) url += `&before=${encodeURIComponent(before)}`;
		const res = await baseClient['request'](url);
		const data = await res.json();
		return data.messages ?? [];
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
