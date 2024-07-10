package convert

import (
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	"github.com/shopspring/decimal"
)

var Request request

type request struct{}

var mappingPBToDomainType = map[gamepb.Type]domain.TypeCode{
	gamepb.Type_T_UNKNOWN:     domain.TypeCodeUnknown,
	gamepb.Type_T_TIC_TAC_TOE: domain.TypeCodeTicTacToe,
}

var mappingPBToDomainAction = map[gamepb.StreamAction]domain.StreamAction{
	gamepb.StreamAction_SA_UNKNOWN:           domain.StreamActionUnknown,
	gamepb.StreamAction_SA_PLAYER_CONNECT:    domain.StreamActionPlayerJoin,
	gamepb.StreamAction_SA_PLAYER_READY:      domain.StreamActionPlayerReady,
	gamepb.StreamAction_SA_PLAYER_MOVE:       domain.StreamActionPlayerMove,
	gamepb.StreamAction_SA_PLAYER_TIMEOUT:    domain.StreamActionPlayerTimeout,
	gamepb.StreamAction_SA_PLAYER_DRAW:       domain.StreamActionPlayerDraw,
	gamepb.StreamAction_SA_PLAYER_GAVE_UP:    domain.StreamActionPlayerGaveUp,
	gamepb.StreamAction_SA_PLAYER_DISCONNECT: domain.StreamActionPlayerDisconnect,
	gamepb.StreamAction_SA_GAME_START:        domain.StreamActionGameStart,
	gamepb.StreamAction_SA_GAME_FINISH:       domain.StreamActionGameFinish,
}

func (request) CreateGame(r *gamepb.CreateGameRequest) domain.CreateParams {
	return domain.CreateParams{
		InitiatorID: r.InitiatorId,
		TypeCode:    mappingPBToDomainType[r.Type],
		Bid:         decimal.NewFromFloat(r.Bid),
		MaxPlayers:  int(r.MaxPlayers),
	}
}

func (request) Stream(r *gamepb.StreamRequest) domain.StreamMessage {
	return domain.StreamMessage{
		GameID: r.GameId,
		UserID: r.UserId,
		Action: mappingPBToDomainAction[r.Action],
		Data:   r.Data.Value,
	}
}
