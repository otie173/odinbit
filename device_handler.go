package main

import (
	"math"

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
		for _, block := range world {
			if rl.CheckCollisionPointRec(mousePos, block.rec) {

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
		if playerPosition.X == (float32(mouseX)*10) && playerPosition.Y == (float32(mouseY)*10) && item != 2 {
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
		if !mouseOnBlock && canBuild && mouseX > -(WORLD_SIZE/2) && mouseX < WORLD_SIZE/2 && mouseY > -(WORLD_SIZE/2) && mouseY < WORLD_SIZE/2 {
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
	if rl.IsKeyPressed(rl.KeyW) && currentScene == GAME {
		for _, block := range world {
			if playerPosition.Y-TILE_SIZE == block.rec.Y && playerPosition.X == block.rec.X && !block.passable {
				canMove = false
				break
			} else {
				canMove = true
			}
		}
		if canMove {
			playerPosition.Y -= TILE_SIZE
		}
	}
	if rl.IsKeyPressed(rl.KeyA) && currentScene == GAME {
		playerDirection = false
		for _, block := range world {
			if playerPosition.X-TILE_SIZE == block.rec.X && playerPosition.Y == block.rec.Y && !block.passable {
				canMove = false
				break
			} else {
				canMove = true
			}
		}
		if canMove {
			playerPosition.X -= TILE_SIZE
		}
	}
	if rl.IsKeyPressed(rl.KeyS) && currentScene == GAME {
		for _, block := range world {
			if playerPosition.Y+TILE_SIZE == block.rec.Y && playerPosition.X == block.rec.X && !block.passable {
				canMove = false
				break
			} else {
				canMove = true
			}
		}
		if canMove {
			playerPosition.Y += TILE_SIZE
		}
	}
	if rl.IsKeyPressed(rl.KeyD) && currentScene == GAME {
		playerDirection = true
		for _, block := range world {
			if playerPosition.X+TILE_SIZE == block.rec.X && playerPosition.Y == block.rec.Y && !block.passable {
				canMove = false
				break
			} else {
				canMove = true
			}
		}
		if canMove {
			playerPosition.X += TILE_SIZE
		}
	}
	if rl.IsKeyPressed(rl.KeyE) && currentScene != TITLE {
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
}
