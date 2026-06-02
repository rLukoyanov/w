import { writable, type Writable } from 'svelte/store';
import { api } from './api';

export type EventType = 
  | 'MESSAGE_CREATE' 
  | 'MESSAGE_UPDATE' 
  | 'MESSAGE_DELETE'
  | 'TYPING_START'
  | 'TYPING_STOP'
  | 'PRESENCE_UPDATE';

export interface WSEvent {
  type: EventType;
  data: any;
}

export interface Message {
  id: string;
  channel_id: string;
  author_id: string;
  content: string;
  created_at: string;
}

class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectInterval = 3000;
  private reconnectTimer: number | null = null;
  
  public connected: Writable<boolean> = writable(false);
  public messages: Writable<Message[]> = writable([]);
  
  private eventHandlers: Map<EventType, Set<(data: any) => void>> = new Map();

  connect(token: string) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const wsUrl = `${protocol}//${host}/ws?token=${token}`;

    console.log('🔌 Connecting to WebSocket:', wsUrl);

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('✅ WebSocket connected');
      this.connected.set(true);
      
      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const wsEvent: WSEvent = JSON.parse(event.data);
        console.log('📨 WebSocket event:', wsEvent);
        
        this.handleEvent(wsEvent);
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err);
      }
    };

    this.ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('🔌 WebSocket closed');
      this.connected.set(false);
      this.scheduleReconnect(token);
    };
  }

  private scheduleReconnect(token: string) {
    if (this.reconnectTimer) return;
    
    console.log(`🔄 Reconnecting in ${this.reconnectInterval}ms...`);
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect(token);
    }, this.reconnectInterval) as unknown as number;
  }

  private handleEvent(event: WSEvent) {
    // Trigger registered handlers
    const handlers = this.eventHandlers.get(event.type);
    if (handlers) {
      handlers.forEach(handler => handler(event.data));
    }

    // Default handling
    switch (event.type) {
      case 'MESSAGE_CREATE':
        this.messages.update(msgs => [...msgs, event.data as Message]);
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

  send(event: WSEvent) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(event));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  sendMessage(channelId: string, content: string) {
    this.send({
      type: 'MESSAGE_CREATE',
      data: { channel_id: channelId, content }
    });
  }

  sendTyping(channelId: string) {
    this.send({
      type: 'TYPING_START',
      data: { channel_id: channelId, user_id: '' } // user_id will be set by server
    });
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.connected.set(false);
  }
}

export const wsClient = new WebSocketClient();
