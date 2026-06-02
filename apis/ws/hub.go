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
		clients:    make(map[string]*Client),
		broadcast:  make(chan *Event, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		store:      store,
		logger:     logger.With().Str("component", "ws-hub").Logger(),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	h.logger.Info().Msg("WebSocket hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.logger.Info().
				Str("user_id", client.UserID).
				Int("total_clients", len(h.clients)).
				Msg("Client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
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
					// Client's send buffer is full, skip
					h.logger.Warn().
						Str("user_id", client.UserID).
						Msg("Client send buffer full, dropping message")
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToChannel sends an event to all clients who have access to a channel
func (h *Hub) BroadcastToChannel(channelID string, event *Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, client := range h.clients {
		// TODO: Check if user has access to the channel
		// For now, send to all connected clients
		_ = userID
		_ = channelID

		select {
		case client.send <- event:
		default:
			h.logger.Warn().
				Str("user_id", client.UserID).
				Msg("Client send buffer full, dropping message")
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}
