package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene Scene = TITLE
	lastScene    Scene = -1
	menuOpen     bool
)

const (
	TITLE Scene = iota
	MODE
	GENERATE
	SAVE
	GAME
	INVENTORY
	MENU
)

type Scene int

func drawScene() {
	switch currentScene {
	case TITLE:
		odinbitLabelSize := rl.MeasureTextEx(fontBold, "Odinbit", 72, 2)
		odinbitLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(odinbitLabelSize.X))/2, float32(rl.GetScreenHeight()-int(odinbitLabelSize.Y))/2-95)
		exitRectangle := rl.NewRectangle(odinbitLabelPos.X+230, odinbitLabelPos.Y+90, 145, 65)
		playRectangle := rl.NewRectangle(odinbitLabelPos.X, odinbitLabelPos.Y+90, 175, 65)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, "Odinbit", odinbitLabelPos, 72, 2, rl.White)
		rl.DrawTextEx(font, "Play", rl.NewVector2(odinbitLabelPos.X, odinbitLabelPos.Y+90), 56, 2, rl.White)
		rl.DrawTextEx(font, "Exit", rl.NewVector2(odinbitLabelPos.X+230, odinbitLabelPos.Y+90), 56, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, playRectangle) {
				currentScene = MODE
			}
			if rl.CheckCollisionPointRec(mousePos, exitRectangle) {
				rl.CloseWindow()
			}
		}
	case MODE:
		modeLabelText := "Select mode"
		modeLabelSize := rl.MeasureTextEx(fontBold, modeLabelText, 72, 2)
		modeLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(modeLabelSize.X))/2, float32(rl.GetScreenHeight()-int(modeLabelSize.Y))/2-175)

		singleplayerText := "1) Singleplayer"
		multiplayerText := "2) Multiplayer"

		singleplayerLabelSize := rl.MeasureTextEx(font, singleplayerText, 48, 2)
		multiplayerLabelSize := rl.MeasureTextEx(font, multiplayerText, 48, 2)

		// Находим максимальную ширину текста пунктов меню
		maxWidth := singleplayerLabelSize.X
		if multiplayerLabelSize.X > maxWidth {
			maxWidth = multiplayerLabelSize.X
		}

		// Вычисляем позицию для текста, центрируя относительно "Select mode"
		textX := modeLabelPos.X + (modeLabelSize.X-maxWidth)/2

		singleplayerLabelPos := rl.NewVector2(textX, float32(rl.GetScreenHeight()-int(singleplayerLabelSize.Y))/2-80)
		multiplayerLabelPos := rl.NewVector2(textX, float32(rl.GetScreenHeight()-int(multiplayerLabelSize.Y))/2)

		singleplayerRectangle := rl.NewRectangle(singleplayerLabelPos.X, singleplayerLabelPos.Y, 475, 65)
		multiplayerRectangle := rl.NewRectangle(multiplayerLabelPos.X, multiplayerLabelPos.Y, 430, 65)
		//exitRectangle := rl.NewRectangle(odinbitLabelPos.X+230, odinbitLabelPos.Y+90, 145, 65)
		//playRectangle := rl.NewRectangle(odinbitLabelPos.X, odinbitLabelPos.Y+90, 175, 65)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawRectangleRec(singleplayerRectangle, rl.Red)
		rl.DrawRectangleRec(multiplayerRectangle, rl.Red)
		rl.DrawTextEx(fontBold, modeLabelText, modeLabelPos, 72, 2, rl.White)
		rl.DrawTextEx(font, singleplayerText, singleplayerLabelPos, 48, 2, rl.White)
		rl.DrawTextEx(font, multiplayerText, multiplayerLabelPos, 48, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, singleplayerRectangle) {
				if checkWorldFile() {
					loadWorldLabelSize := rl.MeasureTextEx(font, "Load world...", 56, 2)
					loadWorldLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(loadWorldLabelSize.X))/2, float32(rl.GetScreenHeight()-int(loadWorldLabelSize.Y))/2)

					rl.BeginDrawing()
					rl.ClearBackground(bkgColor)
					rl.DrawTextEx(font, "Load world...", loadWorldLabelPos, 56, 2, rl.White)
					rl.EndDrawing()
					world = loadWorldFile()
					worldInfo = loadWorldInfo()
					updateWorld()
					loadPlayerFile()
					currentScene = GAME
				} else {
					currentScene = GENERATE
				}
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
	case SAVE:
		saveWorldLabelSize := rl.MeasureTextEx(font, "Save world...", 56, 2)
		saveWorldLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(saveWorldLabelSize.X))/2, float32(rl.GetScreenHeight()-int(saveWorldLabelSize.Y))/2)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(font, "Save world...", saveWorldLabelPos, 56, 2, rl.White)
		rl.EndDrawing()

		saveWorldFile()
		saveWorldInfo()
		savePlayerFile()
		currentScene = TITLE
	case GAME:
		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.BeginMode2D(cam)
		drawWorld(cam)
		drawPlayer()
		rl.EndMode2D()
		rl.DrawRectangleV(rl.NewVector2(0, 0), rl.NewVector2(240, 130), rl.ColorAlpha(rl.Black, 0.65))
		rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, 240, 130), 5, rl.White)
		drawUI()
		rl.DrawTextEx(font, fmt.Sprintf("X: %d Y: %d", int(targetPosition.X)/10, int(targetPosition.Y)/10), rl.NewVector2(50, 100), 16, 2, rl.White)
		//rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 5, 180, 24, rl.White)
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
		rightArrowRec := rl.NewRectangle(rightArrowPosition.X+15, rightArrowPosition.Y+5, 85, 90)
		leftArrowRec := rl.NewRectangle(leftArrowPosition.X-15, leftArrowPosition.Y+5, 85, 90)
		//rl.DrawRectangle(int32(rightArrowRec.X), int32(rightArrowRec.Y), int32(rightArrowRec.Width), int32(rightArrowRec.Height), rl.Red)
		//rl.DrawRectangle(int32(leftArrowRec.X), int32(leftArrowRec.Y), int32(leftArrowRec.Width), int32(leftArrowRec.Height), rl.Red)
		rl.DrawTextureEx(leftArrow, leftArrowPosition, 0, arrowScale, rl.White)
		rl.DrawTextureEx(rightArrow, rightArrowPosition, 0, arrowScale, rl.White)
		pageLabelSize := rl.MeasureTextEx(font, fmt.Sprintf("Page %d/%d", currentPage, maxPage), 24, 2)
		pageLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(pageLabelSize.X))/2, 690.0)
		rl.DrawTextEx(font, fmt.Sprintf("Page %d/%d", currentPage, maxPage), pageLabelPos, 24, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, rightArrowRec) {
				if currentPage+1 <= maxPage {
					currentPage++
				}
			}
			if rl.CheckCollisionPointRec(mousePos, leftArrowRec) {
				if currentPage-1 >= 1 {
					currentPage--
				}
			}
		}
	case MENU:
		menuLabelSize := rl.MeasureTextEx(fontBold, "Game menu", 72, 2)
		menuLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(menuLabelSize.X))/2, float32(rl.GetScreenHeight()-int(menuLabelSize.Y))/2-175)
		backToGameLabelSize := rl.MeasureTextEx(font, "1) Back to game", 48, 2)
		backToGameLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(backToGameLabelSize.X))/2, float32(rl.GetScreenHeight()-int(backToGameLabelSize.Y))/2-80)
		saveAndQuitLabelSize := rl.MeasureTextEx(font, "2) Save and quit", 48, 2)
		saveAndQuitLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(backToGameLabelSize.X))/2, float32(rl.GetScreenHeight()-int(saveAndQuitLabelSize.Y))/2)
		backToGameRectangle := rl.NewRectangle(backToGameLabelPos.X, backToGameLabelPos.Y, 515, 65)
		saveAndQuitRectangle := rl.NewRectangle(backToGameLabelPos.X, saveAndQuitLabelPos.Y, 530, 65)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, "Game menu", rl.NewVector2(menuLabelPos.X, menuLabelPos.Y), 72, 2, rl.White)
		rl.DrawTextEx(font, "1) Back to game", rl.NewVector2(backToGameLabelPos.X, backToGameLabelPos.Y), 48, 2, rl.White)
		rl.DrawTextEx(font, "2) Save and quit", rl.NewVector2(saveAndQuitLabelPos.X, saveAndQuitLabelPos.Y), 48, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, backToGameRectangle) {
				currentScene = GAME
				menuOpen = false
			}
			if rl.CheckCollisionPointRec(mousePos, saveAndQuitRectangle) {
				currentScene = SAVE
				menuOpen = false
			}
		}
	}
}
