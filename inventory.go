package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	inventoryOpen      bool
	slotImage          rl.Texture2D
	inventory          []InventorySlot
	hotInventory       []InventorySlot
	textures           []rl.Texture2D
	otherTextures      []rl.Texture2D
	wood               rl.Texture2D
	stone              rl.Texture2D
	metal              rl.Texture2D
	leftArrow          rl.Texture2D
	rightArrow         rl.Texture2D
	woodCount          int  = 0
	stoneCount         int  = 0
	metalCount         int  = 0
	wallIsOpen         bool = false
	wallWindowIsOpen   bool = false
	floorIsOpen        bool = false
	doorIsOpen         bool = false
	doorOpenIsOpen     bool = false
	chestIsOpen        bool = false
	question           rl.Texture2D
	wallCount          int        = 0
	wallWindowCount    int        = 0
	floorCount         int        = 0
	doorCount          int        = 0
	doorOpenCount      int        = 0
	chestCount         int        = 0
	inventoryZoom      float32    = 5.0
	shovelIsOpen       bool       = false
	pickaxeIsOpen      bool       = false
	axeIsOpen          bool       = false
	leftArrowPosition  rl.Vector2 = rl.NewVector2(0, 485)
	rightArrowPosition rl.Vector2 = rl.NewVector2(0, 485)
	arrowScale         float32    = 9.0
	currentPage        int        = 1
	maxPage            int        = 3
	bigBarrelIsOpen    bool       = false
	bookshelfIsOpen    bool       = false
	chairIsOpen        bool       = false
	closetIsOpen       bool       = false
	fence1IsOpen       bool       = false
	fence2IsOpen       bool       = false
	floor2IsOpen       bool       = false
	floor4IsOpen       bool       = false
	lampIsOpen         bool       = false
	shelfIsOpen        bool       = false
	signIsOpen         bool       = false
	smallBarrelIsOpen  bool       = false
	tableIsOpen        bool       = false
	trashIsOpen        bool       = false
	bigBarrelCount     int        = 0
	bookshelfCount     int        = 0
	chairCount         int        = 0
	closetCount        int        = 0
	fence1Count        int        = 0
	fence2Count        int        = 0
	floor2Count        int        = 0
	floor4Count        int        = 0
	lampCount          int        = 0
	shelfCount         int        = 0
	signCount          int        = 0
	smallBarrelCount   int        = 0
	tableCount         int        = 0
	trashCount         int        = 0
	lootboxIsOpen      bool       = false
	lootboxCount       int        = 0
	tombstoneCount     int        = 0
	tombstoneIsOpen    bool       = false
	saplingIsOpen      bool       = false
	saplingCount       int        = 0
	seedIsOpen         bool       = false
	seedCount          int        = 0
	cabbageIsOpen      bool       = false
	cabbageCount       int        = 0
)

type InventorySlot struct {
	x          float32
	y          float32
	slotNumber int
}

type SlotInfo struct {
	isOpen       bool
	count        int
	x            float32
	y            float32
	textureIndex int
}

func createInventoryRow(startX, startY float32, slots, spacing int, inventory *[]InventorySlot, slotNum *int) {
	for i := 0; i < slots; i++ {
		slotX := startX + float32(i)*float32(slotImage.Width*int32(inventoryZoom)+int32(spacing))
		(*inventory)[*slotNum] = newInventorySlot(slotX, startY, *slotNum)
		*slotNum++
	}
}

func createInventoryRow2(startX, startY float32, spacing int, inventory *[]InventorySlot, slotNum *int) {
	// Расчет ширины слота с учетом масштаба
	slotWidth := float32(slotImage.Width) * inventoryZoom

	// Расчет координаты X для центрального слота
	centerX := startX

	// Расчет координаты X для левого слота
	leftX := centerX - slotWidth - float32(spacing)

	// Расчет координаты X для правого слота
	rightX := centerX + slotWidth + float32(spacing)

	// Добавление левого слота
	(*inventory)[*slotNum] = newInventorySlot(leftX, startY, *slotNum)
	*slotNum++

	// Добавление центрального слота
	(*inventory)[*slotNum] = newInventorySlot(centerX, startY, *slotNum)
	*slotNum++

	// Добавление правого слота
	(*inventory)[*slotNum] = newInventorySlot(rightX, startY, *slotNum)
	*slotNum++

	if startY == 495.0 {
		leftArrowPosition.X = (float32(rl.GetScreenWidth())-float32(leftArrow.Width)*arrowScale)/2.0 - 200
		rightArrowPosition.X = (float32(rl.GetScreenWidth())-float32(rightArrow.Width)*arrowScale)/2.0 + 200
	}
}

func loadInventory() {
	slotImage = loadTexture("assets/images/gui/slot.png")
	inventory = make([]InventorySlot, 3)
	hotInventory = make([]InventorySlot, 9)
	textures = make([]rl.Texture2D, 3)
	otherTextures = make([]rl.Texture2D, 32)
	wood = loadTexture("assets/images/items/wood.png")
	stone = loadTexture("assets/images/items/stone.png")
	metal = loadTexture("assets/images/items/metal.png")
	textures = []rl.Texture2D{wood, stone, metal}
	otherTextures = []rl.Texture2D{wall, floor, door, chest, wallWindow, doorOpen, bigBarrel, bookshelf, chair, closet, fence1, fence2, floor2, floor4, lamp, lootbox, shelf, sign, smallBarrel, table, tombstone, trash, sapling, seed2Small, seed1Big}
	question = loadTexture("assets/images/gui/question.png")
	leftArrow = loadTexture("assets/images/gui/left_arrow.png")
	rightArrow = loadTexture("assets/images/gui/right_arrow.png")

	inventoryLabelSize := rl.MeasureTextEx(fontBold, "Inventory", 72, 2)
	inventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(inventoryLabelSize.X))/2, 75)
	startX1 := inventoryLabelPos.X + 50
	startX2 := (float32(rl.GetMonitorWidth(rl.GetCurrentMonitor())) - float32(slotImage.Width)*inventoryZoom) / 2.0
	inventorySlotNum := 0
	hotInventorySlotNum := 0

	createInventoryRow(startX1, 182.0, 3, 110, &inventory, &inventorySlotNum)
	yPositions := []float32{395.0, 495.0, 595.0}
	for _, yPos := range yPositions {
		createInventoryRow2(startX2, yPos, 32, &hotInventory, &hotInventorySlotNum)
	}
}

func unloadInventory() {
	rl.UnloadTexture(slotImage)
	for i := range textures {
		rl.UnloadTexture(textures[i])
	}
	for i := range otherTextures {
		rl.UnloadTexture(otherTextures[i])
	}
	rl.UnloadTexture(question)
	rl.UnloadTexture(leftArrow)
	rl.UnloadTexture(rightArrow)

}

func newInventorySlot(x, y float32, slotNumber int) InventorySlot {
	return InventorySlot{
		x:          x,
		y:          y,
		slotNumber: slotNumber,
	}
}

func drawSlot(slot int, open bool, textureIndex int) {
	if open {
		if textureIndex >= 0 && textureIndex < len(otherTextures) {
			slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*inventoryZoom
			slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*inventoryZoom
			itemPosX := slotCenterX - (float32(otherTextures[textureIndex].Width)/2)*inventoryZoom
			itemPosY := slotCenterY - (float32(otherTextures[textureIndex].Height)/2)*inventoryZoom
			rl.DrawTextureEx(otherTextures[textureIndex], rl.NewVector2(itemPosX, itemPosY), 0, inventoryZoom, rl.White)
		}
	} else {
		slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*inventoryZoom
		slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*inventoryZoom
		itemPosX := slotCenterX - (float32(question.Width)/2)*inventoryZoom
		itemPosY := slotCenterY - (float32(question.Height)/2)*inventoryZoom
		rl.DrawTextureEx(question, rl.NewVector2(itemPosX, itemPosY), 0, inventoryZoom, rl.White)
	}
}

func drawItems() {
	for i, slot := range inventory {
		slotCenterX := slot.x + (float32(slotImage.Width)/2)*inventoryZoom
		slotCenterY := slot.y + (float32(slotImage.Height)/2)*inventoryZoom
		if i < len(textures) {
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*inventoryZoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*inventoryZoom
			rl.DrawTextureEx(textures[i], rl.NewVector2(itemPosX, itemPosY), 0, inventoryZoom, rl.White)

			var itemCount int
			switch i {
			case 0:
				itemCount = woodCount
			case 1:
				itemCount = stoneCount
			case 2:
				itemCount = metalCount
			}
			drawItemCount(slot.x, slot.y, itemCount, slotImage.Width, slotImage.Height)
		}
	}
	switch currentPage {
	case 1:
		slotData := []SlotInfo{
			{isOpen: wallIsOpen, count: wallCount, x: hotInventory[0].x, y: hotInventory[0].y, textureIndex: 0},
			{isOpen: floorIsOpen, count: floorCount, x: hotInventory[1].x, y: hotInventory[1].y, textureIndex: 1},
			{isOpen: doorIsOpen, count: doorCount, x: hotInventory[2].x, y: hotInventory[2].y, textureIndex: 2},
			{isOpen: chestIsOpen, count: chestCount, x: hotInventory[3].x, y: hotInventory[3].y, textureIndex: 3},
			{isOpen: wallWindowIsOpen, count: wallWindowCount, x: hotInventory[4].x, y: hotInventory[4].y, textureIndex: 4},
			{isOpen: doorOpenIsOpen, count: doorOpenCount, x: hotInventory[5].x, y: hotInventory[5].y, textureIndex: 5},
			{isOpen: bigBarrelIsOpen, count: bigBarrelCount, x: hotInventory[6].x, y: hotInventory[6].y, textureIndex: 6},
			{isOpen: bookshelfIsOpen, count: bookshelfCount, x: hotInventory[7].x, y: hotInventory[7].y, textureIndex: 7},
			{isOpen: chairIsOpen, count: chairCount, x: hotInventory[8].x, y: hotInventory[8].y, textureIndex: 8},
		}

		for i, slot := range slotData {
			if slot.isOpen {
				drawSlot(i, true, slot.textureIndex)
				drawItemCount(slot.x, slot.y, slot.count, slotImage.Width, slotImage.Height)
			}
			if !slot.isOpen {
				drawSlot(i, false, slot.textureIndex)
			}
		}

		// Сбросить состояние неиспользуемых слотов
		for i := len(slotData); i < len(hotInventory); i++ {
			drawSlot(i, false, 0)
		}
	case 2:
		slotData := []SlotInfo{
			{isOpen: closetIsOpen, count: closetCount, x: hotInventory[0].x, y: hotInventory[0].y, textureIndex: 9},
			{isOpen: fence1IsOpen, count: fence1Count, x: hotInventory[1].x, y: hotInventory[1].y, textureIndex: 10},
			{isOpen: fence2IsOpen, count: fence2Count, x: hotInventory[2].x, y: hotInventory[2].y, textureIndex: 11},
			{isOpen: floor2IsOpen, count: floor2Count, x: hotInventory[3].x, y: hotInventory[3].y, textureIndex: 12},
			{isOpen: floor4IsOpen, count: floor4Count, x: hotInventory[4].x, y: hotInventory[4].y, textureIndex: 13},
			{isOpen: lampIsOpen, count: lampCount, x: hotInventory[5].x, y: hotInventory[5].y, textureIndex: 14},
			{isOpen: lootboxIsOpen, count: lootboxCount, x: hotInventory[6].x, y: hotInventory[6].y, textureIndex: 15},
			{isOpen: shelfIsOpen, count: shelfCount, x: hotInventory[7].x, y: hotInventory[7].y, textureIndex: 16},
			{isOpen: signIsOpen, count: signCount, x: hotInventory[8].x, y: hotInventory[8].y, textureIndex: 17},
		}

		for i, slot := range slotData {
			if slot.isOpen {
				drawSlot(i, true, slot.textureIndex)
				drawItemCount(slot.x, slot.y, slot.count, slotImage.Width, slotImage.Height)
			}
			if !slot.isOpen {
				drawSlot(i, false, slot.textureIndex)
			}
		}

		// Сбросить состояние неиспользуемых слотов
		for i := len(slotData); i < len(hotInventory); i++ {
			drawSlot(i, false, 0)
		}
	case 3:
		slotData := []SlotInfo{
			{isOpen: smallBarrelIsOpen, count: smallBarrelCount, x: hotInventory[0].x, y: hotInventory[0].y, textureIndex: 18},
			{isOpen: tableIsOpen, count: tableCount, x: hotInventory[1].x, y: hotInventory[1].y, textureIndex: 19},
			{isOpen: tombstoneIsOpen, count: tombstoneCount, x: hotInventory[2].x, y: hotInventory[2].y, textureIndex: 20},
			{isOpen: trashIsOpen, count: trashCount, x: hotInventory[3].x, y: hotInventory[3].y, textureIndex: 21},
			{isOpen: saplingIsOpen, count: saplingCount, x: hotInventory[4].x, y: hotInventory[4].y, textureIndex: 22},
			{isOpen: seedIsOpen, count: seedCount, x: hotInventory[5].x, y: hotInventory[5].y, textureIndex: 23},
			{isOpen: cabbageIsOpen, count: cabbageCount, x: hotInventory[6].x, y: hotInventory[6].y, textureIndex: 24},
		}

		for i, slot := range slotData {
			if slot.isOpen {
				drawSlot(i, true, slot.textureIndex)
				drawItemCount(slot.x, slot.y, slot.count, slotImage.Width, slotImage.Height)
			}
			if !slot.isOpen {
				drawSlot(i, false, slot.textureIndex)
			}
		}

		// Сбросить состояние неиспользуемых слотов
		for i := len(slotData); i < 3; i++ {
			drawSlot(i, false, 0)
		}
	}
}

func drawItemCount(x, y float32, count int, slotWidth, slotHeight int32) {
	itemCountText := fmt.Sprintf("%d", count)
	fontSize := float32(16)
	fontSpacing := float32(2)
	var itemCountTextWidth float32
	for {
		itemCountTextWidth = rl.MeasureTextEx(font, itemCountText, fontSize, fontSpacing).X
		if itemCountTextWidth+14 <= float32(slotWidth)*inventoryZoom {
			break
		}
		fontSize -= 1
		if fontSize <= 8 {
			break
		}
	}
	itemCountPosX := x + float32(slotWidth)*inventoryZoom - itemCountTextWidth - 7
	itemCountPosY := y + float32(slotHeight)*inventoryZoom - fontSize - 4
	rl.DrawTextEx(font, itemCountText, rl.NewVector2(itemCountPosX, itemCountPosY), fontSize, fontSpacing, rl.White)
}
