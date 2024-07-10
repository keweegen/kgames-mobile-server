package store

import (
	"github.com/keweegen/tic-toe/internal/broadcaster"
	rediscache "github.com/keweegen/tic-toe/internal/cache/redis"
	gamedomain "github.com/keweegen/tic-toe/internal/domain/game"
	"github.com/keweegen/tic-toe/internal/domain/game/mover"
	gameservice "github.com/keweegen/tic-toe/internal/domain/game/service"
	userdomain "github.com/keweegen/tic-toe/internal/domain/user"
	userservice "github.com/keweegen/tic-toe/internal/domain/user/service"
	"github.com/shopspring/decimal"
	"time"
)

type Service struct {
	Game       gamedomain.Service
	GameStream gamedomain.StreamService
	User       userdomain.Service
}

func NewService(
	gameBroadcaster broadcaster.Broadcaster,
	repo *Repository,
	cacheClient rediscache.CacheClient,
) *Service {
	feePercent := decimal.NewFromFloat(1.5)
	gameTTL := time.Minute * 30

	gameService := gameservice.NewGame(
		feePercent,
		repo.Game,
		repo.GameMove,
		mover.New(repo.GameMove, feePercent),
		gameBroadcaster,
	)

	return &Service{
		Game:       gameService,
		GameStream: gameservice.NewGameStream(gameTTL, cacheClient, gameService),
		User:       userservice.NewUserService(repo.User),
	}
}
