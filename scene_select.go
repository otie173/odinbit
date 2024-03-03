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
	INVENTORY
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
		drawPlayer()
		rl.EndMode2D()
		rl.EndDrawing()
	case INVENTORY:
		inventoryLabelSize := rl.MeasureTextEx(fontBold, "Inventory", 72, 2)
		inventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(inventoryLabelSize.X))/2, 75)
		hotInventoryLabelSize := rl.MeasureTextEx(fontBold, "Hot inventory", 72, 2)
		hotInventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(hotInventoryLabelSize.X))/2, 650)
		slotWidth := float32(slotImage.Width * int32(cam.Zoom))
		slotSpacing := float32(32)
		startX1 := inventoryLabelPos.X - 36
		startX2 := startX1

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, "Inventory", rl.NewVector2(inventoryLabelPos.X, inventoryLabelPos.Y), 72, 2, rl.White)
		for i := 0; i < 6; i++ {
			slotX := startX1 + float32(i)*(slotWidth+slotSpacing)
			rl.DrawTextureEx(slotImage, rl.NewVector2(slotX, 200), 0, cam.Zoom, rl.White)
			rl.DrawTextureEx(slotImage, rl.NewVector2(slotX, 300), 0, cam.Zoom, rl.White)
			rl.DrawTextureEx(slotImage, rl.NewVector2(slotX, 400), 0, cam.Zoom, rl.White)
			rl.DrawTextureEx(slotImage, rl.NewVector2(slotX, 500), 0, cam.Zoom, rl.White)
		}
		rl.DrawTextEx(fontBold, "Hot inventory", rl.NewVector2(hotInventoryLabelPos.X, hotInventoryLabelPos.Y), 72, 2, rl.White)
		for i := 0; i < 6; i++ {
			slotX := startX2 + float32(i)*(slotWidth+slotSpacing)
			rl.DrawTextureEx(slotImage, rl.NewVector2(slotX, 775), 0, cam.Zoom, rl.White)
		}
		rl.EndDrawing()
	}
}
