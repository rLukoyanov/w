package models

import "time"

type Message struct {
	ID        string    `json:"id" db:"id"`
	ChannelID string    `json:"channel_id" db:"channel_id"`
	AuthorID  string    `json:"author_id" db:"author_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
