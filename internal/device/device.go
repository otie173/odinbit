package device

import (
	"odinbit/internal/entity"
	"odinbit/internal/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboard() {
	var playerSpeed float32 = 45.0
	delta := rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyW) {
		player.Roben.Pos.Y -= playerSpeed * delta
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.Roben.Dir = entity.LEFT
		player.Roben.Pos.X -= playerSpeed * delta
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.Roben.Pos.Y += playerSpeed * delta
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.Roben.Dir = entity.RIGHT
		player.Roben.Pos.X += playerSpeed * delta
	}
}

func HandleMouse() {

}
