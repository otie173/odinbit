package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene int = TITLE
)

const (
	TITLE int = iota
	GENERATE
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
				currentScene = GENERATE
			}
			if rl.CheckCollisionPointRec(mousePos, exitRectangle) {
				rl.CloseWindow()
			}
		}
	case GENERATE:
		generatingWorldLabelSize := rl.MeasureTextEx(font, "Generating world...", 56, 2)
		generatingWorldLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(generatingWorldLabelSize.X))/2, float32(rl.GetScreenHeight()-int(generatingWorldLabelSize.Y))/2)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(font, "Generating world...", generatingWorldLabelPos, 56, 2, rl.White)
		rl.EndDrawing()

		generateWorld()
		if worldGenerated {
			currentScene = GAME
		}

	case GAME:
		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 0, 0, 32, rl.White)
		rl.BeginMode2D(cam)
		drawWorld()
		drawPlayer()
		rl.EndMode2D()
		rl.EndDrawing()
	case INVENTORY:
		inventoryLabelSize := rl.MeasureTextEx(fontBold, "Resources", 72, 2)
		inventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(inventoryLabelSize.X))/2, 75)
		hotInventoryLabelSize := rl.MeasureTextEx(fontBold, "Items & Blocks", 72, 2)
		hotInventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(hotInventoryLabelSize.X))/2, 290)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, "Resources", rl.NewVector2(inventoryLabelPos.X, inventoryLabelPos.Y), 72, 2, rl.White)
		rl.DrawTextEx(fontBold, "Items & Blocks", rl.NewVector2(hotInventoryLabelPos.X, hotInventoryLabelPos.Y), 72, 2, rl.White)
		for i := range inventory {
			rl.DrawTextureEx(slotImage, rl.NewVector2(inventory[i].x, inventory[i].y), 0, inventoryZoom, rl.White)
		}
		for i := range hotInventory {
			rl.DrawTextureEx(slotImage, rl.NewVector2(hotInventory[i].x, hotInventory[i].y), 0, inventoryZoom, rl.White)

		}
		drawItems()
		rl.EndDrawing()
	}
}
