package main

import (
	"odinbit/utils/build"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
	build.SetBuildType(build.Debug)
}

func update() {

}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.EndDrawing()
}

func main() {
	screenWidth, screenHeight := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())

	rl.SetConfigFlags(rl.FlagFullscreenMode | rl.FlagVsyncHint | rl.FlagWindowUnfocused)
	rl.InitWindow(screenWidth, screenHeight, "Odinbit")

	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		go update()
		draw()
	}
	rl.CloseWindow()
}
