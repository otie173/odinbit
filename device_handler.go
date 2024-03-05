package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	mousePos     rl.Vector2
	mouseOnBlock bool
	item         int
)

func mouseHandler() {
	canBuild = true
	mouseOnBlock = false
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && currentScene == GAME {
		mousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)
		for _, block := range world {
			if rl.CheckCollisionPointRec(mousePos, block.rec) {
				removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
				soundBlockRemove()
			}
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && currentScene == GAME {
		mousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)
		mouseX := int(math.Floor(float64(mousePos.X / TILE_SIZE)))
		mouseY := int(math.Floor(float64(mousePos.Y / TILE_SIZE)))
		if playerPosition.X == (float32(mouseX)*10) && playerPosition.Y == (float32(mouseY)*10) {
			canBuild = false
		} else {
			canBuild = true
		}
		for _, block := range world {
			if rl.CheckCollisionPointRec(mousePos, block.rec) {
				mouseOnBlock = true
				break
			} else {
				mouseOnBlock = false
			}
		}
		if !mouseOnBlock && canBuild {
			switch item {
			case 1:
				addBlock(wall, float32(mouseX), float32(mouseY), false)
				soundBlockAdd()
			case 2:
				addBlock(floor, float32(mouseX), float32(mouseY), true)
				soundBlockAdd()
			case 3:
				addBlock(chestClose, float32(mouseX), float32(mouseY), false)
				soundBlockAdd()
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
		item = 1
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		item = 2
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		item = 3
	}
}
