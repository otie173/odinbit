package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	blockX     float32 = 10
	blockY     float32 = 10
	blockScale float32 = 8
	heart      rl.Texture2D
	heartX     float32 = 105
	heartY     float32 = 10
	heartScale float32 = 4
	pickaxe    rl.Texture2D
	axe        rl.Texture2D
	shovel     rl.Texture2D
	toolX      float32 = 105
	toolY      float32 = 50
	toolScale  float32 = 4
)

func loadUI() {
	heart = loadTexture("assets/images/gui/heart.png")
	pickaxe = loadTexture("assets/images/gui/pickaxe.png")
	axe = loadTexture("assets/images/gui/axe.png")
	shovel = loadTexture("assets/images/gui/shovel.png")
}

func drawUI() {
	// Keep It Simple Stupid :)
	itemStates := map[int]bool{
		WALL:        wallIsOpen,
		FLOOR:       floorIsOpen,
		DOOR:        doorIsOpen,
		CHEST:       chestIsOpen,
		WALLWINDOW:  wallWindowIsOpen,
		DOOROPEN:    doorOpenIsOpen,
		BIGBARREL:   bigBarrelIsOpen,
		BOOKSHELF:   bookshelfIsOpen,
		CHAIR:       chairIsOpen,
		CLOSET:      closetIsOpen,
		FENCE1:      fence1IsOpen,
		FENCE2:      fence2IsOpen,
		FLOOR2:      floor2IsOpen,
		FLOOR4:      floor4IsOpen,
		LAMP:        lampIsOpen,
		LOOTBOX:     lootboxIsOpen,
		SHELF:       shelfIsOpen,
		SIGN:        signIsOpen,
		SMALLBARREL: smallBarrelIsOpen,
		TABLE:       tableIsOpen,
		TOMBSTONE:   tombstoneIsOpen,
		TRASH:       trashIsOpen,
		SAPLING:     saplingIsOpen,
		SEED2SMALL:  seedIsOpen,
		SEED1BIG:    cabbageIsOpen,
	}

	if state, ok := itemStates[item]; ok && state {
		rl.DrawTextureEx(id[item], rl.NewVector2(blockX, blockY), 0, blockScale, rl.White)
	} else {
		rl.DrawTextureEx(question, rl.NewVector2(blockX, blockY), 0, blockScale, rl.White)
	}

	switch playerHealth {
	case 3:
		rl.DrawTextureEx(heart, rl.NewVector2(heartX, heartY), 0, heartScale, rl.White)
		rl.DrawTextureEx(heart, rl.NewVector2(heartX+40, heartY), 0, heartScale, rl.White)
		rl.DrawTextureEx(heart, rl.NewVector2(heartX+80, heartY), 0, heartScale, rl.White)
	case 2:
		rl.DrawTextureEx(heart, rl.NewVector2(heartX, heartY), 0, heartScale, rl.White)
		rl.DrawTextureEx(heart, rl.NewVector2(heartX+40, heartY), 0, heartScale, rl.White)
	case 1:
		rl.DrawTextureEx(heart, rl.NewVector2(heartX, heartY), 0, heartScale, rl.White)
	}

	if pickaxeIsOpen {
		rl.DrawTextureEx(pickaxe, rl.NewVector2(toolX, toolY), 0, toolScale, rl.White)
	} else {
		rl.DrawTextureEx(question, rl.NewVector2(toolX, toolY), 0, toolScale, rl.White)
	}

	if axeIsOpen {
		rl.DrawTextureEx(axe, rl.NewVector2(toolX+40, toolY), 0, toolScale, rl.White)
	} else {
		rl.DrawTextureEx(question, rl.NewVector2(toolX+40, toolY), 0, toolScale, rl.White)
	}

	if shovelIsOpen {
		rl.DrawTextureEx(shovel, rl.NewVector2(toolX+80, toolY), 0, toolScale, rl.White)
	} else {
		rl.DrawTextureEx(question, rl.NewVector2(toolX+80, toolY), 0, toolScale, rl.White)
	}
}
