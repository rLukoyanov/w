package core

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rLukoyanov/w/apis"
	"github.com/rLukoyanov/w/apis/ws"
	"github.com/rLukoyanov/w/config"
	"github.com/rLukoyanov/w/storage"
	"github.com/rLukoyanov/w/store"
	"github.com/rs/zerolog/log"
)

type App struct {
	Config *config.Config
	store  store.Store
	server *apis.Server
}

func New(cfg *config.Config) *App {
	if cfg == nil {
		cfg = &config.Config{}
	}

	if cfg.DBPath == "" {
		cfg.DBPath = "data.db"
	}

	if cfg.Port == 0 {
		cfg.Port = 8090
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "change-me-in-production"
	}

	if cfg.UploadDir == "" {
		cfg.UploadDir = filepath.Join(filepath.Dir(cfg.DBPath), "uploads")
	}

	return &App{
		Config: cfg,
	}
}

func (a *App) Run() error {
	// Initialize store
	st, err := store.New("sqlite", a.Config.DBPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize store")
		return err
	}
	a.store = st
	defer a.store.Close()

	// Apply migrations
	log.Info().Msg("Start migrating database")
	err = a.store.Migrate("./store/sqlite/migrations")
	if err != nil {
		log.Error().Err(err).Msg("failed to apply migrations")
		return err
	}
	log.Info().Msg("migrations applied")

	// Initialize file storage
	fileStorage, err := storage.NewLocalStorage(a.Config.UploadDir)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize file storage")
		return err
	}

	// Initialize WebSocket hub
	hub := ws.NewHub(a.store, log.Logger)
	go hub.Run()

	// Initialize HTTP server
	a.server = apis.NewServer(a.store, hub, a.Config.Port, a.Config.JWTSecret, fileStorage)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Info().Msgf("Starting HTTP server on :%d", a.Config.Port)
		if err := a.server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or server error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Error().Err(err).Msg("server error")
		return err
	case sig := <-sigChan:
		log.Info().Msgf("received signal: %v", sig)
	}

	// Graceful shutdown
	log.Info().Msg("shutting down server...")
	if err := a.server.Shutdown(); err != nil {
		log.Error().Err(err).Msg("error during shutdown")
		return err
	}

	log.Info().Msg("server stopped gracefully")
	return nil
}
