package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rLukoyanov/w/config"
	"github.com/rLukoyanov/w/core"
	"github.com/rLukoyanov/w/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "w"}
	var dbPath string
	var logLevel string

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return logger.Setup(logLevel)
	}

	serve := &cobra.Command{
		Use:   "serve",
		Short: "Запустить сервер",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msg("starting application")

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			cfg := &config.Config{DBPath: dbPath}

			app := core.New(cfg)
			err := app.Run()
			if err != nil {
				return err
			}

			<-ctx.Done()
			return nil
		},
	}

	serve.Flags().StringVar(&dbPath, "db", "data.db", "SQLite database file path")
	root.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level: debug|info|warn|error")

	root.AddCommand(serve)
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command failed")
	}
}
