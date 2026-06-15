package response

import "time"

type AuthorInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AttachmentInfo struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type MessageItem struct {
	ID          string           `json:"id"`
	ChannelID   string           `json:"channel_id"`
	AuthorID    string           `json:"author_id"`
	Author      AuthorInfo       `json:"author"`
	Content     string           `json:"content"`
	Attachments []AttachmentInfo `json:"attachments,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type Message struct {
	Messages []MessageItem `json:"messages"`
}
