package game

type PlayerStateCode string

const (
	PlayerStatusCodeUnknown     PlayerStateCode = ""
	PlayerStateCodeJoined       PlayerStateCode = "CONNECTED"
	PlayerStateCodeReady        PlayerStateCode = "READY"
	PlayerStateCodePlaying      PlayerStateCode = "PLAYING"
	PlayerStateCodeDisconnected PlayerStateCode = "DISCONNECTED"
)

func (p PlayerStateCode) Unknown() bool {
	return p == PlayerStatusCodeUnknown
}

func (p PlayerStateCode) Connected() bool {
	return p == PlayerStateCodeJoined
}

func (p PlayerStateCode) Ready() bool {
	return p == PlayerStateCodeReady
}

func (p PlayerStateCode) Playing() bool {
	return p == PlayerStateCodePlaying
}

func (p PlayerStateCode) Disconnected() bool {
	return p == PlayerStateCodeDisconnected
}
