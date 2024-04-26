package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	inventoryOpen bool
	slotImage     rl.Texture2D
	inventory     []InventorySlot
	hotInventory  []InventorySlot
	textures      []rl.Texture2D
	otherTextures []rl.Texture2D
	wood          rl.Texture2D
	stone         rl.Texture2D
	metal         rl.Texture2D
	woodCount     int  = 0
	stoneCount    int  = 0
	metalCount    int  = 0
	wallIsOpen    bool = false
	floorIsOpen   bool = false
	doorIsOpen    bool = false
	chestIsOpen   bool = false
	question      rl.Texture2D
	wallCount     int     = 0
	floorCount    int     = 0
	doorCount     int     = 0
	chestCount    int     = 0
	inventoryZoom float32 = 5.0
	shovelIsOpen  bool    = false
	pickaxeIsOpen bool    = false
	axeIsOpen     bool    = false
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

func loadInventory() {
	slotImage = loadTexture("assets/images/gui/slot.png")
	inventory = make([]InventorySlot, 3)
	hotInventory = make([]InventorySlot, 32)
	textures = make([]rl.Texture2D, 3)
	otherTextures = make([]rl.Texture2D, 32)
	wood = loadTexture("assets/images/items/wood.png")
	stone = loadTexture("assets/images/items/stone.png")
	metal = loadTexture("assets/images/items/metal.png")
	textures = []rl.Texture2D{wood, stone, metal}
	otherTextures = []rl.Texture2D{wall, floor, door, chest}
	question = loadTexture("assets/images/gui/question.png")

	inventoryLabelSize := rl.MeasureTextEx(fontBold, "Inventory", 72, 2)
	inventoryLabelPos := rl.NewVector2(float32(rl.GetScreenWidth()-int(inventoryLabelSize.X))/2, 75)
	startX1 := inventoryLabelPos.X + 50
	inventorySlotNum := 0
	hotInventorySlotNum := 0

	createInventoryRow(startX1, 182.0, 3, 110, &inventory, &inventorySlotNum)
	yPositions := []float32{395.0, 495.0, 595.0, 695.0}
	for _, yPos := range yPositions {
		createInventoryRow(startX1-200, yPos, 8, 32, &hotInventory, &hotInventorySlotNum)
	}
}

func unloadInventory() {
	rl.UnloadTexture(slotImage)
	for i := range textures {
		rl.UnloadTexture(textures[i])
	}
	rl.UnloadTexture(question)

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
