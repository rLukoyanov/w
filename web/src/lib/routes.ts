export const ROUTES = {
  HOME: "/",
  APP: {
    INDEX: "/app",
    SERVERS: "/app/servers",
    SETTINGS: "/app/settings",
  },
  AUTH: {
    LOGIN: "/auth/login",
    REGISTER: "/auth/register",
  },
  SERVER: {
    INDEX: "/app/servers",
    DETAIL: (id: string) => `/app/servers/${id}`,
    CHANNEL: (serverId: string, channelId: string) =>
      `/app/servers/${serverId}/channels/${channelId}`,
  },
  SETTINGS: "/app/settings",
} as const;
