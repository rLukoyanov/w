package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rLukoyanov/w/core"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "w"}

	serve := &cobra.Command{
		Use:   "serve",
		Short: "Запустить сервер",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("W запущен на :8090")

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			app := core.New(nil)
			err := app.Run()
			if err != nil {
				return err
			}

			<-ctx.Done()
			return nil
		},
	}

	root.AddCommand(serve)
	root.Execute()
}
