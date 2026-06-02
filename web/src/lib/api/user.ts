import { baseClient } from './client';
import type { AuthResponse, LoginRequest, RegisterRequest, User } from './types';

class UserClient {
	async login(data: LoginRequest): Promise<AuthResponse> {
		const response = await baseClient['request']<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		baseClient.setToken(response.token);
		return response;
	}

	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await baseClient['request']<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		baseClient.setToken(response.token);
		return response;
	}

	async getMe(): Promise<User> {
		return baseClient['request']<User>('/auth/me');
	}

	getToken(): string | null {
		return baseClient.getToken();
	}

	clearToken(): void {
		baseClient.clearToken();
	}
}

export const userClient = new UserClient();
