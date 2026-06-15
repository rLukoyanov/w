import { writable, type Writable } from "svelte/store";
import { userClient } from "./api/user";

export type EventType =
  | "MESSAGE_CREATE"
  | "MESSAGE_UPDATE"
  | "MESSAGE_DELETE"
  | "TYPING_START"
  | "TYPING_STOP"
  | "PRESENCE_UPDATE"
  | "SUBSCRIBE"
  | "UNSUBSCRIBE";

export interface WSEvent {
  type: EventType;
  data: any;
}

export interface MessagePayload {
  id: string;
  channel_id: string;
  author_id: string;
  author_username?: string;
  author?: {
    id: string;
    username: string;
    email: string;
  };
  content: string;
  created_at: string;
}

export interface TypingPayload {
  channel_id: string;
  user_id: string;
}

const WS_HOST =
  typeof window !== "undefined"
    ? import.meta.env.VITE_WS_HOST ?? `${window.location.hostname}:8090`
    : "localhost:8090";

class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectInterval = 3000;
  private reconnectTimer: number | null = null;
  private currentToken: string | null = null;
  private subscribedChannelId: string | null = null;

  public connected: Writable<boolean> = writable(false);
  public messages: Writable<MessagePayload[]> = writable([]);
  public typingUsers: Writable<Map<string, { user_id: string; channel_id: string; username?: string }>> = writable(new Map());

  private eventHandlers: Map<EventType, Set<(data: any) => void>> = new Map();

  connect(token: string) {
    if (this.ws?.readyState === WebSocket.OPEN) return;

    this.currentToken = token;
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${WS_HOST}/ws?token=${token}`;

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      this.connected.set(true);

      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }

      if (this.subscribedChannelId) {
        this.sendRaw({ type: "SUBSCRIBE", data: { channel_id: this.subscribedChannelId } });
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const wsEvent: WSEvent = JSON.parse(event.data);
        this.handleEvent(wsEvent);
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
      }
    };

    this.ws.onerror = () => {};

    this.ws.onclose = () => {
      this.connected.set(false);
      this.scheduleReconnect();
    };
  }

  private scheduleReconnect() {
    if (this.reconnectTimer || !this.currentToken) return;

    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect(this.currentToken!);
    }, this.reconnectInterval) as unknown as number;
  }

  private handleEvent(event: WSEvent) {
    const handlers = this.eventHandlers.get(event.type);
    if (handlers) {
      handlers.forEach((handler) => handler(event.data));
    }

    switch (event.type) {
      case "MESSAGE_CREATE":
        this.messages.update((msgs) => [...msgs, event.data as MessagePayload]);
        break;
    }
  }

  on(eventType: EventType, handler: (data: any) => void) {
    if (!this.eventHandlers.has(eventType)) {
      this.eventHandlers.set(eventType, new Set());
    }
    this.eventHandlers.get(eventType)!.add(handler);
  }

  off(eventType: EventType, handler: (data: any) => void) {
    const handlers = this.eventHandlers.get(eventType);
    if (handlers) {
      handlers.delete(handler);
    }
  }

  private sendRaw(event: WSEvent) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(event));
    }
  }

  sendMessage(channelId: string, content: string) {
    this.sendRaw({
      type: "MESSAGE_CREATE",
      data: { channel_id: channelId, content },
    });
  }

  private typingTimeout: number | null = null;

  sendTyping(channelId: string) {
    this.sendRaw({
      type: "TYPING_START",
      data: { channel_id: channelId },
    });

    if (this.typingTimeout) {
      clearTimeout(this.typingTimeout);
    }

    this.typingTimeout = setTimeout(() => {
      this.sendRaw({
        type: "TYPING_STOP",
        data: { channel_id: channelId },
      });
    }, 3000) as unknown as number;
  }

  sendTypingStop(channelId: string) {
    if (this.typingTimeout) {
      clearTimeout(this.typingTimeout);
      this.typingTimeout = null;
    }
    this.sendRaw({
      type: "TYPING_STOP",
      data: { channel_id: channelId },
    });
  }

  subscribe(channelId: string) {
    if (this.subscribedChannelId === channelId) return;

    if (this.subscribedChannelId) {
      this.sendRaw({ type: "UNSUBSCRIBE", data: { channel_id: this.subscribedChannelId } });
    }

    this.subscribedChannelId = channelId;
    this.sendRaw({ type: "SUBSCRIBE", data: { channel_id: channelId } });
  }

  unsubscribe() {
    if (this.subscribedChannelId) {
      this.sendRaw({ type: "UNSUBSCRIBE", data: { channel_id: this.subscribedChannelId } });
      this.subscribedChannelId = null;
    }
  }

  disconnect() {
    this.subscribedChannelId = null;
    this.currentToken = null;

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.typingTimeout) {
      clearTimeout(this.typingTimeout);
      this.typingTimeout = null;
    }

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.connected.set(false);
  }
}

export const wsClient = new WebSocketClient();

export const wsConnected = wsClient.connected;
export const wsMessages = wsClient.messages;
