package sqlite

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(dsn string) (*Store, error) {
	goose.SetDialect("sqlite")
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Migrate(migrationsDir string) error {
	return goose.Up(s.db, migrationsDir)
}

func (s *Store) Close() error {
	return s.db.Close()
}
