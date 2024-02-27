package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene  int          = TITLE
	gameRectangle rl.Rectangle = rl.NewRectangle(775, 500, 175, 65)
	exitRectangle rl.Rectangle = rl.NewRectangle(1005, 500, 145, 65)
)

const (
	TITLE int = iota
	GAME
)

func drawScene() {
	switch currentScene {
	case TITLE:
		rl.BeginDrawing()
		rl.DrawTextEx(fontBold, "Odinbit", rl.NewVector2(775, 400), 72, 2, rl.White)
		rl.DrawTextEx(font, "Game", rl.NewVector2(775, 500), 56, 2, rl.White)
		rl.DrawTextEx(font, "Exit", rl.NewVector2(1005, 500), 56, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, gameRectangle) {
				currentScene = GAME
			}
			if rl.CheckCollisionPointRec(mousePos, exitRectangle) {
				rl.CloseWindow()
			}
		}
	case GAME:
		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.BeginMode2D(cam)
		drawWorld()
		rl.EndMode2D()
		rl.EndDrawing()
	}
}
