package main

import (
	"context"
	"github.com/keweegen/tic-toe/internal/kapp"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "run", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "service is started")
	<-ctx.Done()
	slog.InfoContext(ctx, "service is shutting down")
}

func run(ctx context.Context) error {
	app, err := kapp.NewApp(ctx)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if err := app.Shutdown(); err != nil {
			slog.ErrorContext(ctx, "shutdown", slog.Any("error", err))
			return
		}
	}()

	return nil
}
