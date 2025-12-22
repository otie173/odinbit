package common

type Scene int
type Material int

var (
	MaterialMap = map[Material]string{
		-1:    "None",
		Wood:  "Wood",
		Stone: "Stone",
		Metal: "Metal",
	}
)

const (
	Title Scene = iota
	Connect
	Game
	ConnClosed
)

const (
	Wood Material = iota
	Stone
	Metal
)

const (
	TileSize         float32 = 12
	WorldSize        int     = 1024
	ScreenWidth      int32   = 1920
	ScreenHeight     int32   = 1080
	BaseRenderWidth  int32   = 1920
	BaseRenderHeight int32   = 1080
)
