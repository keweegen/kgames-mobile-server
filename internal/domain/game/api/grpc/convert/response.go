package convert

import (
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var Response response

type response struct{}

var mappingDomainToPBType = map[domain.TypeCode]gamepb.Type{
	domain.TypeCodeUnknown:   gamepb.Type_T_UNKNOWN,
	domain.TypeCodeTicTacToe: gamepb.Type_T_TIC_TAC_TOE,
}

var mappingDomainToPBState = map[domain.StateCode]gamepb.State{
	domain.StateCodeUnknown:  gamepb.State_S_UNKNOWN,
	domain.StateCodeCreated:  gamepb.State_S_CREATED,
	domain.StateCodeActive:   gamepb.State_S_ACTIVE,
	domain.StateCodeFinished: gamepb.State_S_FINISHED,
}

var mappingDomainToPBFinishReason = map[domain.FinishReasonCode]gamepb.FinishReason{
	domain.FinishReasonCodeUnknown: gamepb.FinishReason_FR_UNKNOWN,
	domain.FinishReasonCodeDefault: gamepb.FinishReason_FR_UNKNOWN,
	domain.FinishReasonCodeDraw:    gamepb.FinishReason_FR_DRAW,
	domain.FinishReasonCodeGaveUp:  gamepb.FinishReason_FR_GAVE_UP,
	domain.FinishReasonCodeInvalid: gamepb.FinishReason_FR_INVALID,
}

func (response) Game(g domain.Game) *gamepb.GameResponse {
	return &gamepb.GameResponse{
		Id:           g.ID,
		Type:         mappingDomainToPBType[g.TypeCode],
		StateCode:    mappingDomainToPBState[g.StateCode],
		Bid:          g.Bid.InexactFloat64(),
		MaxPlayers:   int32(g.MaxPlayers),
		FinishReason: mappingDomainToPBFinishReason[g.FinishReasonCode],
		Players:      Response.gamePlayerSlice(g.Players),
		CreatedAt:    timestamppb.New(g.CreatedAt),
		StartedAt:    timestamppb.New(g.StartedAt.Time),
		FinishedAt:   timestamppb.New(g.FinishedAt.Time),
	}
}

func (response) gamePlayer(p domain.Player) *gamepb.GamePlayerResponse {
	return &gamepb.GamePlayerResponse{
		Id:       p.ID,
		Position: int32(p.Position),
	}
}

func (response) gamePlayerSlice(ps []domain.Player) []*gamepb.GamePlayerResponse {
	res := make([]*gamepb.GamePlayerResponse, len(ps))
	for i, p := range ps {
		res[i] = Response.gamePlayer(p)
	}
	return res
}

func (response) Error(err error) *gamepb.StreamResponse {
	return &gamepb.StreamResponse{Ok: err != nil}
}
