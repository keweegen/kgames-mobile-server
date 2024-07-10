package game

import (
	"strings"
)

const (
	UnitCodeTicTacToeX UnitCode = "X"
	UnitCodeTicTacToeO UnitCode = "O"
)

type UnitCode string

func (c UnitCode) String() string {
	return string(c)
}

func (c UnitCode) Valid() bool {
	return strings.TrimSpace(c.String()) != ""
}
