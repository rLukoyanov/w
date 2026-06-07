package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/apis/requests"
	"github.com/rLukoyanov/w/apis/response"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/store"
)

type ServersHandler struct {
	store store.Store
}

func NewServersHandler(store store.Store) *ServersHandler {
	return &ServersHandler{store: store}
}

func (h *ServersHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req requests.CreateServer
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	server := &models.Server{
		ID:        uuid.New().String(),
		Name:      req.Name,
		OwnerID:   userID,
		CreatedAt: time.Now(),
	}

	if err := h.store.Servers().Create(server); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create server",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(server)
}

func (h *ServersHandler) GetByID(c *fiber.Ctx) error {
	serverID := c.Params("id")

	server, err := h.store.Servers().GetByID(serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	if server == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "server not found",
		})
	}

	// Get channels for this server
	channels, err := h.store.Channels().GetByServerID(serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get channels",
		})
	}

	// Return empty array instead of null if no channels
	if channels == nil {
		channels = []*models.Channel{}
	}

	resp := &response.ServerWithChannels{
		Server:   server,
		Channels: channels,
	}

	return c.JSON(resp)
}

func (h *ServersHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	servers, err := h.store.Servers().GetByOwnerID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get servers",
		})
	}

	// Return empty array instead of null if no servers
	if servers == nil {
		servers = []*models.Server{}
	}

	return c.JSON(servers)
}
