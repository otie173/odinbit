package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/camera"
	"github.com/otie173/odinbit/internal/device"
	"github.com/otie173/odinbit/internal/player"
	"github.com/otie173/odinbit/internal/world"
	"github.com/otie173/odinbit/resource"
)

func New() {
	rl.SetConfigFlags(rl.FlagFullscreenMode | rl.FlagVsyncHint | rl.FlagWindowUnfocused)
	rl.InitWindow(1920, 1080, "Odinbit")
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)
}

func Load() {
	resource.Load()
	player.Load()
	world.Load()
	camera.Load()
}

func update() {
	device.HandleMouse()
	device.HandleKeyboard()
	camera.Update()
}

func render() {
	rl.BeginDrawing()
	rl.BeginMode2D(camera.Cam)
	rl.ClearBackground(rl.Black)
	world.Draw()
	player.Draw()
	rl.EndMode2D()
	rl.EndDrawing()
}

func Run() {
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		update()
		render()
	}
}
