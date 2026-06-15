package sqlite

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/rLukoyanov/w/repository"
	_ "modernc.org/sqlite"
)

type Store struct {
	db            *sql.DB
	users         repository.UsersRepository
	servers       repository.ServersRepository
	channels      repository.ChannelsRepository
	messages      repository.MessagesRepository
	serverMembers repository.ServerMembersRepository
	serverInvites repository.ServerInvitesRepository
	attachments   repository.AttachmentsRepository
}

func New(dsn string) (*Store, error) {
	goose.SetDialect("sqlite")
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:            db,
		users:         &UsersRepository{db: db},
		servers:       &ServersRepository{db: db},
		channels:      &ChannelsRepository{db: db},
		messages:      &MessagesRepository{db: db},
		serverMembers: &ServerMembersRepository{db: db},
		serverInvites: &ServerInvitesRepository{db: db},
		attachments:   &AttachmentsRepository{db: db},
	}, nil
}

func (s *Store) Users() repository.UsersRepository {
	return s.users
}

func (s *Store) Servers() repository.ServersRepository {
	return s.servers
}

func (s *Store) Channels() repository.ChannelsRepository {
	return s.channels
}

func (s *Store) Messages() repository.MessagesRepository {
	return s.messages
}

func (s *Store) ServerMembers() repository.ServerMembersRepository {
	return s.serverMembers
}

func (s *Store) ServerInvites() repository.ServerInvitesRepository {
	return s.serverInvites
}

func (s *Store) Attachments() repository.AttachmentsRepository {
	return s.attachments
}

func (s *Store) Migrate(migrationsDir string) error {
	return goose.Up(s.db, migrationsDir)
}

func (s *Store) Close() error {
	return s.db.Close()
}
