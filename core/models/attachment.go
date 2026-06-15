package models

import "time"

type Attachment struct {
	ID          string    `json:"id" db:"id"`
	ChannelID   string    `json:"channel_id" db:"channel_id"`
	AuthorID    string    `json:"author_id" db:"author_id"`
	Filename    string    `json:"filename" db:"filename"`
	Size        int64     `json:"size" db:"size"`
	ContentType string    `json:"content_type" db:"content_type"`
	StoragePath string    `json:"-" db:"storage_path"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
