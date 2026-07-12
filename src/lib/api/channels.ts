import { pb } from "./pocketbase";
import type { RecordModel } from "pocketbase";

export type ChannelType = "text" | "voice";

export interface Channel extends RecordModel {
  name: string;
  type: ChannelType;
  server: string;
  sort_order: number;
}

export const channels = {
  async list(serverId: string) {
    return await pb.collection("channels").getFullList<Channel>({
      filter: `server_id = "${serverId}"`,
    });
  },

  async get(id: string) {
    return await pb.collection("channels").getOne<Channel>(id);
  },
};
