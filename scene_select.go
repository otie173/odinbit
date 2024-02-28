package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene int = TITLE
)

const (
	TITLE int = iota
	GAME
)

func drawScene() {
	switch currentScene {
	case TITLE:
		odinbitLabelSize := rl.MeasureTextEx(fontBold, "Odinbit", 72, 2)
		odinbitLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(odinbitLabelSize.X))/2, float32(rl.GetScreenHeight()-int(odinbitLabelSize.Y))/2-95)
		exitRectangle := rl.NewRectangle(odinbitLabelPos.X+230, odinbitLabelPos.Y+90, 145, 65)
		gameRectangle := rl.NewRectangle(odinbitLabelPos.X, odinbitLabelPos.Y+90, 175, 65)

		rl.BeginDrawing()
		rl.DrawTextEx(fontBold, "Odinbit", odinbitLabelPos, 72, 2, rl.White)
		rl.DrawTextEx(font, "Game", rl.NewVector2(odinbitLabelPos.X, odinbitLabelPos.Y+90), 56, 2, rl.White)
		rl.DrawTextEx(font, "Exit", rl.NewVector2(odinbitLabelPos.X+230, odinbitLabelPos.Y+90), 56, 2, rl.White)
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
