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
)

// Event is the WebSocket message envelope
type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MessageCreateData represents a new message event
type MessageCreateData struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TypingStartData represents a typing indicator event
type TypingStartData struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
