package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	world map[rl.Rectangle]Block
	wall  rl.Texture2D
	floor rl.Texture2D
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
}

func unloadWorld() {
	rl.UnloadTexture(wall)
	rl.UnloadTexture(floor)
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

func drawWorld() {
	rl.DrawTextureRec(player, playerRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	for _, block := range world {
		rl.DrawTextureRec(block.img, block.rec, rl.NewVector2(block.rec.X, block.rec.Y), rl.White)
	}
}
