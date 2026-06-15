export interface User {
	id: string;
	username: string;
	email: string;
	created_at: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
}

export interface Server {
	id: string;
	name: string;
	owner_id: string;
	created_at: string;
}

export interface ServerWithChannels extends Server {
	channels: Channel[];
}

export interface Channel {
	id: string;
	server_id: string;
	name: string;
	type: string;
	created_at: string;
}

export interface Attachment {
	id: string;
	filename: string;
	size: number;
	content_type: string;
	url: string;
}

export interface Message {
  id: string;
  channel_id: string;
  author_id: string;
  author: {
    id: string;
    username: string;
    email: string;
  };
  content: string;
  attachments?: Attachment[];
  created_at: string;
}
