package mover

import (
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/game/errors"
	"github.com/shopspring/decimal"
)

type ticTacToeMover struct {
	unitCode   domain.UnitCode
	feePercent decimal.Decimal
}

var _ mover = (*ticTacToeMover)(nil)

func newTicTacToeMover(unitCode domain.UnitCode, feePercent decimal.Decimal) ticTacToeMover {
	return ticTacToeMover{
		unitCode:   unitCode,
		feePercent: feePercent,
	}
}

func (m ticTacToeMover) do(moves []domain.Move, newMove *domain.Move, players domain.Players) (finished bool, err error) {
	if len(moves) == 0 {
		return false, nil
	}
	if moves[len(moves)-1].UnitCode == newMove.UnitCode {
		return false, domainerrs.ErrNotYourTurn
	}
	if newMove.UnitCode != m.unitCode {
		return false, domainerrs.ErrMoveNotFound
	}

	if len(moves) >= 5 {
		return true, nil
	}

	grid := m.grid(moves)
	if grid[newMove.Position.X][newMove.Position.Y].Valid() {
		return false, domainerrs.ErrMoveExists
	}

	grid[newMove.Position.X][newMove.Position.Y] = newMove.UnitCode

	return false, nil
}

func (m ticTacToeMover) setBatterID(players domain.Players, newMove *domain.Move) {
	sortedPlayers := players.SortByPositions()

	for _, player := range sortedPlayers {
		if player.ID == newMove.StrikerID {
			nextPosition := player.Position + 1

			if nextPosition >= len(sortedPlayers) {
				nextPosition = 0
			}

			newMove.BatterID = sortedPlayers[nextPosition].ID

			break
		}
	}
}

func (m ticTacToeMover) grid(moves []domain.Move) [][]domain.UnitCode {
	grid := make([][]domain.UnitCode, 3)

	for i := 0; i < 3; i++ {
		grid[i] = make([]domain.UnitCode, 3)
	}

	for _, move := range moves {
		grid[move.Position.X][move.Position.Y] = move.UnitCode
	}

	return grid
}

func (m ticTacToeMover) results(game domain.Game) (domain.FinishParams, error) {
	grid := m.grid(game.Moves)
	lastMove := game.Moves[len(game.Moves)-1]
	playersCount := game.Players.LenDecimal()
	totalBid := game.Bid.Mul(playersCount)

	if m.draw(grid) {
		return m.drawResults(game, totalBid), nil
	}

	if m.winner(grid) {
		return m.winnerResults(game, lastMove, totalBid), nil
	}

	return domain.FinishParams{}, domainerrs.ErrGameInternal
}

func (m ticTacToeMover) drawResults(
	game domain.Game,
	totalBid decimal.Decimal,
) domain.FinishParams {
	fee := totalBid.Mul(m.feePercent)
	profit := totalBid.Sub(fee).Div(game.Players.LenDecimal())

	playersParams := make([]domain.PlayerFinishParams, len(game.Players))
	for i, player := range game.Players {
		playersParams[i] = domain.PlayerFinishParams{
			PlayerID: player.ID,
			Profit:   profit,
			Fee:      fee,
		}
	}

	return domain.FinishParams{
		ReasonCode: domain.FinishReasonCodeDraw,
		Players:    playersParams,
	}
}

func (m ticTacToeMover) winnerResults(
	game domain.Game,
	lastMove domain.Move,
	totalBid decimal.Decimal,
) domain.FinishParams {
	fee := totalBid.Mul(m.feePercent)
	winnerProfit := totalBid.Sub(fee)

	playerIDs := game.Players.IDs()
	playersParams := make([]domain.PlayerFinishParams, len(playerIDs))

	for i, playerID := range playerIDs {
		if playerID == lastMove.StrikerID {
			playersParams[i] = domain.PlayerFinishParams{
				PlayerID: playerID,
				Profit:   winnerProfit,
				Fee:      fee,
			}
			continue
		}

		playersParams[i] = domain.PlayerFinishParams{
			PlayerID: playerID,
			Profit:   game.Bid.Neg(),
			Fee:      decimal.Zero,
		}
	}

	return domain.FinishParams{
		ReasonCode: domain.FinishReasonCodeDefault,
		Players:    playersParams,
	}
}

func (m ticTacToeMover) winner(grid [][]domain.UnitCode) (winner bool) {
	for i := 0; i < 3; i++ {
		if grid[i][0] == m.unitCode && grid[i][1] == m.unitCode && grid[i][2] == m.unitCode {
			return true
		}
		if grid[0][i] == m.unitCode && grid[1][i] == m.unitCode && grid[2][i] == m.unitCode {
			return true
		}
	}

	if grid[0][0] == m.unitCode && grid[1][1] == m.unitCode && grid[2][2] == m.unitCode {
		return true
	}

	if grid[0][2] == m.unitCode && grid[1][1] == m.unitCode && grid[2][0] == m.unitCode {
		return true
	}

	return false
}

func (m ticTacToeMover) draw(grid [][]domain.UnitCode) (draw bool) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[i][j] != domain.UnitCodeTicTacToeX && grid[i][j] != domain.UnitCodeTicTacToeO {
				return false
			}
		}
	}

	return true
}
