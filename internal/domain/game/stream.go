package game

import "time"

type StreamDataGameStart struct {
	StrikerID   string
	TimeForMove time.Duration
}

type StreamDataPlayerMove struct {
	StrikerID string
	UnitCode  UnitCode
}
