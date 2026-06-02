const API_URL = 'http://localhost:8090/api';

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

export interface Channel {
	id: string;
	server_id: string;
	name: string;
	type: string;
	created_at: string;
}

export interface Message {
	id: string;
	channel_id: string;
	author_id: string;
	content: string;
	created_at: string;
}

export class ApiClient {
	private token: string | null = null;

	constructor() {
		if (typeof window !== 'undefined') {
			this.token = localStorage.getItem('token');
		}
	}

	setToken(token: string) {
		this.token = token;
		if (typeof window !== 'undefined') {
			localStorage.setItem('token', token);
		}
	}

	clearToken() {
		this.token = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem('token');
		}
	}

	private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers
		};

		if (this.token) {
			headers['Authorization'] = `Bearer ${this.token}`;
		}

		const response = await fetch(`${API_URL}${endpoint}`, {
			...options,
			headers
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Unknown error' }));
			throw new Error(error.error || `HTTP ${response.status}`);
		}

		return response.json();
	}

	async login(data: LoginRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		this.setToken(response.token);
		return response;
	}

	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		this.setToken(response.token);
		return response;
	}

	async getMe(): Promise<User> {
		return this.request<User>('/auth/me');
	}

	async createServer(name: string): Promise<Server> {
		return this.request<Server>('/servers', {
			method: 'POST',
			body: JSON.stringify({ name })
		});
	}

	async getServer(id: string): Promise<Server> {
		return this.request<Server>(`/servers/${id}`);
	}

	async createChannel(serverId: string, name: string, type: string = 'text'): Promise<Channel> {
		return this.request<Channel>(`/servers/${serverId}/channels`, {
			method: 'POST',
			body: JSON.stringify({ name, type })
		});
	}

	async getChannel(id: string): Promise<Channel> {
		return this.request<Channel>(`/channels/${id}`);
	}

	async getMessages(channelId: string, limit: number = 50): Promise<Message[]> {
		return this.request<Message[]>(`/channels/${channelId}/messages?limit=${limit}`);
	}

	async createMessage(channelId: string, content: string): Promise<Message> {
		return this.request<Message>(`/channels/${channelId}/messages`, {
			method: 'POST',
			body: JSON.stringify({ content })
		});
	}

	getToken(): string | null {
		return this.token;
	}
}

export const api = new ApiClient();
