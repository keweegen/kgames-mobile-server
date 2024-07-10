package game

type StreamAction string

const (
	StreamActionUnknown          StreamAction = "unknown"
	StreamActionPlayerJoin       StreamAction = "player.join"
	StreamActionPlayerReady      StreamAction = "player.ready"
	StreamActionPlayerMove       StreamAction = "player.move"
	StreamActionPlayerTimeout    StreamAction = "player.timeout"
	StreamActionPlayerDraw       StreamAction = "player.draw"
	StreamActionPlayerGaveUp     StreamAction = "player.gave_up"
	StreamActionPlayerDisconnect StreamAction = "player.disconnect"
	StreamActionGameStart        StreamAction = "game.start"
	StreamActionGameFinish       StreamAction = "game.finish"
)

func (s StreamAction) String() string {
	return string(s)
}
