import { baseClient } from "./client";
import type { Server } from "./types";

class ServersClient {
  async getAll(): Promise<Server[]> {
    const res = await baseClient["request"]("/servers");
    return res.json();
  }

  async create(name: string): Promise<Server> {
    const res = await baseClient["request"]("/servers", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
    return res.json();
  }

  async get(id: string): Promise<Server> {
    const res = await baseClient["request"](`/servers/${id}`);
    return res.json();
  }
}

export const serversClient = new ServersClient();
