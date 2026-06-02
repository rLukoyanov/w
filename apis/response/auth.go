package response

import "github.com/rLukoyanov/w/core/models"

type Auth struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}
