package core

import (
	"fmt"
	"os"

	"github.com/rLukoyanov/w/config"
	"github.com/rLukoyanov/w/store"
)

type App struct {
	Config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
	}
}

func (a *App) Run() error {
	store, err := store.New("sqlite", "test.db")
	if err != nil {
		return err
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("Current working dir:", dir)
	err = store.Migrate("./store/sqlite/migrations")
	if err != nil {
		return err
	}

	defer store.Close()

	return nil
}
