package common

type Scene int

const (
	Title Scene = iota
	Connect
	Game
	ConnClosed
)

const (
	ScreenWidth      int32 = 1920
	ScreenHeight     int32 = 1080
	BaseRenderWidth  int32 = 1920
	BaseRenderHeight int32 = 1080
)
