package ws

import (
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/store"
	"github.com/rs/zerolog"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

// Client represents a single WebSocket connection
type Client struct {
	// User ID of the connected client
	UserID string

	// WebSocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan *Event

	// Hub instance
	hub *Hub

	// Store for database access
	store store.Store

	// Logger
	logger zerolog.Logger
}

// NewClient creates a new Client instance
func NewClient(userID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		UserID: userID,
		conn:   conn,
		send:   make(chan *Event, 256),
		hub:    hub,
		store:  hub.store,
		logger: hub.logger.With().Str("user_id", userID).Logger(),
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var event Event
		if err := c.conn.ReadJSON(&event); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}

		c.logger.Debug().
			Str("event_type", string(event.Type)).
			Msg("Received WebSocket event")

		if err := c.handleEvent(&event); err != nil {
			c.logger.Error().Err(err).
				Str("event_type", string(event.Type)).
				Msg("Error handling event")
			continue
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(event); err != nil {
				c.logger.Error().Err(err).Msg("Error writing WebSocket message")
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleEvent processes incoming WebSocket events
func (c *Client) handleEvent(event *Event) error {
	switch event.Type {
	case MessageCreate:
		return c.handleMessageCreate(event.Data)
	case TypingStart:
		return c.handleTypingStart(event.Data)
	default:
		c.logger.Warn().Str("event_type", string(event.Type)).Msg("Unknown event type")
		return nil
	}
}

// handleMessageCreate handles MESSAGE_CREATE events
func (c *Client) handleMessageCreate(data json.RawMessage) error {
	var msgData struct {
		ChannelID string `json:"channel_id"`
		Content   string `json:"content"`
	}

	if err := json.Unmarshal(data, &msgData); err != nil {
		return err
	}

	// TODO: Validate user has access to channel

	// Create message in database
	msg := &models.Message{
		ID:        uuid.New().String(),
		ChannelID: msgData.ChannelID,
		AuthorID:  c.UserID,
		Content:   msgData.Content,
		CreatedAt: time.Now(),
	}

	if err := c.store.Messages().Create(msg); err != nil {
		return err
	}

	// Broadcast to all clients in channel
	broadcastData := MessageCreateData{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		AuthorID:  msg.AuthorID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}

	dataBytes, err := json.Marshal(broadcastData)
	if err != nil {
		return err
	}

	c.hub.BroadcastToChannel(msg.ChannelID, &Event{
		Type: MessageCreate,
		Data: dataBytes,
	})

	return nil
}

// handleTypingStart handles TYPING_START events
func (c *Client) handleTypingStart(data json.RawMessage) error {
	var typingData TypingStartData
	if err := json.Unmarshal(data, &typingData); err != nil {
		return err
	}

	// Just broadcast typing indicator, don't persist
	c.hub.BroadcastToChannel(typingData.ChannelID, &Event{
		Type: TypingStart,
		Data: data,
	})

	return nil
}
