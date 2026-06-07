import { baseClient } from "./client";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "./types";

class UserClient {
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await baseClient["request"]("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
    const res: AuthResponse = await response.json();
    baseClient.setToken(res.token);
    return res;
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await baseClient["request"]("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
    const res: AuthResponse = await response.json();
    baseClient.setToken(res.token);
    return res;
  }

  async getMe(): Promise<User> {
    const res = await baseClient["request"]("/auth/me");
    return res.json();
  }

  getToken(): string | null {
    return baseClient.getToken();
  }

  clearToken(): void {
    baseClient.clearToken();
  }
}

export const userClient = new UserClient();
