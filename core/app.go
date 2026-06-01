package core

import (
	"github.com/rLukoyanov/w/config"
	"github.com/rLukoyanov/w/store"
	"github.com/rs/zerolog/log"
)

type App struct {
	Config *config.Config
}

func New(cfg *config.Config) *App {
	if cfg == nil {
		cfg = &config.Config{}
	}

	if cfg.DBPath == "" {
		cfg.DBPath = "data.db"
	}

	return &App{
		Config: cfg,
	}
}

func (a *App) Run() error {
	store, err := store.New("sqlite", a.Config.DBPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize store")
		return err
	}

	log.Info().Msg("Start migrating database")
	err = store.Migrate("./store/sqlite/migrations")
	if err != nil {
		log.Error().Err(err).Msg("failed to apply migrations")
		return err
	}
	log.Info().Msg("migrations applied")

	defer store.Close()

	return nil
}
