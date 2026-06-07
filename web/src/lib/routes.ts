export const ROUTES = {
  HOME: "/",
  AUTH: {
    LOGIN: "/auth/login",
    REGISTER: "/auth/register",
  },
  SERVER: {
    INDEX: "/servers",
    DETAIL: (id: string) => `/servers/${id}`,
    CHANNEL: (serverId: string, channelId: string) =>
      `/servers/${serverId}/channels/${channelId}`,
  },
  SETTINGS: "/settings",
} as const;
