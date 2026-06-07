package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/apis/requests"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/store"
)

type ChannelsHandler struct {
	store store.Store
}

func NewChannelsHandler(store store.Store) *ChannelsHandler {
	return &ChannelsHandler{store: store}
}

func (h *ChannelsHandler) Create(c *fiber.Ctx) error {
	serverID := c.Params("id")

	var req requests.CreateChannel
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Verify server exists
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

	channel := &models.Channel{
		ID:        uuid.New().String(),
		ServerID:  serverID,
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: time.Now(),
	}

	if err := h.store.Channels().Create(channel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create channel",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(channel)
}

func (h *ChannelsHandler) GetByID(c *fiber.Ctx) error {
	channelID := c.Params("id")

	channel, err := h.store.Channels().GetByID(channelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	if channel == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "channel not found",
		})
	}

	return c.JSON(channel)
}

func (h *ChannelsHandler) GetByServerID(c *fiber.Ctx) error {
	serverID := c.Params("id")

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

	return c.JSON(channels)
}

func (h *ChannelsHandler) Delete(c *fiber.Ctx) error {
	channelID := c.Params("id")

	// Verify channel exists
	channel, err := h.store.Channels().GetByID(channelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	if channel == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "channel not found",
		})
	}

	if err := h.store.Channels().Delete(channelID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete channel",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
