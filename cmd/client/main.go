package main

import (
	"odinbit/internal/device"
	"odinbit/internal/player"
	"odinbit/internal/scene"
	"odinbit/internal/world"
	"odinbit/utils/build"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func load() {
	build.SetBuildType(build.Debug)

	world.LoadTexture()

	player.LoadTexture()
	player.RegisterCam()

	world.AddBlock(1, 2, world.Tree1)
}

func update() {
	player.UpdateCamera()

	device.HandleMouse()
	device.HandleKeyboard()
}

func main() {
	screenWidth, screenHeight := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())

	rl.SetConfigFlags(rl.FlagFullscreenMode | rl.FlagWindowUnfocused | rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Odinbit")

	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	load()
	for !rl.WindowShouldClose() {
		go update()
		scene.DrawScene()
	}
	rl.CloseWindow()
}
