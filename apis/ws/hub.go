package ws

import (
	"sync"

	"github.com/rLukoyanov/w/store"
	"github.com/rs/zerolog"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients mapped by user ID
	clients map[string]*Client

	// Subscriptions: channelID → {userID → client}
	subscriptions map[string]map[string]*Client

	// Inbound messages from clients
	broadcast chan *Event

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Store for database access
	store store.Store

	// Logger
	logger zerolog.Logger

	// Mutex for thread-safe access to clients
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub(store store.Store, logger zerolog.Logger) *Hub {
	return &Hub{
		clients:       make(map[string]*Client),
		subscriptions: make(map[string]map[string]*Client),
		broadcast:     make(chan *Event, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		store:         store,
		logger:        logger.With().Str("component", "ws-hub").Logger(),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	h.logger.Info().Msg("WebSocket hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Close existing connection for the same user if any
			if existing, ok := h.clients[client.UserID]; ok {
				close(existing.send)
				existing.conn.Close()
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.logger.Info().
				Str("user_id", client.UserID).
				Int("total_clients", len(h.clients)).
				Msg("Client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				h.unsubscribeAll(client.UserID)
				delete(h.clients, client.UserID)
				close(client.send)
				h.logger.Info().
					Str("user_id", client.UserID).
					Int("total_clients", len(h.clients)).
					Msg("Client unregistered")
			}
			h.mu.Unlock()

		case event := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- event:
				default:
					h.logger.Warn().
						Str("user_id", client.UserID).
						Msg("Client send buffer full, dropping message")
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Subscribe adds a client to a channel's subscription list
func (h *Hub) Subscribe(userID, channelID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[userID]
	if !ok {
		return
	}

	// Unsubscribe from previous channel if any
	if client.SubscribedChannel != "" && client.SubscribedChannel != channelID {
		if subs, ok := h.subscriptions[client.SubscribedChannel]; ok {
			delete(subs, userID)
			if len(subs) == 0 {
				delete(h.subscriptions, client.SubscribedChannel)
			}
		}
	}

	if h.subscriptions[channelID] == nil {
		h.subscriptions[channelID] = make(map[string]*Client)
	}
	h.subscriptions[channelID][userID] = client
	client.SubscribedChannel = channelID

	h.logger.Debug().
		Str("user_id", userID).
		Str("channel_id", channelID).
		Msg("Client subscribed to channel")
}

// Unsubscribe removes a client from a channel's subscription list
func (h *Hub) Unsubscribe(userID, channelID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs, ok := h.subscriptions[channelID]; ok {
		delete(subs, userID)
		if len(subs) == 0 {
			delete(h.subscriptions, channelID)
		}
	}

	if client, ok := h.clients[userID]; ok && client.SubscribedChannel == channelID {
		client.SubscribedChannel = ""
	}

	h.logger.Debug().
		Str("user_id", userID).
		Str("channel_id", channelID).
		Msg("Client unsubscribed from channel")
}

// unsubscribeAll removes a client from all channel subscriptions (internal, lock held)
func (h *Hub) unsubscribeAll(userID string) {
	for channelID, subs := range h.subscriptions {
		delete(subs, userID)
		if len(subs) == 0 {
			delete(h.subscriptions, channelID)
		}
	}
	if client, ok := h.clients[userID]; ok {
		client.SubscribedChannel = ""
	}
}

// UnsubscribeAll removes a client from all channel subscriptions
func (h *Hub) UnsubscribeAll(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unsubscribeAll(userID)
}

// BroadcastToChannel sends an event to all clients subscribed to a channel
func (h *Hub) BroadcastToChannel(channelID string, event *Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subs, ok := h.subscriptions[channelID]
	if !ok {
		return
	}

	for _, client := range subs {
		select {
		case client.send <- event:
		default:
			h.logger.Warn().
				Str("user_id", client.UserID).
				Msg("Client send buffer full, dropping message")
		}
	}
}

// ConnectedUserIDs returns a copy of the currently connected user IDs
func (h *Hub) ConnectedUserIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		ids = append(ids, userID)
	}
	return ids
}

// ConnectedUserInfo holds user ID and their subscribed channel
type ConnectedUserInfo struct {
	UserID            string
	SubscribedChannel string
}

// ConnectedUsers returns info about all connected users
func (h *Hub) ConnectedUsers() []ConnectedUserInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]ConnectedUserInfo, 0, len(h.clients))
	for userID, client := range h.clients {
		users = append(users, ConnectedUserInfo{
			UserID:            userID,
			SubscribedChannel: client.SubscribedChannel,
		})
	}
	return users
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}
