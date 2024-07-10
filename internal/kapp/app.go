package kapp

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	"github.com/keweegen/tic-toe/internal/config"
	"github.com/keweegen/tic-toe/internal/store"
)

type App struct {
	ctx         context.Context
	cfg         *config.Config
	db          *sql.DB
	redisClient *redis.Client

	gameBroadcaster broadcaster.Broadcaster
	repositories    *store.Repository
	services        *store.Service
}

func NewApp(ctx context.Context) (*App, error) {
	app := new(App)
	app.ctx = ctx

	if err := app.setup(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) setup() error {
	for _, h := range setupHandlers {
		if err := h(a); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Shutdown() error {
	for i := len(shutdownHandlers) - 1; i >= 0; i-- {
		if err := shutdownHandlers[i](); err != nil {
			return err
		}
	}

	return nil
}
