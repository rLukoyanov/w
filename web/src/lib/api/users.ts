import { baseClient } from "./client";

export interface UserWithStatus {
  id: string;
  username: string;
  email: string;
  created_at: string;
  connected: boolean;
  subscribed_channel?: string;
  subscribed_channel_name?: string;
}

class UsersClient {
  async getConnected(): Promise<UserWithStatus[]> {
    const res = await baseClient["request"]("/users/connected");
    return res.json();
  }
}

export const usersClient = new UsersClient();
