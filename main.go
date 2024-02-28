package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	bkgColor       rl.Color = rl.NewColor(0, 0, 0, 255)
	fontBold, font rl.Font
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
	unloadAudio()
	rl.UnloadFont(fontBold)
	rl.UnloadFont(font)
}

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), "Odinbit")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	fontBold = rl.LoadFont("assets/fonts/pypx/pypx_bold.ttf")
	font = rl.LoadFont("assets/fonts/pypx/pypx.ttf")
	loadWorld()
	loadPlayer()
	loadAudio()
}

func main() {
	for !rl.WindowShouldClose() {
		update()
		render()
	}
	exit()
}
