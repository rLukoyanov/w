package ws

import (
	"encoding/json"
	"time"
)

// EventType represents the type of WebSocket event
type EventType string

const (
	MessageCreate  EventType = "MESSAGE_CREATE"
	MessageUpdate  EventType = "MESSAGE_UPDATE"
	MessageDelete  EventType = "MESSAGE_DELETE"
	TypingStart    EventType = "TYPING_START"
	TypingStop     EventType = "TYPING_STOP"
	PresenceUpdate EventType = "PRESENCE_UPDATE"
	Subscribe      EventType = "SUBSCRIBE"
	Unsubscribe    EventType = "UNSUBSCRIBE"
)

// Event is the WebSocket message envelope
type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MessageCreateData represents a new message event
type MessageCreateData struct {
	ID             string    `json:"id"`
	ChannelID      string    `json:"channel_id"`
	AuthorID       string    `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// MessageUpdateData represents a message edit event
type MessageUpdateData struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

// MessageDeleteData represents a message delete event
type MessageDeleteData struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}

// TypingStartData represents a typing indicator event
type TypingStartData struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
}

// SubscribeData represents a channel subscription request
type SubscribeData struct {
	ChannelID string `json:"channel_id"`
}

// UnsubscribeData represents a channel unsubscription request
type UnsubscribeData struct {
	ChannelID string `json:"channel_id"`
}
