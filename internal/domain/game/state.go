package game

type StateCode string

const (
	StateCodeUnknown  StateCode = ""
	StateCodeCreated  StateCode = "CREATED"
	StateCodeActive   StateCode = "ACTIVE"
	StateCodeFinished StateCode = "FINISHED"
)

func (s StateCode) String() string {
	return string(s)
}
