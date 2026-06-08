package handlers

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rLukoyanov/w/apis/requests"
	"github.com/rLukoyanov/w/core/models"
	"github.com/rLukoyanov/w/store"
)

type InvitesHandler struct {
	store store.Store
}

func NewInvitesHandler(store store.Store) *InvitesHandler {
	return &InvitesHandler{store: store}
}

const codeChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateCode(length int) (string, error) {
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeChars))))
		if err != nil {
			return "", err
		}
		code[i] = codeChars[n.Int64()]
	}
	return string(code), nil
}

type InviteResponse struct {
	ID        string     `json:"id"`
	ServerID  string     `json:"server_id"`
	Code      string     `json:"code"`
	URL       string     `json:"url"`
	CreatedBy string     `json:"created_by"`
	MaxUses   *int       `json:"max_uses"`
	UseCount  int        `json:"use_count"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func inviteToResponse(invite *models.ServerInvite) InviteResponse {
	return InviteResponse{
		ID:        invite.ID,
		ServerID:  invite.ServerID,
		Code:      invite.Code,
		URL:       "/invite/" + invite.Code,
		CreatedBy: invite.CreatedBy,
		MaxUses:   invite.MaxUses,
		UseCount:  invite.UseCount,
		ExpiresAt: invite.ExpiresAt,
		CreatedAt: invite.CreatedAt,
	}
}

func (h *InvitesHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	serverID := c.Params("id")

	server, err := h.store.Servers().GetByID(serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if server == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "server not found"})
	}
	if server.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only the owner can create invites"})
	}

	var req requests.CreateInvite
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	code, err := generateCode(8)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate code"})
	}

	var expiresAt *time.Time
	if req.ExpiresInHours != nil && *req.ExpiresInHours > 0 {
		t := time.Now().Add(time.Duration(*req.ExpiresInHours) * time.Hour)
		expiresAt = &t
	}

	invite := &models.ServerInvite{
		ID:        uuid.New().String(),
		ServerID:  serverID,
		CreatedBy: userID,
		Code:      code,
		MaxUses:   req.MaxUses,
		UseCount:  0,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := h.store.ServerInvites().Create(invite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create invite"})
	}

	return c.Status(fiber.StatusCreated).JSON(inviteToResponse(invite))
}

func (h *InvitesHandler) ListByServer(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	serverID := c.Params("id")

	server, err := h.store.Servers().GetByID(serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if server == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "server not found"})
	}
	if server.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only the owner can view invites"})
	}

	invites, err := h.store.ServerInvites().GetByServerID(serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get invites"})
	}

	result := make([]InviteResponse, 0, len(invites))
	for _, inv := range invites {
		result = append(result, inviteToResponse(inv))
	}

	return c.JSON(result)
}

func (h *InvitesHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	inviteID := c.Params("id")

	// Get invite to check ownership
	// We need to look it up — but our repo only has GetByCode and GetByServerID
	// For now, get all invites for the user's servers and find this one
	// Better approach: add GetByID to repo, but let's use a workaround
	// Actually, let's check if the user owns the server by looking up the invite

	// We'll add a GetByID later — for now, iterate through user's servers
	// This is inefficient but works for MVP
	servers, err := h.store.Servers().GetByOwnerID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	found := false
	for _, s := range servers {
		invites, err := h.store.ServerInvites().GetByServerID(s.ID)
		if err != nil {
			continue
		}
		for _, inv := range invites {
			if inv.ID == inviteID {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invite not found"})
	}

	if err := h.store.ServerInvites().Delete(inviteID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete invite"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type JoinResponse struct {
	ServerID   string `json:"server_id"`
	ServerName string `json:"server_name"`
}

func (h *InvitesHandler) Join(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	code := c.Params("code")

	invite, err := h.store.ServerInvites().GetByCode(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if invite == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invite not found"})
	}

	// Check expiry
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "invite has expired"})
	}

	// Check max uses
	if invite.MaxUses != nil && invite.UseCount >= *invite.MaxUses {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "invite has reached max uses"})
	}

	// Check if already a member
	already, err := h.store.ServerMembers().IsMember(invite.ServerID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}
	if already {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "you are already a member of this server"})
	}

	// Add member
	if err := h.store.ServerMembers().Add(invite.ServerID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to join server"})
	}

	// Increment use count
	if err := h.store.ServerInvites().IncrementUseCount(invite.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update invite"})
	}

	// Get server name
	server, err := h.store.Servers().GetByID(invite.ServerID)
	if err != nil || server == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get server"})
	}

	return c.JSON(JoinResponse{
		ServerID:   server.ID,
		ServerName: server.Name,
	})
}
