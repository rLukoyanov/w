package store

import (
	"fmt"

	"github.com/rLukoyanov/w/store/sqlite"
)

func New(driver, dsn string) (Store, error) {
	switch driver {
	case "sqlite":
		return sqlite.New()
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}
