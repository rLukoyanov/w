package response

import "time"

type AuthorInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type MessageItem struct {
	ID        string     `json:"id"`
	ChannelID string     `json:"channel_id"`
	AuthorID  string     `json:"author_id"`
	Author    AuthorInfo `json:"author"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}

type Message struct {
	Messages []MessageItem `json:"messages"`
}
