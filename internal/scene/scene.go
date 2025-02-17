package scene

import (
	"odinbit/internal/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene Scene = GAME
	lastScene    Scene = -1
)

const (
	TITLE Scene = iota
	MODE
	GENERATE
	SAVE
	LOAD
	GAME
	MENU
)

type Scene int

func DrawScene() {
	switch currentScene {
	case GAME:
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(player.Cam)
		player.DrawPlayer()
		rl.EndMode2D()
		rl.EndDrawing()
	}
}
