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
}

export const api = new ApiClient();
