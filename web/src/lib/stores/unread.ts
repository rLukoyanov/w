import { writable, derived } from "svelte/store";

function key(serverId: string, channelId: string) {
  return `${serverId}.${channelId}`;
}

function createUnreadStore() {
  const counts = writable<Map<string, number>>(new Map());

  return {
    subscribe: counts.subscribe,

    get(serverId: string, channelId: string): number {
      let val = 0;
      counts.subscribe((m) => (val = m.get(key(serverId, channelId)) ?? 0))();
      return val;
    },

    increment(serverId: string, channelId: string) {
      counts.update((m) => {
        const k = key(serverId, channelId);
        m.set(k, (m.get(k) ?? 0) + 1);
        return new Map(m);
      });
    },

    markRead(serverId: string, channelId: string) {
      counts.update((m) => {
        m.delete(key(serverId, channelId));
        return new Map(m);
      });
    },

    markAllRead(serverId: string) {
      counts.update((m) => {
        const prefix = `${serverId}.`;
        for (const k of m.keys()) {
          if (k.startsWith(prefix)) m.delete(k);
        }
        return new Map(m);
      });
    },
  };
}

export const unread = createUnreadStore();
