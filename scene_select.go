package main

import (
	"fmt"
	"strings"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	currentScene            Scene = TITLE
	lastScene               Scene = -1
	menuOpen                bool
	ipAddress               string
	backspaceHeldTime       float32
	backspaceRepeatDelay    float32 = 0.25
	backspaceRepeatInterval float32 = 0.05
)

const (
	TITLE Scene = iota
	MODE
	GENERATE
	SAVE
	GAME
	INVENTORY
	MENU
	IP_INPUT
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
				if !checkWorldFile("world.odn") || !checkWorldFile("world_send.odn") {
					currentScene = GENERATE
				} else {
					currentScene = MODE
				}
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

		maxWidth := singleplayerLabelSize.X
		if multiplayerLabelSize.X > maxWidth {
			maxWidth = multiplayerLabelSize.X
		}

		textX := modeLabelPos.X + (modeLabelSize.X-maxWidth)/2

		singleplayerLabelPos := rl.NewVector2(textX, float32(rl.GetScreenHeight()-int(singleplayerLabelSize.Y))/2-80)
		multiplayerLabelPos := rl.NewVector2(textX, float32(rl.GetScreenHeight()-int(multiplayerLabelSize.Y))/2)

		singleplayerRectangle := rl.NewRectangle(singleplayerLabelPos.X, singleplayerLabelPos.Y, 475, 65)
		multiplayerRectangle := rl.NewRectangle(multiplayerLabelPos.X, multiplayerLabelPos.Y, 430, 65)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, modeLabelText, modeLabelPos, 72, 2, rl.White)
		rl.DrawTextEx(font, singleplayerText, singleplayerLabelPos, 48, 2, rl.White)
		rl.DrawTextEx(font, multiplayerText, multiplayerLabelPos, 48, 2, rl.White)
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, singleplayerRectangle) {
				gameMode = SINGLEPLAYER
				if checkWorldFile("world.odn") {
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
				}
			}

			if rl.CheckCollisionPointRec(mousePos, multiplayerRectangle) {
				gameMode = MULTIPLAYER
				currentScene = IP_INPUT
			}
		}
	case GENERATE:
		generatingWorldLabelSize := rl.MeasureTextEx(font, "Generating world...", 56, 2)
		generatingWorldLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(generatingWorldLabelSize.X))/2, float32(rl.GetScreenHeight()-int(generatingWorldLabelSize.Y))/2)

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(font, "Generating world...", generatingWorldLabelPos, 56, 2, rl.White)
		rl.EndDrawing()

		// Проверка миров
		if !checkWorldFile("world.odn") {
			gameMode = SINGLEPLAYER
			generateWorld()
			saveWorldFile()
			saveWorldInfo()
			savePlayerFile()
			clear(world)
		}
		if !checkWorldFile("world_send.odn") {
			gameMode = MULTIPLAYER
			generateWorld()
			saveWorldFile()
			clear(world)
		}
		worldGenerated = true

		if worldGenerated {
			currentScene = MODE
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

				if gameMode == MULTIPLAYER {
					socket.Close()
					currentScene = TITLE
				} else {
					currentScene = SAVE
				}
				menuOpen = false
			}
		}
	case IP_INPUT:
		titleText := "Enter Server Details"
		titleSize := rl.MeasureTextEx(fontBold, titleText, 72, 2)
		titlePos := rl.NewVector2(float32(rl.GetScreenWidth()-int(titleSize.X))/2, 50)

		inputBoxWidth := 400
		inputBoxHeight := 50
		inputBoxX := (rl.GetScreenWidth() - inputBoxWidth) / 2
		inputBoxY := rl.GetScreenHeight()/2 - 150

		instructionText := "Press ENTER to connect or ESC to go back"
		instructionSize := rl.MeasureTextEx(font, instructionText, 24, 2)
		instructionPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(instructionSize.X))/2, float32(rl.GetScreenHeight()-50))

		var lineThickness float32 = 3.0

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTextEx(fontBold, titleText, titlePos, 72, 2, rl.White)

		// Nickname input
		nicknameBoxRec := rl.Rectangle{X: float32(inputBoxX), Y: float32(inputBoxY), Width: float32(inputBoxWidth), Height: float32(inputBoxHeight)}
		rl.DrawTextEx(font, "Nickname:", rl.Vector2{X: float32(inputBoxX), Y: float32(inputBoxY - 40)}, 24, 2, rl.White)
		rl.DrawRectangleLinesEx(nicknameBoxRec, lineThickness, rl.White)
		rl.DrawTextEx(font, nickname, rl.Vector2{X: float32(inputBoxX + 10), Y: float32(inputBoxY + 15)}, 24, 2, rl.LightGray)

		// Password input
		passwordBoxRec := rl.Rectangle{X: float32(inputBoxX), Y: float32(inputBoxY + 120), Width: float32(inputBoxWidth), Height: float32(inputBoxHeight)}
		rl.DrawTextEx(font, "Password:", rl.Vector2{X: float32(inputBoxX), Y: float32(inputBoxY + 80)}, 24, 2, rl.White)
		rl.DrawRectangleLinesEx(passwordBoxRec, lineThickness, rl.White)
		passwordText := strings.Repeat("*", len(password))
		maxPasswordWidth := inputBoxWidth - 20 // 20 пикселей отступа
		for len(passwordText) > 0 && rl.MeasureTextEx(font, passwordText, 24, 2).X > float32(maxPasswordWidth) {
			passwordText = passwordText[1:] // Обрезаем слева
		}
		rl.DrawTextEx(font, passwordText, rl.Vector2{X: float32(inputBoxX + 10), Y: float32(inputBoxY + 135)}, 24, 2, rl.LightGray)

		// IP address input
		ipBoxRec := rl.Rectangle{X: float32(inputBoxX), Y: float32(inputBoxY + 240), Width: float32(inputBoxWidth), Height: float32(inputBoxHeight)}
		rl.DrawTextEx(font, "IP Address:", rl.Vector2{X: float32(inputBoxX), Y: float32(inputBoxY + 200)}, 24, 2, rl.White)
		rl.DrawRectangleLinesEx(ipBoxRec, lineThickness, rl.White)
		rl.DrawTextEx(font, ipAddress, rl.Vector2{X: float32(inputBoxX + 10), Y: float32(inputBoxY + 255)}, 24, 2, rl.LightGray)

		rl.DrawTextEx(font, instructionText, instructionPos, 24, 2, rl.Gray)
		rl.EndDrawing()

		// Handle mouse input for selecting input fields
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePos := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePos, nicknameBoxRec) {
				activeInput = 0
			} else if rl.CheckCollisionPointRec(mousePos, passwordBoxRec) {
				activeInput = 1
			} else if rl.CheckCollisionPointRec(mousePos, ipBoxRec) {
				activeInput = 2
			}
		}

		// Handle keyboard input
		key := rl.GetCharPressed()
		if key >= 32 && key <= 126 {
			switch activeInput {
			case 0:
				if len(nickname) < 16 {
					nickname += string(key)
				}
			case 1:
				if len(password) < 32 {
					password += string(key)
				}
			case 2:
				if len(ipAddress) < 21 {
					ipAddress += string(key)
				}
			}
		}

		// Handle backspace
		if rl.IsKeyDown(rl.KeyBackspace) {
			backspaceHeldTime += rl.GetFrameTime()
			if backspaceHeldTime > backspaceRepeatDelay {
				if backspaceHeldTime-backspaceRepeatDelay > backspaceRepeatInterval {
					backspaceHeldTime = backspaceRepeatDelay
					switch activeInput {
					case 0:
						if len(nickname) > 0 {
							nickname = nickname[:len(nickname)-1]
						}
					case 1:
						if len(password) > 0 {
							password = password[:len(password)-1]
						}
					case 2:
						if len(ipAddress) > 0 {
							ipAddress = ipAddress[:len(ipAddress)-1]
						}
					}
				}
			} else if rl.IsKeyPressed(rl.KeyBackspace) {
				switch activeInput {
				case 0:
					if len(nickname) > 0 {
						nickname = nickname[:len(nickname)-1]
					}
				case 1:
					if len(password) > 0 {
						password = password[:len(password)-1]
					}
				case 2:
					if len(ipAddress) > 0 {
						ipAddress = ipAddress[:len(ipAddress)-1]
					}
				}
			}
		} else {
			backspaceHeldTime = 0
		}

		if rl.IsKeyPressed(rl.KeyEnter) && len(nickname) > 0 && len(password) > 0 && len(ipAddress) > 0 {
			if atomic.LoadInt32(&connectedToServer) == 0 {
				succesAuth := authPlayer()

				if succesAuth {
					checkStatusRest()
					connectServer("ws://" + ipAddress + "/ws")
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			currentScene = MODE
		}
	}
}
