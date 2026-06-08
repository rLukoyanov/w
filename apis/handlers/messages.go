package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/apis/requests"
	"github.com/rLukoyanov/w/apis/response"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/store"
)

type MessagesHandler struct {
	store store.Store
}

func NewMessagesHandler(store store.Store) *MessagesHandler {
	return &MessagesHandler{store: store}
}

func (h *MessagesHandler) GetByChannelID(c *fiber.Ctx) error {
	channelID := c.Params("id")

	// Get limit from query params (default 50)
	limit := 50
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

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

	messages, err := h.store.Messages().GetByChannelID(channelID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get messages",
		})
	}

	return c.JSON(response.Message{
		Messages: messages,
	})
}

func (h *MessagesHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	channelID := c.Params("id")

	var req requests.CreateMessage
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

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

	message := &models.Message{
		ID:        uuid.New().String(),
		ChannelID: channelID,
		AuthorID:  userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.store.Messages().Create(message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}

func (h *MessagesHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	messageID := c.Params("id")

	msg, err := h.store.Messages().GetByID(messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	if msg == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "message not found",
		})
	}

	if msg.AuthorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not the author of this message",
		})
	}

	var req requests.UpdateMessage
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	msg.Content = req.Content

	if err := h.store.Messages().Update(msg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update message",
		})
	}

	return c.JSON(msg)
}

func (h *MessagesHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	messageID := c.Params("id")

	msg, err := h.store.Messages().GetByID(messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	if msg == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "message not found",
		})
	}

	if msg.AuthorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not the author of this message",
		})
	}

	if err := h.store.Messages().Delete(messageID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete message",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
