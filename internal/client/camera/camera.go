package camera

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/player"
)

var (
	Camera rl.Camera2D
)

func LoadCamera() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	baseWidth := float32(common.BaseRenderWidth)
	baseHeight := float32(common.BaseRenderHeight)

	var scale float32 = 1.0
	if baseWidth > 0 && baseHeight > 0 {
		widthScale := screenWidth / baseWidth
		heightScale := screenHeight / baseHeight
		scale = float32(math.Min(float64(widthScale), float64(heightScale)))
		if scale <= 0 {
			scale = 1
		}
	}

	const baseZoom float32 = 8.0
	zoom := baseZoom * scale
	offset := rl.NewVector2(screenWidth/2, screenHeight/2)

	Camera = rl.NewCamera2D(offset, rl.NewVector2(256, 256), 0.0, zoom)
}

func UpdateCamera() {
	player.PlayerMu.Lock()
	targetX := player.GamePlayer.CurrentX*12 + 12/2
	targetY := player.GamePlayer.CurrentY*12 + 12/2
	player.PlayerMu.Unlock()

	xVec := rl.NewVector2(targetX, Camera.Target.Y)
	yVec := rl.NewVector2(Camera.Target.X, targetY)

	duration := 0.25
	step := rl.GetFrameTime() / float32(duration)

	newX := rl.Vector2Lerp(Camera.Target, xVec, step).X
	newY := rl.Vector2Lerp(Camera.Target, yVec, step).Y

	Camera.Target.X = newX
	Camera.Target.Y = newY
}
