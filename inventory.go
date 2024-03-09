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
	wallIsOpen    bool = true
	floorIsOpen   bool = true
	question      rl.Texture2D
)

type InventorySlot struct {
	x          float32
	y          float32
	slotNumber int
}

func createInventoryRow(startX, startY float32, slots, spacing int, inventory *[]InventorySlot, slotNum *int) {
	for i := 0; i < slots; i++ {
		slotX := startX + float32(i)*float32(slotImage.Width*int32(cam.Zoom)+int32(spacing))
		(*inventory)[*slotNum] = newInventorySlot(slotX, startY, *slotNum)
		*slotNum++
	}
}

func loadInventory() {
	slotImage = rl.LoadTexture("assets/images/gui/slot.png")
	inventory = make([]InventorySlot, 3)
	hotInventory = make([]InventorySlot, 32)
	textures = make([]rl.Texture2D, 3)
	otherTextures = make([]rl.Texture2D, 32)
	wood = rl.LoadTexture("assets/images/items/wood.png")
	stone = rl.LoadTexture("assets/images/items/stone.png")
	metal = rl.LoadTexture("assets/images/items/metal.png")
	textures = []rl.Texture2D{wood, stone, metal}
	otherTextures = []rl.Texture2D{wall, floor}
	question = rl.LoadTexture("assets/images/gui/question.png")

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
		slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*cam.Zoom
		slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*cam.Zoom
		itemPosX := slotCenterX - (float32(otherTextures[slot].Width)/2)*cam.Zoom
		itemPosY := slotCenterY - (float32(otherTextures[slot].Height)/2)*cam.Zoom
		rl.DrawTextureEx(otherTextures[slot], rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
	} else {
		slotCenterX := hotInventory[slot].x + (float32(slotImage.Width)/2)*cam.Zoom
		slotCenterY := hotInventory[slot].y + (float32(slotImage.Height)/2)*cam.Zoom
		itemPosX := slotCenterX - (float32(textures[slot].Width)/2)*cam.Zoom
		itemPosY := slotCenterY - (float32(textures[slot].Height)/2)*cam.Zoom
		rl.DrawTextureEx(question, rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
	}
}

func drawItems() {
	for i, slot := range inventory {
		slotCenterX := slot.x + (float32(slotImage.Width)/2)*cam.Zoom
		slotCenterY := slot.y + (float32(slotImage.Height)/2)*cam.Zoom
		if i < len(textures) {
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*cam.Zoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*cam.Zoom
			rl.DrawTextureEx(textures[i], rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)

			var itemCount int
			switch i {
			case 0:
				itemCount = woodCount
			case 1:
				itemCount = stoneCount
			case 2:
				itemCount = metalCount
			}
			itemCountText := fmt.Sprintf("%d", itemCount)
			fontSize := float32(16)
			fontSpacing := float32(2)
			var itemCountTextWidth float32
			for {
				itemCountTextWidth = rl.MeasureTextEx(font, itemCountText, fontSize, fontSpacing).X
				if itemCountTextWidth+14 <= float32(slotImage.Width)*cam.Zoom {
					break
				}
				fontSize -= 1
				if fontSize <= 8 {
					break
				}
			}
			itemCountPosX := slot.x + float32(slotImage.Width)*cam.Zoom - itemCountTextWidth - 7
			itemCountPosY := slot.y + float32(slotImage.Height)*cam.Zoom - fontSize - 4
			rl.DrawTextEx(font, itemCountText, rl.NewVector2(itemCountPosX, itemCountPosY), fontSize, fontSpacing, rl.White)
		}
		if wallIsOpen {
			drawSlot(0, true)
		} else {
			drawSlot(0, false)
		}
		if floorIsOpen {
			drawSlot(1, true)
		} else {
			drawSlot(1, false)
		}
	}
}
