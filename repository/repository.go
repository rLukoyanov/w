package repository

import "github.com/rLukoyanov/w/core/models"

type UsersRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type ServersRepository interface {
	Create(server *models.Server) error
	GetByID(id string) (*models.Server, error)
	GetByOwnerID(ownerID string) ([]*models.Server, error)
	Update(server *models.Server) error
	Delete(id string) error
}

type ChannelsRepository interface {
	Create(channel *models.Channel) error
	GetByID(id string) (*models.Channel, error)
	GetByServerID(serverID string) ([]*models.Channel, error)
	Update(channel *models.Channel) error
	Delete(id string) error
}

type MessagesRepository interface {
	Create(message *models.Message) error
	GetByID(id string) (*models.Message, error)
	GetByChannelID(channelID string, limit int) ([]*models.Message, error)
	Update(message *models.Message) error
	Delete(id string) error
}
