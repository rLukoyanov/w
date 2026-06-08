package models

import "time"

type ServerInvite struct {
	ID        string     `json:"id" db:"id"`
	ServerID  string     `json:"server_id" db:"server_id"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	Code      string     `json:"code" db:"code"`
	MaxUses   *int       `json:"max_uses" db:"max_uses"`
	UseCount  int        `json:"use_count" db:"use_count"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
