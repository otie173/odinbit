package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	bkgColor       rl.Color = rl.NewColor(0, 0, 0, 255)
	fontBold, font rl.Font
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
	drawScene()
}

func exit() {
	rl.CloseWindow()

	unloadWorld()
	unloadPlayer()
	rl.UnloadFont(fontBold)
	rl.UnloadFont(font)
}

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Odinbit")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	fontBold = rl.LoadFont("assets/fonts/pypx/pypx_bold.ttf")
	font = rl.LoadFont("assets/fonts/pypx/pypx.ttf")
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
