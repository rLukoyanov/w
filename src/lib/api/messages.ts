import { pb } from './pocketbase';
import type { RecordModel } from 'pocketbase';

export interface Message extends RecordModel {
	channel_id: string;
	user_id: string;
	content: string;
	images: string[];
	is_deleted: boolean;

	expand?: {
		user_id: {
			id: string;
			username: string;
			avatar?: string;
		};
	};
}

export const messages = {
	async list(channelId: string, page = 1, perPage = 10) {
		return await pb.collection('messages').getList<Message>(page, perPage, {
			filter: `channel_id = "${channelId}"`,
			sort: '-created',
			expand: 'user_id',
		});
	},

	async create(data: { content: string; channel_id: string }) {
		return await pb.collection('messages').create<Message>(data);
	},
};
