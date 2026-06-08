import { baseClient } from "./client";

export interface Invite {
  id: string;
  server_id: string;
  code: string;
  url: string;
  created_by: string;
  max_uses: number | null;
  use_count: number;
  expires_at: string | null;
  created_at: string;
}

export interface JoinResponse {
  server_id: string;
  server_name: string;
}

class InvitesClient {
  async create(
    serverId: string,
    max_uses?: number | null,
    expires_in_hours?: number | null,
  ): Promise<Invite> {
    const res = await baseClient["request"](`/servers/${serverId}/invites`, {
      method: "POST",
      body: JSON.stringify({
        max_uses: max_uses ?? null,
        expires_in_hours: expires_in_hours ?? null,
      }),
    });
    return res.json();
  }

  async list(serverId: string): Promise<Invite[]> {
    const res = await baseClient["request"](`/servers/${serverId}/invites`);
    return res.json();
  }

  async delete(inviteId: string): Promise<void> {
    await baseClient["request"](`/invites/${inviteId}`, {
      method: "DELETE",
    });
  }

  async join(code: string): Promise<JoinResponse> {
    const res = await baseClient["request"](`/invites/${code}/join`, {
      method: "POST",
    });
    return res.json();
  }
}

export const invitesClient = new InvitesClient();
