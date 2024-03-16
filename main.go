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
	cam.Target = rl.Vector2Lerp(cam.Target, rl.NewVector2(playerPosition.X+playerRectangle.Width/2, playerPosition.Y+playerRectangle.Height/2), 0.05)
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
	unloadInventory()
}

func init() {
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), "Odinbit")
	rl.SetExitKey(0)
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	fontBold = rl.LoadFont("assets/fonts/pypx/pypx_bold.ttf")
	font = rl.LoadFont("assets/fonts/pypx/pypx.ttf")
	loadWorld()
	loadPlayer()
	loadAudio()
	loadInventory()
	generateWorld()
}

func main() {
	for !rl.WindowShouldClose() {
		update()
		render()
	}
	exit()
}
