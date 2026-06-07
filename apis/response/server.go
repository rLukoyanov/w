package response

import "github.com/rLukoyanov/w/core/models"

type ServerWithChannels struct {
	*models.Server
	Channels []*models.Channel `json:"channels"`
}
