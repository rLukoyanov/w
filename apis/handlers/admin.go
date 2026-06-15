package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rLukoyanov/w/apis/ws"
	"github.com/rLukoyanov/w/store"
)

type StatsResponse struct {
	TotalUsers    int `json:"total_users"`
	TotalServers  int `json:"total_servers"`
	TotalChannels int `json:"total_channels"`
	TotalMessages int `json:"total_messages"`
	OnlineUsers   int `json:"online_users"`
}

type AdminUserItem struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type AdminServerItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	Channels  int    `json:"channels"`
	CreatedAt string `json:"created_at"`
}

type AdminHandler struct {
	store store.Store
	hub   *ws.Hub
}

func NewAdminHandler(st store.Store, hub *ws.Hub) *AdminHandler {
	return &AdminHandler{store: st, hub: hub}
}

func (h *AdminHandler) Stats(c *fiber.Ctx) error {
	users, _ := h.store.Users().GetAll()
	servers, _ := h.store.Servers().GetAll()

	totalChannels := 0
	totalMessages := 0

	return c.JSON(StatsResponse{
		TotalUsers:    len(users),
		TotalServers:  len(servers),
		TotalChannels: totalChannels,
		TotalMessages: totalMessages,
		OnlineUsers:   len(h.hub.ConnectedUsers()),
	})
}

func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.store.Users().GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list users"})
	}

	items := make([]AdminUserItem, 0, len(users))
	for _, u := range users {
		items = append(items, AdminUserItem{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}
	return c.JSON(items)
}

func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Role != "user" && req.Role != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "role must be 'user' or 'admin'"})
	}

	user, err := h.store.Users().GetByID(userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	user.Role = req.Role
	if err := h.store.Users().Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user"})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.store.Users().Delete(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminHandler) ListServers(c *fiber.Ctx) error {
	servers, err := h.store.Servers().GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list servers"})
	}

	items := make([]AdminServerItem, 0, len(servers))
	for _, s := range servers {
		ownerName := ""
		if owner, err := h.store.Users().GetByID(s.OwnerID); err == nil && owner != nil {
			ownerName = owner.Username
		}

		channels, _ := h.store.Channels().GetByServerID(s.ID)

		items = append(items, AdminServerItem{
			ID:        s.ID,
			Name:      s.Name,
			OwnerID:   s.OwnerID,
			OwnerName: ownerName,
			Channels:  len(channels),
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
		})
	}
	return c.JSON(items)
}

func (h *AdminHandler) DeleteServer(c *fiber.Ctx) error {
	serverID := c.Params("id")
	if err := h.store.Servers().Delete(serverID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "server not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
