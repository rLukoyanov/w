package store

type (
	Store interface {
		Migrate(migrationsDir string) error
		Close() error
	}
)
