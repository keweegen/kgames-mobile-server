package kapp

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	grpcbroadcaster "github.com/keweegen/tic-toe/internal/broadcaster/grpc"
	rediscache "github.com/keweegen/tic-toe/internal/cache/redis"
	"github.com/keweegen/tic-toe/internal/config"
	"github.com/keweegen/tic-toe/internal/logger"
	"github.com/keweegen/tic-toe/internal/server/grpc"
	"github.com/keweegen/tic-toe/internal/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log/slog"
	"os"
)

type setupHandler func(app *App) error

var setupHandlers = []setupHandler{
	setup.config,
	setup.logger,
	setup.database,
	setup.redis,
	setup.broadcaster,
	setup.repositories,
	setup.services,
	setup.background(setup.grpcServer),
}

var setup installer

type installer struct{}

func (installer) config(app *App) error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	app.cfg = cfg

	return nil
}

func (installer) logger(app *App) error {
	level := slog.LevelInfo

	if app.cfg.App.Debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	modifiedLogger := slog.New(handler).With(
		slog.String(logger.FieldAppEnvironment, app.cfg.App.Environment),
		slog.String(logger.FieldAppName, app.cfg.App.Name),
		slog.String(logger.FieldAppVersion, app.cfg.App.Version),
		slog.Bool(logger.FieldAppDebug, app.cfg.App.Debug),
	)

	slog.SetDefault(modifiedLogger)

	return nil
}

func (installer) database(app *App) error {
	conn, err := sql.Open("postgres", app.cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}
	if err = conn.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	app.db = conn
	appendShutdownHandler(app.db.Close)

	boil.SetDB(app.db)
	boil.DebugMode = app.cfg.App.Debug

	return nil
}

func (installer) redis(app *App) error {
	opts := &redis.Options{
		Addr: app.cfg.Redis.Address,
		DB:   app.cfg.Redis.DB,
	}

	app.redisClient = redis.NewClient(opts)
	appendShutdownHandler(app.redisClient.Close)

	return nil
}

func (installer) broadcaster(app *App) error {
	app.gameBroadcaster = grpcbroadcaster.NewBroadcaster("game", app.redisClient)
	appendShutdownHandler(app.gameBroadcaster.Shutdown)
	return nil

}

func (installer) repositories(app *App) error {
	app.repositories = store.NewRepository(app.db)
	return nil
}

func (installer) services(app *App) error {
	cacheClient := rediscache.NewCacheClient(app.cfg.App.Name, app.redisClient)
	app.services = store.NewService(app.gameBroadcaster, app.repositories, cacheClient)
	return nil
}

func (installer) background(h setupHandler) setupHandler {
	return func(app *App) error {
		go func() {
			if err := h(app); err != nil {
				slog.Error("background setup", slog.Any("error", err))
			}
		}()
		return nil
	}
}

func (installer) grpcServer(app *App) error {
	s := grpc.NewServer(app.gameBroadcaster, app.services)
	addr := fmt.Sprintf(":%d", app.cfg.GRPCServer.Port)

	appendShutdownHandler(s.Shutdown)
	if err := s.Run(addr); err != nil {
		return err
	}

	return nil
}
