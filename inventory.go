package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	inventoryOpen bool
	slotImage     rl.Texture2D
)

func loadInventory() {
	slotImage = rl.LoadTexture("assets/images/gui/slot.png")
}

func unloadInventory() {
	rl.UnloadTexture(slotImage)
}
