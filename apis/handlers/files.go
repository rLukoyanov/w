package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/storage"
	"github.com/rLukoyanov/w/store"
)

const maxUploadSize = 50 << 20 // 50 MB

type AttachmentResponse struct {
	ID          string `json:"id"`
	ChannelID   string `json:"channel_id"`
	AuthorID    string `json:"author_id"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
}

type FilesHandler struct {
	store   store.Store
	storage storage.FileStorage
}

func NewFilesHandler(st store.Store, fileStorage storage.FileStorage) *FilesHandler {
	return &FilesHandler{store: st, storage: fileStorage}
}

func (h *FilesHandler) Upload(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	channelID := c.Params("id")

	channel, err := h.store.Channels().GetByID(channelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if channel == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "channel not found"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing file"})
	}

	if file.Size > maxUploadSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "file too large (max 50MB)"})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file"})
	}
	defer f.Close()

	limited := io.LimitReader(f, maxUploadSize)
	storagePath, size, err := h.storage.Save(file.Filename, limited)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save file"})
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	attachment := &models.Attachment{
		ID:          uuid.New().String(),
		ChannelID:   channelID,
		AuthorID:    userID,
		Filename:    file.Filename,
		Size:        size,
		ContentType: contentType,
		StoragePath: storagePath,
		CreatedAt:   time.Now(),
	}

	if err := h.store.Attachments().Create(attachment); err != nil {
		h.storage.Delete(storagePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save attachment metadata"})
	}

	return c.Status(fiber.StatusCreated).JSON(AttachmentResponse{
		ID:          attachment.ID,
		ChannelID:   attachment.ChannelID,
		AuthorID:    attachment.AuthorID,
		Filename:    attachment.Filename,
		Size:        attachment.Size,
		ContentType: attachment.ContentType,
		URL:         "/api/attachments/" + attachment.ID,
		CreatedAt:   attachment.CreatedAt.Format(time.RFC3339),
	})
}

func (h *FilesHandler) Download(c *fiber.Ctx) error {
	attachmentID := c.Params("id")

	attachment, err := h.store.Attachments().GetByID(attachmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if attachment == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "attachment not found"})
	}

	reader, err := h.storage.Open(attachment.StoragePath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
	}
	defer reader.Close()

	c.Set("Content-Type", attachment.ContentType)
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, attachment.Filename))
	c.Set("Content-Length", fmt.Sprintf("%d", attachment.Size))

	return c.SendStream(reader)
}
