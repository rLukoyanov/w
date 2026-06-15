package repository

import (
	"time"

	"github.com/rLukoyanov/w/core/models"
)

type UsersRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type ServersRepository interface {
	Create(server *models.Server) error
	GetByID(id string) (*models.Server, error)
	GetByOwnerID(ownerID string) ([]*models.Server, error)
	GetAll() ([]*models.Server, error)
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
	GetByChannelIDBefore(channelID string, before time.Time, limit int) ([]*models.Message, error)
	CountByChannelID(channelID string) (int, error)
	Update(message *models.Message) error
	Delete(id string) error
}

type ServerMembersRepository interface {
	Add(serverID, userID string) error
	Remove(serverID, userID string) error
	GetMembers(serverID string) ([]*models.ServerMember, error)
	IsMember(serverID, userID string) (bool, error)
	GetServerIDsByUser(userID string) ([]string, error)
}

type ServerInvitesRepository interface {
	Create(invite *models.ServerInvite) error
	GetByCode(code string) (*models.ServerInvite, error)
	GetByServerID(serverID string) ([]*models.ServerInvite, error)
	IncrementUseCount(id string) error
	Delete(id string) error
}

type AttachmentsRepository interface {
	Create(attachment *models.Attachment) error
	GetByID(id string) (*models.Attachment, error)
	GetByChannelID(channelID string) ([]*models.Attachment, error)
	Delete(id string) error
}
