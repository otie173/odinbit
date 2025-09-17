package camera

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/player"
)

var (
	Camera rl.Camera2D
)

func LoadCamera() {
	Camera = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), rl.NewVector2(256, 256), 0.0, 8.0)
}

func UpdateCamera() {
	targetX := player.GamePlayer.CurrentX*12 + 12/2
	targetY := player.GamePlayer.CurrentY*12 + 12/2

	xVec := rl.NewVector2(targetX, Camera.Target.Y)
	yVec := rl.NewVector2(Camera.Target.X, targetY)

	duration := 0.5
	step := rl.GetFrameTime() / float32(duration)

	newX := rl.Vector2Lerp(Camera.Target, xVec, step).X
	newY := rl.Vector2Lerp(Camera.Target, yVec, step).Y

	Camera.Target.X = newX
	Camera.Target.Y = newY
}
