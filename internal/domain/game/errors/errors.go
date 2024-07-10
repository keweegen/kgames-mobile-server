package errors

import "errors"

var (
	ErrInvalidStreamData = errors.New("invalid stream data")

	ErrGameNotFound            = errors.New("game not found")
	ErrGameIsFinished          = errors.New("game is finished")
	ErrGameIsNotActive         = errors.New("game is not active")
	ErrGameInternal            = errors.New("game internal error")
	ErrGameHasMaxPlayers       = errors.New("game has max players")
	ErrGameHasNotEnoughPlayers = errors.New("game has not enough players")

	ErrPlayerNotFound          = errors.New("player not found")
	ErrPlayerHasIncorrectState = errors.New("player has incorrect state")

	ErrNotYourTurn  = errors.New("not your turn")
	ErrMoveNotFound = errors.New("move not found")
	ErrMoveExists   = errors.New("move already exists")
)
