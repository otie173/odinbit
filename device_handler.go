package main

import (
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	mousePos     rl.Vector2
	mouseOnBlock bool
)

func mouseHandler() {
	canBuild = true
	mouseOnBlock = false
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && currentScene == GAME {
		mousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)
		//mouseX := mousePos.X / TILE_SIZE
		//mouseY := mousePos.Y / TILE_SIZE

		playerX := int(math.Floor(float64(targetPosition.X / TILE_SIZE)))
		playerY := int(math.Floor(float64(targetPosition.Y / TILE_SIZE)))
		for _, block := range world {
			blockX := block.rec.X / TILE_SIZE
			blockY := block.rec.Y / TILE_SIZE
			if rl.CheckCollisionPointRec(mousePos, block.rec) {
				blockDistance := distanceInBlocks(float32(playerX)*TILE_SIZE, float32(playerY)*TILE_SIZE, blockX*TILE_SIZE, blockY*TILE_SIZE, float32(playerBlockDistance))
				if blockDistance {
					// Воспроизведение звука и разрушение объекта в зависимости от текстуры
					switch block.img {
					case smallTree, normalTree, bigTree, stone1, stone2, stone3, stone4:
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						pickupResourceSound()
						// Генерация травы на месте сломанного блока, чтобы не было просто пустого места
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
					case grass1, grass2, grass3, grass4, grass5, grass6, barrier:
						break
					default:
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						soundBlockAction()
						// Генерация травы на месте сломанного блока, чтобы не было просто пустого места
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)

					}
				}
				// Открытие слота инвентаря, если блок был неизвестен
				switch block.img {
				case wall:
					if !wallIsOpen {
						wallIsOpen = true
					}
					wallCount++
				case floor:
					if !floorIsOpen {
						floorIsOpen = true
					}
					floorCount++
				case door:
					if !doorIsOpen {
						doorIsOpen = true
					}
					doorCount++
				case chest:
					if !chestIsOpen {
						chestIsOpen = true
					}
					chestCount++
				case smallTree, normalTree, bigTree:
					woodCount += 1
				case stone1, stone2, stone3, stone4:
					stoneCount += 1
				}
			}
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && currentScene == GAME {
		mousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)
		mouseX := int(math.Floor(float64(mousePos.X / TILE_SIZE)))
		mouseY := int(math.Floor(float64(mousePos.Y / TILE_SIZE)))
		playerX := int(math.Floor(float64(targetPosition.X / TILE_SIZE)))
		playerY := int(math.Floor(float64(targetPosition.Y / TILE_SIZE)))
		blockDistance := distanceInBlocks(float32(playerX)*TILE_SIZE, float32(playerY)*TILE_SIZE, float32(mouseX)*TILE_SIZE, float32(mouseY)*TILE_SIZE, float32(playerBlockDistance))
		if targetPosition.X == (float32(mouseX)*TILE_SIZE) && targetPosition.Y == (float32(mouseY)*TILE_SIZE) && item != 2 {
			canBuild = false
		} else {
			canBuild = true
		}
		for _, block := range world {
			if rl.CheckCollisionPointRec(mousePos, block.rec) && block.img != grass1 && block.img != grass2 && block.img != grass3 && block.img != grass4 && block.img != grass5 && block.img != grass6 {
				mouseOnBlock = true
				break
			} else {
				mouseOnBlock = false
			}
		}
		if !mouseOnBlock && canBuild && mouseX > -(WORLD_SIZE/2) && mouseX < WORLD_SIZE/2 && mouseY > -(WORLD_SIZE/2) && mouseY < WORLD_SIZE/2 && blockDistance {
			switch item {
			case WALL:
				if wallIsOpen && wallCount != 0 {
					addBlock(wall, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					wallCount--
				}
			case FLOOR:
				if floorIsOpen && floorCount != 0 {
					addBlock(floor, float32(mouseX), float32(mouseY), true)
					soundBlockAction()
					floorCount--
				}
			case DOOR:
				if doorIsOpen && doorCount != 0 {
					addBlock(door, float32(mouseX), float32(mouseY), true)
					soundBlockAction()
					doorCount--
				}
			case CHEST:
				if chestIsOpen && chestCount != 0 {
					addBlock(chest, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					chestCount--
				}
			}
		}
	}
	// Крафт блоков из материалов
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && currentScene == INVENTORY {
		mousePos := rl.GetMousePosition()
		// Крафт стены
		if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
			if wallIsOpen && stoneCount-9 >= 0 {
				wallCount++
				stoneCount -= 9
			}
		}
		//Крафт пола
		if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
			if floorIsOpen && woodCount-9 >= 0 {
				floorCount++
				woodCount -= 9
			}
		}
	}
}

func keyboardHandler() {
	var moveX, moveY float32
	var shouldMove bool

	if rl.IsKeyDown(rl.KeyW) {
		moveY = -TILE_SIZE
		shouldMove = true
	} else if rl.IsKeyDown(rl.KeyS) {
		moveY = TILE_SIZE
		shouldMove = true
	}
	if rl.IsKeyDown(rl.KeyA) {
		playerDirection = false
		moveX = -TILE_SIZE
		shouldMove = true
	} else if rl.IsKeyDown(rl.KeyD) {
		playerDirection = true
		moveX = TILE_SIZE
		shouldMove = true
	}

	if shouldMove && canMoveAgain() {
		// Предполагаем, что движение возможно
		canMove := true
		// Проверяем диагональное движение только если оба направления активны
		if moveX != 0 && moveY != 0 {
			// Проверяем наличие блоков на диагональных путях
			for _, block := range world {
				if (!block.passable) && ((targetPosition.X+moveX == block.rec.X && targetPosition.Y == block.rec.Y) ||
					(targetPosition.Y+moveY == block.rec.Y && targetPosition.X == block.rec.X) ||
					(targetPosition.X+moveX == block.rec.X && targetPosition.Y+moveY == block.rec.Y)) {
					canMove = false
					break
				}
			}
		} else {
			// Проверяем наличие блоков на прямом пути
			for _, block := range world {
				if (!block.passable) && (targetPosition.X+moveX == block.rec.X && targetPosition.Y+moveY == block.rec.Y) {
					canMove = false
					break
				}
			}
		}

		if canMove {
			targetPosition.X += moveX
			targetPosition.Y += moveY
			lastMoveTime = time.Now()
		}
	}
	if rl.IsKeyPressed(rl.KeyE) && currentScene != TITLE && currentScene != MENU {
		switch inventoryOpen {
		case true:
			inventoryOpen = false
		case false:
			inventoryOpen = true
		}
		if inventoryOpen {
			currentScene = INVENTORY
		} else if !inventoryOpen {
			currentScene = GAME
		}
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		item = WALL
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		item = FLOOR
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		item = DOOR
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		item = CHEST
	}
	if rl.IsKeyPressed(rl.KeyI) {
		saveWorldFile()
	}
	if rl.IsKeyPressed(rl.KeyEscape) && currentScene != TITLE && currentScene != INVENTORY {
		switch menuOpen {
		case true:
			menuOpen = false
		case false:
			menuOpen = true
		}

		if menuOpen {
			currentScene = MENU
		} else if !menuOpen {
			currentScene = GAME
		}
	}
}
