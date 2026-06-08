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
    SETTINGS: (id: string) => `/app/servers/${id}/settings`,
    CHANNEL: (serverId: string, channelId: string) =>
      `/app/servers/${serverId}/channels/${channelId}`,
  },
  INVITE: {
    JOIN: (code: string) => `/app/invite/${code}`,
  },
  SETTINGS: "/app/settings",
} as const;
