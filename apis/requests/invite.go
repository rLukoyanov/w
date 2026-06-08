package requests

type CreateInvite struct {
	MaxUses        *int `json:"max_uses"`
	ExpiresInHours *int `json:"expires_in_hours"`
}
