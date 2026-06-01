package store

import "github.com/rLukoyanov/w/repository"

type Store interface {
	Migrate(migrationsDir string) error
	Close() error

	Users() repository.UsersRepository
	Servers() repository.ServersRepository
	Channels() repository.ChannelsRepository
	Messages() repository.MessagesRepository
}
