package device

import (
	"odinbit/internal/entity"
	"odinbit/internal/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboard() {
	if rl.IsKeyDown(rl.KeyW) {
		player.Roben.Pos.Y -= 0.5
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.Roben.Dir = entity.LEFT
		player.Roben.Pos.X -= 0.5
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.Roben.Pos.Y += 0.5
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.Roben.Dir = entity.RIGHT
		player.Roben.Pos.X += 0.5
	}
}

func HandleMouse() {

}
