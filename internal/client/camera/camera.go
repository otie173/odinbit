package camera

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	Camera rl.Camera2D
)

func LoadCamera() {
	Camera = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), rl.NewVector2(256, 256), 0.0, 8.0)
}

func UpdateCamera() {
	Camera.Target.X = 256*12 + 12/2
	Camera.Target.Y = 256*12 + 12/2
}
