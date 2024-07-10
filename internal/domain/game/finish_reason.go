package game

type FinishReasonCode string

const (
	FinishReasonCodeUnknown FinishReasonCode = ""
	FinishReasonCodeDefault FinishReasonCode = "DEFAULT"
	FinishReasonCodeDraw    FinishReasonCode = "DRAW"
	FinishReasonCodeGaveUp  FinishReasonCode = "GAVE_UP"
	FinishReasonCodeInvalid FinishReasonCode = "INVALID"
)

func (c FinishReasonCode) String() string {
	return string(c)
}

func (c FinishReasonCode) Unknown() bool {
	return c == FinishReasonCodeUnknown
}
