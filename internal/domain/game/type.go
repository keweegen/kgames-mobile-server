package game

type TypeCode string

const (
	TypeCodeUnknown   TypeCode = ""
	TypeCodeTicTacToe TypeCode = "TIC-TAC-TOE"
)

func (t TypeCode) String() string {
	return string(t)
}
