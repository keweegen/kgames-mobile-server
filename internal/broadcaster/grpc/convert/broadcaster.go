package convert

import (
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	gamedomain "github.com/keweegen/tic-toe/internal/domain/game"
)

var Broadcaster broadcasterConverter

type broadcasterConverter struct{}

var mappingDomainToPBAction = map[gamedomain.StreamAction]gamepb.StreamAction{
	gamedomain.StreamActionUnknown:          gamepb.StreamAction_SA_UNKNOWN,
	gamedomain.StreamActionPlayerJoin:       gamepb.StreamAction_SA_PLAYER_CONNECT,
	gamedomain.StreamActionPlayerReady:      gamepb.StreamAction_SA_PLAYER_READY,
	gamedomain.StreamActionPlayerMove:       gamepb.StreamAction_SA_PLAYER_MOVE,
	gamedomain.StreamActionPlayerTimeout:    gamepb.StreamAction_SA_PLAYER_TIMEOUT,
	gamedomain.StreamActionPlayerDraw:       gamepb.StreamAction_SA_PLAYER_DRAW,
	gamedomain.StreamActionPlayerGaveUp:     gamepb.StreamAction_SA_PLAYER_GAVE_UP,
	gamedomain.StreamActionPlayerDisconnect: gamepb.StreamAction_SA_PLAYER_DISCONNECT,
	gamedomain.StreamActionGameStart:        gamepb.StreamAction_SA_GAME_START,
	gamedomain.StreamActionGameFinish:       gamepb.StreamAction_SA_GAME_FINISH,
}

func (broadcasterConverter) Response(msg broadcaster.Message) *gamepb.StreamResponse {
	return &gamepb.StreamResponse{
		Ok:          msg.Error != nil,
		Description: msg.Description,
		GameId:      msg.GameID,
		UserId:      msg.UserID,
		Action:      mappingDomainToPBAction[msg.Action],
		Data:        nil, // TODO: convert by action
	}
}
