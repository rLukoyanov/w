package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rLukoyanov/w/apis/ws"
	"github.com/rLukoyanov/w/store"
)

type UsersHandler struct {
	store store.Store
	hub   *ws.Hub
}

func NewUsersHandler(store store.Store, hub *ws.Hub) *UsersHandler {
	return &UsersHandler{store: store, hub: hub}
}

type UserWithStatus struct {
	ID                  string `json:"id"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	CreatedAt           string `json:"created_at"`
	Connected           bool   `json:"connected"`
	SubscribedChannel   string `json:"subscribed_channel,omitempty"`
	SubscribedChannelName string `json:"subscribed_channel_name,omitempty"`
}

func (h *UsersHandler) ListConnected(c *fiber.Ctx) error {
	connected := h.hub.ConnectedUsers()

	users := make([]UserWithStatus, 0, len(connected))

	for _, info := range connected {
		user, err := h.store.Users().GetByID(info.UserID)
		if err != nil || user == nil {
			continue
		}

		channelName := ""
		if info.SubscribedChannel != "" {
			channel, err := h.store.Channels().GetByID(info.SubscribedChannel)
			if err == nil && channel != nil {
				channelName = channel.Name
			}
		}

		users = append(users, UserWithStatus{
			ID:                    user.ID,
			Username:              user.Username,
			Email:                 user.Email,
			CreatedAt:             user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Connected:             true,
			SubscribedChannel:     info.SubscribedChannel,
			SubscribedChannelName: channelName,
		})
	}

	return c.JSON(users)
}
