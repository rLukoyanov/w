import { baseClient } from "./client";
import type { Channel } from "./types";

class ChannelsClient {
  async create(
    serverId: string,
    name: string,
    type: string = "text",
  ): Promise<Channel> {
    const res = await baseClient["request"](`/servers/${serverId}/channels`, {
      method: "POST",
      body: JSON.stringify({ name, type }),
    });
    return res.json();
  }

  async get(id: string): Promise<Channel> {
    const res = await baseClient["request"](`/channels/${id}`);
    return res.json();
  }

  async delete(id: string): Promise<number> {
    const res = await baseClient["request"](`/channels/${id}`, {
      method: "DELETE",
    });
    return res.status;
  }
}

export const channelsClient = new ChannelsClient();
