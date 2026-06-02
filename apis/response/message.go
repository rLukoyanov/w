package response

import "github.com/rLukoyanov/w/core/models"

type Message struct {
	Messages []*models.Message `json:"messages"`
}
