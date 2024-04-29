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
	chestIsOpen        bool = false
	question           rl.Texture2D
	wallCount          int        = 0
	wallWindowCount    int        = 0
	floorCount         int        = 0
	doorCount          int        = 0
	chestCount         int        = 0
	inventoryZoom      float32    = 5.0
	shovelIsOpen       bool       = false
	pickaxeIsOpen      bool       = false
	axeIsOpen          bool       = false
	leftArrowPosition  rl.Vector2 = rl.NewVector2(0, 485)
	rightArrowPosition rl.Vector2 = rl.NewVector2(0, 485)
	arrowScale         float32    = 9.0
	currentPage        int        = 1
	maxPage            int        = 1
	pageLabelPos       rl.Vector2 = rl.NewVector2(0, 0)
)

type InventorySlot struct {
	x          float32
	y          float32
	slotNumber int
}

func createInventoryRow(startX, startY float32, slots, spacing int, inventory *[]InventorySlot, slotNum *int) {
	for i := 0; i < slots; i++ {
		slotX := startX + float32(i)*float32(slotImage.Width*int32(inventoryZoom)+int32(spacing))
		(*inventory)[*slotNum] = newInventorySlot(slotX, startY, *slotNum)
		*slotNum++
	}
}

func createInventoryRow2(startX, startY float32, slots, spacing int, inventory *[]InventorySlot, slotNum *int) {
	// Расчет ширины слота с учетом масштаба
	slotWidth := float32(slotImage.Width) * inventoryZoom

	// Расчет координаты X для центрального слота
	centerX := startX

	// Добавление центрального слота
	(*inventory)[*slotNum] = newInventorySlot(centerX, startY, *slotNum)
	*slotNum++

	// Расчет координаты X для левого слота
	leftX := centerX - slotWidth - float32(spacing)

	// Добавление левого слота
	(*inventory)[*slotNum] = newInventorySlot(leftX, startY, *slotNum)
	*slotNum++

	// Расчет координаты X для правого слота
	rightX := centerX + slotWidth + float32(spacing)

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
	otherTextures = []rl.Texture2D{wall, floor, door, chest, wallWindow}
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
		createInventoryRow2(startX2, yPos, 3, 32, &hotInventory, &hotInventorySlotNum)
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

func drawSlot(slot int, open bool) {
	if open {
		slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*inventoryZoom
		slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*inventoryZoom
		itemPosX := slotCenterX - (float32(otherTextures[slot].Width)/2)*inventoryZoom
		itemPosY := slotCenterY - (float32(otherTextures[slot].Height)/2)*inventoryZoom
		rl.DrawTextureEx(otherTextures[slot], rl.NewVector2(itemPosX, itemPosY), 0, inventoryZoom, rl.White)
	} else {
		slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*inventoryZoom
		slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*inventoryZoom
		itemPosX := slotCenterX - (float32(otherTextures[slot].Width)/2)*inventoryZoom
		itemPosY := slotCenterY - (float32(otherTextures[slot].Height)/2)*inventoryZoom
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

	if wallIsOpen {
		drawSlot(0, true)
		drawItemCount(hotInventory[0].x, hotInventory[0].y, wallCount, slotImage.Width, slotImage.Height)
	} else {
		drawSlot(0, false)
	}
	if floorIsOpen {
		drawSlot(1, true)
		drawItemCount(hotInventory[1].x, hotInventory[1].y, floorCount, slotImage.Width, slotImage.Height)
	} else {
		drawSlot(1, false)
	}
	if doorIsOpen {
		drawSlot(2, true)
		drawItemCount(hotInventory[2].x, hotInventory[2].y, doorCount, slotImage.Width, slotImage.Height)
	} else {
		drawSlot(2, false)
	}
	if chestIsOpen {
		drawSlot(3, true)
		drawItemCount(hotInventory[3].x, hotInventory[3].y, chestCount, slotImage.Width, slotImage.Height)
	} else {
		drawSlot(3, false)
	}
	if wallWindowIsOpen {
		drawSlot(4, true)
		drawItemCount(hotInventory[4].x, hotInventory[4].y, wallWindowCount, slotImage.Width, slotImage.Height)
	} else {
		drawSlot(4, false)
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
