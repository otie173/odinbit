package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	bkgColor rl.Color = rl.NewColor(0, 0, 0, 255)
)

const (
	SCREEN_WIDTH  int32 = 1920
	SCREEN_HEIGHT int32 = 1080
)

func update() {
	keyboardHandler()
	mouseHandler()
	cam.Target = rl.NewVector2(playerPosition.X, playerPosition.Y)
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)
	drawWorld()
	rl.EndMode2D()
	rl.EndDrawing()
}

func exit() {
	rl.CloseWindow()

	unloadWorld()
	unloadPlayer()
}

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Odinbit")
	rl.SetExitKey(0)
	rl.SetTargetFPS(144)

	loadWorld()
	loadPlayer()
}

func main() {
	for !rl.WindowShouldClose() {
		update()
		render()
	}
	exit()
}
