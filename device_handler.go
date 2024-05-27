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

const (
	MATERIAL_FOR_CRAFT int = 25
)

func mouseHandler() {
	canBuild = true
	mouseOnBlock = false
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && currentScene == GAME {
		mousePos = rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)

		playerX := int(math.Floor(float64(targetPosition.X / TILE_SIZE)))
		playerY := int(math.Floor(float64(targetPosition.Y / TILE_SIZE)))
		for _, block := range world {
			blockX := block.rec.X / TILE_SIZE
			blockY := block.rec.Y / TILE_SIZE
			if rl.CheckCollisionPointRec(mousePos, block.rec) {
				blockDistance := distanceInBlocks(float32(playerX)*TILE_SIZE, float32(playerY)*TILE_SIZE, blockX*TILE_SIZE, blockY*TILE_SIZE, float32(playerBlockDistance))
				if blockDistance {
					switch block.img {
					case smallTree, normalTree, bigTree:
						if axeIsOpen {
							removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							pickupResourceSound()
							generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							woodCount += 5
						}
					case stone1, stone2, stone3, stone4, bigStone1, bigStone2, bigStone3, bigStone4, bigStone5:
						if pickaxeIsOpen {
							removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							pickupResourceSound()
							generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							stoneCount += 5
						}
					case grass1, grass2, grass3, grass4, grass5, grass6, barrier:
						break
					case wall, wallWindow, tombstone, trash, lamp, floor4, lootbox:
						if pickaxeIsOpen {
							removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							soundBlockAction()
							generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)

							if block.img == wall {
								wallIsOpen = true
								wallCount++
							}
							if block.img == wallWindow {
								wallWindowIsOpen = true
								wallWindowCount++
							}
							if block.img == tombstone {
								tombstoneIsOpen = true
								tombstoneCount++
							}
							if block.img == trash {
								trashIsOpen = true
								trashCount++
							}
							if block.img == lamp {
								lampIsOpen = true
								lampCount++
							}
							if block.img == floor4 {
								floor4IsOpen = true
								floor4Count++
							}
							if block.img == lootbox {
								lootboxIsOpen = true
								lootboxCount++
							}
						}
					case floor, door, chest, doorOpen, bigBarrel, bookshelf, chair, closet, fence1, fence2, floor2, shelf, sign, smallBarrel, table:
						if axeIsOpen {
							removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
							soundBlockAction()
							generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)

							if block.img == floor {
								floorIsOpen = true
								floorCount++
							}
							if block.img == door {
								doorIsOpen = true
								doorCount++
							}
							if block.img == chest {
								chestIsOpen = true
								chestCount++
							}
							if block.img == doorOpen {
								doorOpenIsOpen = true
								doorOpenCount++
							}
							if block.img == bigBarrel {
								bigBarrelIsOpen = true
								bigBarrelCount++
							}
							if block.img == bookshelf {
								bookshelfIsOpen = true
								bookshelfCount++
							}
							if block.img == chair {
								chairIsOpen = true
								chairCount++
							}
							if block.img == closet {
								closetIsOpen = true
								closetCount++
							}
							if block.img == fence1 {
								fence1IsOpen = true
								fence1Count++
							}
							if block.img == fence2 {
								fence2IsOpen = true
								fence2Count++
							}
							if block.img == floor2 {
								floor2IsOpen = true
								floor2Count++
							}
							if block.img == shelf {
								shelfIsOpen = true
								shelfCount++
							}
							if block.img == sign {
								signIsOpen = true
								signCount++
							}
							if block.img == smallBarrel {
								smallBarrelIsOpen = true
								smallBarrelCount++
							}
							if block.img == table {
								tableIsOpen = true
								tableCount++
							}
						}
					case bones1, bones2, bones3, bones4, bones5:
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						pickupResourceSound()
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
					case pickaxe:
						pickaxeIsOpen = true
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						pickupResourceSound()
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
					case axe:
						axeIsOpen = true
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						pickupResourceSound()
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
					case shovel:
						shovelIsOpen = true
						removeBlock(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
						pickupResourceSound()
						generateGrass(block.rec.X/TILE_SIZE, block.rec.Y/TILE_SIZE)
					}
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
			case WALLWINDOW:
				if wallWindowIsOpen && wallWindowCount != 0 {
					addBlock(wallWindow, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					wallWindowCount--
				}
			case DOOROPEN:
				if doorOpenIsOpen && doorOpenCount != 0 {
					addBlock(doorOpen, float32(mouseX), float32(mouseY), true)
					soundBlockAction()
					doorOpenCount--
				}
			case BIGBARREL:
				if bigBarrelIsOpen && bigBarrelCount != 0 {
					addBlock(bigBarrel, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					bigBarrelCount--
				}
			case BOOKSHELF:
				if bookshelfIsOpen && bookshelfCount != 0 {
					addBlock(bookshelf, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					bookshelfCount--
				}
			case CHAIR:
				if chairIsOpen && chairCount != 0 {
					addBlock(chair, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					chairCount--
				}
			case CLOSET:
				if closetIsOpen && closetCount != 0 {
					addBlock(closet, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					closetCount--
				}
			case FENCE1:
				if fence1IsOpen && fence1Count != 0 {
					addBlock(fence1, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					fence1Count--
				}
			case FENCE2:
				if fence2IsOpen && fence2Count != 0 {
					addBlock(fence2, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					fence2Count--
				}
			case FLOOR2:
				if floor2IsOpen && floor2Count != 0 {
					addBlock(floor2, float32(mouseX), float32(mouseY), true)
					soundBlockAction()
					floor2Count--
				}
			case FLOOR4:
				if floor4IsOpen && floor4Count != 0 {
					addBlock(floor4, float32(mouseX), float32(mouseY), true)
					soundBlockAction()
					floor4Count--
				}
			case LAMP:
				if lampIsOpen && lampCount != 0 {
					addBlock(lamp, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					lampCount--
				}
			case LOOTBOX:
				if lootboxIsOpen && lootboxCount != 0 {
					addBlock(lootbox, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					lootboxCount--
				}
			case SHELF:
				if shelfIsOpen && shelfCount != 0 {
					addBlock(shelf, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					shelfCount--
				}
			case SIGN:
				if signIsOpen && signCount != 0 {
					addBlock(sign, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					signCount--
				}
			case SMALLBARREL:
				if smallBarrelIsOpen && smallBarrelCount != 0 {
					addBlock(smallBarrel, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					smallBarrelCount--
				}
			case TABLE:
				if tableIsOpen && tableCount != 0 {
					addBlock(table, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					tableCount--
				}
			case TOMBSTONE:
				if tombstoneIsOpen && tombstoneCount != 0 {
					addBlock(tombstone, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					tombstoneCount--
				}
			case TRASH:
				if trashIsOpen && trashCount != 0 {
					addBlock(trash, float32(mouseX), float32(mouseY), false)
					soundBlockAction()
					trashCount--
				}
			}

		}
	}
	// Крафт блоков из материалов
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && currentScene == INVENTORY {
		mousePos := rl.GetMousePosition()
		switch currentPage {
		case 1:
			// Крафт стены
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if wallIsOpen && stoneCount-MATERIAL_FOR_CRAFT >= 0 {
					wallCount++
					stoneCount -= MATERIAL_FOR_CRAFT
				}
			}
			//Крафт пола
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floorIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					floorCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			//Крафт двери
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if doorIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					doorCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			//Крафт сундука
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if chestIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					chestCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт окна
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[4].x, hotInventory[4].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if wallWindowIsOpen && stoneCount-MATERIAL_FOR_CRAFT >= 0 {
					wallWindowCount++
					stoneCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт дверного проёма
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[5].x, hotInventory[5].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if doorOpenIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					doorOpenCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт большой бочки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[6].x, hotInventory[6].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if bigBarrelIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					bigBarrelCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт книжной полки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[7].x, hotInventory[7].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if bookshelfIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					bookshelfCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт стула
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[8].x, hotInventory[8].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if chairIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					chairCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
		case 2:
			// Крафт шкафа
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if closetIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					closetCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт забора1
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if fence1IsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					fence1Count++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт забора2
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if fence2IsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					fence2Count++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт пол2
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floor2IsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					floor2Count++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт пол4
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[4].x, hotInventory[4].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floor4IsOpen && stoneCount-MATERIAL_FOR_CRAFT >= 0 {
					floor4Count++
					stoneCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт лампы
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[5].x, hotInventory[5].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if lampIsOpen && metalCount-MATERIAL_FOR_CRAFT >= 0 {
					lampCount++
					metalCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт кувшина
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[6].x, hotInventory[6].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if lootboxIsOpen && stoneCount-MATERIAL_FOR_CRAFT >= 0 {
					lootboxCount++
					stoneCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт полки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[7].x, hotInventory[7].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if shelfIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					shelfCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт таблички
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[8].x, hotInventory[8].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if signIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					signCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
		case 3:
			// Крафт маленькой бочки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if smallBarrelIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					smallBarrelCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт стола
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if tableIsOpen && woodCount-MATERIAL_FOR_CRAFT >= 0 {
					tableCount++
					woodCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт надгробия
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if tombstoneIsOpen && stoneCount-MATERIAL_FOR_CRAFT >= 0 {
					tombstoneCount++
					stoneCount -= MATERIAL_FOR_CRAFT
				}
			}
			// Крафт мусорки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if trashIsOpen && metalCount-MATERIAL_FOR_CRAFT >= 0 {
					trashCount++
					metalCount -= MATERIAL_FOR_CRAFT
				}
			}
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && currentScene == INVENTORY {
		mousePos := rl.GetMousePosition()
		switch currentPage {
		case 1:
			// Декрафт стены
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if wallIsOpen && wallCount > 0 {
					wallCount--
					stoneCount += MATERIAL_FOR_CRAFT
				}
			}
			//Декрафт пола
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floorIsOpen && floorCount > 0 {
					floorCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			//Декрафт двери
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if doorIsOpen && doorCount > 0 {
					doorCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			//Декрафт сундука
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if chestIsOpen && chestCount > 0 {
					chestCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт окна
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[4].x, hotInventory[4].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if wallWindowIsOpen && wallWindowCount > 0 {
					wallWindowCount--
					stoneCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт дверного проёма
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[5].x, hotInventory[5].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if doorOpenIsOpen && doorOpenCount > 0 {
					doorOpenCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт большой бочки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[6].x, hotInventory[6].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if bigBarrelIsOpen && bigBarrelCount > 0 {
					bigBarrelCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт книжной полки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[7].x, hotInventory[7].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if bookshelfIsOpen && bookshelfCount > 0 {
					bookshelfCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт стула
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[8].x, hotInventory[8].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if chairIsOpen && chairCount > 0 {
					chairCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
		case 2:
			// Декрафт шкафа
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if closetIsOpen && closetCount > 0 {
					closetCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт забора1
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if fence1IsOpen && fence1Count > 0 {
					fence1Count--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт забора2
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if fence2IsOpen && fence2Count > 0 {
					fence2Count--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт пол2
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floor2IsOpen && floor2Count > 0 {
					floor2Count--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт пол4
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[4].x, hotInventory[4].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if floor4IsOpen && floor4Count > 0 {
					floor4Count--
					stoneCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт лампы
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[5].x, hotInventory[5].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if lampIsOpen && lampCount > 0 {
					lampCount--
					metalCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт кувшина
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[6].x, hotInventory[6].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if lootboxIsOpen && lootboxCount > 0 {
					lootboxCount--
					stoneCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт полки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[7].x, hotInventory[7].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if shelfIsOpen && shelfCount > 0 {
					shelfCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт таблички
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[8].x, hotInventory[8].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if signIsOpen && signCount > 0 {
					signCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
		case 3:
			// Декрафт маленькой бочки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[0].x, hotInventory[0].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if smallBarrelIsOpen && smallBarrelCount > 0 {
					smallBarrelCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт стола
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[1].x, hotInventory[1].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if tableIsOpen && tableCount > 0 {
					tableCount--
					woodCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт надгробия
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[2].x, hotInventory[2].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if tombstoneIsOpen && tombstoneCount > 0 {
					tombstoneCount--
					stoneCount += MATERIAL_FOR_CRAFT
				}
			}
			// Декрафт мусорки
			if rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(hotInventory[3].x, hotInventory[3].y, float32(slotImage.Width)*inventoryZoom, float32(slotImage.Height)*cam.Zoom)) {
				if trashIsOpen && trashCount > 0 {
					trashCount--
					metalCount += MATERIAL_FOR_CRAFT
				}
			}
		}
	}
}

func keyboardHandler() {
	var moveX, moveY float32
	var shouldMove bool

	if currentScene == GAME {
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
			canMove := true
			if moveX != 0 && moveY != 0 {
				for _, block := range world {
					if (!block.passable) && ((targetPosition.X+moveX == block.rec.X && targetPosition.Y == block.rec.Y) ||
						(targetPosition.Y+moveY == block.rec.Y && targetPosition.X == block.rec.X) ||
						(targetPosition.X+moveX == block.rec.X && targetPosition.Y+moveY == block.rec.Y)) {
						canMove = false
						break
					}
				}
			} else {
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
	switch currentPage {
	case 1:
		itemMap := map[int32]int{
			rl.KeyOne:   WALL,
			rl.KeyTwo:   FLOOR,
			rl.KeyThree: DOOR,
			rl.KeyFour:  CHEST,
			rl.KeyFive:  WALLWINDOW,
			rl.KeySix:   DOOROPEN,
			rl.KeySeven: BIGBARREL,
			rl.KeyEight: BOOKSHELF,
			rl.KeyNine:  CHAIR,
		}

		for key, value := range itemMap {
			if rl.IsKeyPressed(key) {
				item = value
				break
			}
		}
	case 2:
		itemMap := map[int32]int{
			rl.KeyOne:   CLOSET,
			rl.KeyTwo:   FENCE1,
			rl.KeyThree: FENCE2,
			rl.KeyFour:  FLOOR2,
			rl.KeyFive:  FLOOR4,
			rl.KeySix:   LAMP,
			rl.KeySeven: LOOTBOX,
			rl.KeyEight: SHELF,
			rl.KeyNine:  SIGN,
		}

		for key, value := range itemMap {
			if rl.IsKeyPressed(key) {
				item = value
				break
			}
		}
	case 3:
		itemMap := map[int32]int{
			rl.KeyOne:   SMALLBARREL,
			rl.KeyTwo:   TABLE,
			rl.KeyThree: TOMBSTONE,
			rl.KeyFour:  TRASH,
		}

		for key, value := range itemMap {
			if rl.IsKeyPressed(key) {
				item = value
				break
			}
		}
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
