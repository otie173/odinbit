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
	floorIsOpen   bool
	question      rl.Texture2D
)

type InventorySlot struct {
	x          float32
	y          float32
	slotNumber int
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
	firstSlot := startX1 - 200
	for i := 0; i < 3; i++ {
		slotX := startX1 + float32(i)*float32((slotImage.Width*int32(cam.Zoom))+110)
		inventory[inventorySlotNum] = newInventorySlot(slotX, 182.0, inventorySlotNum)
		inventorySlotNum++
	}
	for i := 0; i < 8; i++ {
		slotX := firstSlot + float32(i)*float32((slotImage.Width*int32(cam.Zoom))+32)
		hotInventory[hotInventorySlotNum] = newInventorySlot(slotX, 395.0, hotInventorySlotNum)
		hotInventorySlotNum++
	}
	for i := 0; i < 8; i++ {
		slotX := firstSlot + float32(i)*float32((slotImage.Width*int32(cam.Zoom))+32)
		hotInventory[hotInventorySlotNum] = newInventorySlot(slotX, 495.0, hotInventorySlotNum)
		hotInventorySlotNum++
	}
	for i := 0; i < 8; i++ {
		slotX := firstSlot + float32(i)*float32((slotImage.Width*int32(cam.Zoom))+32)
		hotInventory[hotInventorySlotNum] = newInventorySlot(slotX, 595.0, hotInventorySlotNum)
		hotInventorySlotNum++
	}
	for i := 0; i < 8; i++ {
		slotX := firstSlot + float32(i)*float32((slotImage.Width*int32(cam.Zoom))+32)
		hotInventory[hotInventorySlotNum] = newInventorySlot(slotX, 695.0, hotInventorySlotNum)
		hotInventorySlotNum++
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
			slotCenterX := hotInventory[0].x + (float32(slotImage.Width)/2)*cam.Zoom
			slotCenterY := hotInventory[0].y + (float32(slotImage.Height)/2)*cam.Zoom
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*cam.Zoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*cam.Zoom
			rl.DrawTextureEx(wall, rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
		} else {
			slotCenterX := hotInventory[0].x + (float32(slotImage.Width)/2)*cam.Zoom
			slotCenterY := hotInventory[0].y + (float32(slotImage.Height)/2)*cam.Zoom
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*cam.Zoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*cam.Zoom
			rl.DrawTextureEx(question, rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
		}
		if floorIsOpen {
			slotCenterX := hotInventory[1].x + (float32(slotImage.Width)/2)*cam.Zoom
			slotCenterY := hotInventory[1].y + (float32(slotImage.Height)/2)*cam.Zoom
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*cam.Zoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*cam.Zoom
			rl.DrawTextureEx(floor, rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
		} else {
			slotCenterX := hotInventory[1].x + (float32(slotImage.Width)/2)*cam.Zoom
			slotCenterY := hotInventory[1].y + (float32(slotImage.Height)/2)*cam.Zoom
			itemPosX := slotCenterX - (float32(textures[i].Width)/2)*cam.Zoom
			itemPosY := slotCenterY - (float32(textures[i].Height)/2)*cam.Zoom
			rl.DrawTextureEx(question, rl.NewVector2(itemPosX, itemPosY), 0, cam.Zoom, rl.White)
		}
	}
}
