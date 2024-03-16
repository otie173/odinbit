package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	world      map[rl.Rectangle]Block
	wall       rl.Texture2D
	floor      rl.Texture2D
	door       rl.Texture2D
	chest      rl.Texture2D
	smallTree  rl.Texture2D
	stone1     rl.Texture2D
	stone2     rl.Texture2D
	stone3     rl.Texture2D
	stone4     rl.Texture2D
	normalTree rl.Texture2D
	bigTree    rl.Texture2D
)

const (
	TILE_SIZE float32 = 10.0
)

type Block struct {
	img      rl.Texture2D
	rec      rl.Rectangle
	passable bool
}

func loadWorld() {
	world = make(map[rl.Rectangle]Block)
	wall = rl.LoadTexture("assets/images/blocks/wall.png")
	floor = rl.LoadTexture("assets/images/blocks/floor.png")
	door = rl.LoadTexture("assets/images/blocks/door.png")
	chest = rl.LoadTexture("assets/images/blocks/chest.png")
	smallTree = rl.LoadTexture("assets/images/world/small_tree.png")
	stone1 = rl.LoadTexture("assets/images/world/stone1.png")
	stone2 = rl.LoadTexture("assets/images/world/stone2.png")
	stone3 = rl.LoadTexture("assets/images/world/stone3.png")
	stone4 = rl.LoadTexture("assets/images/world/stone4.png")
	normalTree = rl.LoadTexture("assets/images/world/normal_tree.png")
	bigTree = rl.LoadTexture("assets/images/world/big_tree.png")

	rl.SetTextureFilter(smallTree, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone1, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone2, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone3, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone4, rl.TextureFilterNearest)
	rl.SetTextureFilter(normalTree, rl.TextureFilterNearest)
	rl.SetTextureFilter(bigTree, rl.TextureFilterNearest)
}

func unloadWorld() {
	rl.UnloadTexture(wall)
	rl.UnloadTexture(floor)
	rl.UnloadTexture(chest)
	rl.UnloadTexture(smallTree)
	rl.UnloadTexture(stone1)
	rl.UnloadTexture(stone2)
	rl.UnloadTexture(stone3)
	rl.UnloadTexture(stone4)
	rl.UnloadTexture(normalTree)
	rl.UnloadTexture(bigTree)
}

func addBlock(img rl.Texture2D, x, y float32, passable bool) {
	block := Block{
		img:      img,
		rec:      rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE),
		passable: passable,
	}
	world[block.rec] = block
}

func removeBlock(x, y float32) {
	delete(world, rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE))
}

func generateStructure(x, y, structure int) {
	switch structure {
	case 1:
		addBlock(wall, float32(x), float32(y), false)
		addBlock(wall, float32(x+1), float32(y), false)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(floor, float32(x+3), float32(y+2), true)
		addBlock(wall, float32(x+4), float32(y+1), false)
		addBlock(floor, float32(x+3), float32(y+4), true)
	case 2:
		addBlock(wall, float32(x), float32(y), false)
		addBlock(door, float32(x+1), float32(y), true)
		addBlock(wall, float32(x+2), float32(y), false)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+2), float32(y+1), true)
		addBlock(floor, float32(x), float32(y+2), true)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(wall, float32(x+1), float32(y+3), false)
		addBlock(floor, float32(x+2), float32(y+3), true)
	case 3:
		addBlock(floor, float32(x), float32(y), true)
		addBlock(floor, float32(x+1), float32(y), true)
		addBlock(floor, float32(x+3), float32(y), true)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+2), float32(y+1), true)
		addBlock(wall, float32(x+3), float32(y+1), false)
		addBlock(wall, float32(x), float32(y+2), false)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(door, float32(x+3), float32(y+2), true)
		addBlock(floor, float32(x+4), float32(y+2), true)
		addBlock(wall, float32(x+1), float32(y+3), false)
		addBlock(wall, float32(x+2), float32(y+3), false)
		addBlock(wall, float32(x+3), float32(y+3), false)

	}
}

func generateWorld() {
	// Генерация данжа1
	x1 := rand.Intn(65) - 32
	y1 := rand.Intn(65) - 32
	generateStructure(x1, y1, 1)

	//Генерация данжа2
	x2 := rand.Intn(65) - 32
	y2 := rand.Intn(65) - 32
	generateStructure(x2, y2, 2)

	// Генерация данжа3
	x3 := rand.Intn(65) - 32
	y3 := rand.Intn(65) - 32
	generateStructure(x3, y3, 3)

	// Генерация деревьев
	for i := 0; i < 96; i++ {
		x := rand.Intn(65) - 32
		y := rand.Intn(65) - 32

		treeImg := rand.Intn(3) + 1
		switch treeImg {
		case 1:
			addBlock(smallTree, float32(x), float32(y), false)
		case 2:
			addBlock(normalTree, float32(x), float32(y), false)
		case 3:
			addBlock(bigTree, float32(x), float32(y), false)
		}
	}

	// Генерация камней
	for i := 0; i < 96; i++ {
		x := rand.Intn(65) - 32
		y := rand.Intn(65) - 32

		stoneImg := rand.Intn(4) + 1
		switch stoneImg {
		case 1:
			addBlock(stone1, float32(x), float32(y), false)
		case 2:
			addBlock(stone2, float32(x), float32(y), false)
		case 3:
			addBlock(stone3, float32(x), float32(y), false)
		case 4:
			addBlock(stone4, float32(x), float32(y), false)
		}
	}
}

func drawWorld() {
	for _, block := range world {
		rl.DrawTextureRec(block.img, block.rec, rl.NewVector2(block.rec.X, block.rec.Y), rl.White)
	}
}
