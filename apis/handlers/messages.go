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

func (h *MessagesHandler) messageToItem(msg *models.Message) (response.MessageItem, error) {
	item := response.MessageItem{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		AuthorID:  msg.AuthorID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}

	author, err := h.store.Users().GetByID(msg.AuthorID)
	if err != nil || author == nil {
		item.Author = response.AuthorInfo{
			ID: msg.AuthorID,
		}
	} else {
		item.Author = response.AuthorInfo{
			ID:       author.ID,
			Username: author.Username,
			Email:    author.Email,
		}
	}

	return item, nil
}

func (h *MessagesHandler) GetByChannelID(c *fiber.Ctx) error {
	channelID := c.Params("id")

	limit := 50
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

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

	var messages []*models.Message
	if beforeParam := c.Query("before"); beforeParam != "" {
		before, err := time.Parse(time.RFC3339, beforeParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid before timestamp, use RFC3339 format",
			})
		}
		messages, err = h.store.Messages().GetByChannelIDBefore(channelID, before, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get messages",
			})
		}
	} else {
		messages, err = h.store.Messages().GetByChannelID(channelID, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get messages",
			})
		}
	}

	items := make([]response.MessageItem, 0, len(messages))
	for _, msg := range messages {
		item, err := h.messageToItem(msg)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return c.JSON(response.Message{
		Messages: items,
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

	msg := &models.Message{
		ID:        uuid.New().String(),
		ChannelID: channelID,
		AuthorID:  userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.store.Messages().Create(msg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create message",
		})
	}

	item, err := h.messageToItem(msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to build response",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
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

	item, err := h.messageToItem(msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to build response",
		})
	}

	return c.JSON(item)
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
