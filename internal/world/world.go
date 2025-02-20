package world

import (
	"odinbit/resource"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	world map[rl.Vector2]Block = make(map[rl.Vector2]Block, WorldSize*WorldSize)

	Tree1, Tree2 rl.Texture2D
)

const (
	TileSize  float32 = 10
	WorldSize int     = 512
)

type Block struct {
	Rec      rl.Rectangle
	Texture  rl.Texture2D
	Passable bool
}

func LoadTexture() {
	Tree1 = resource.LoadTexture("block/tree1.png")
	Tree2 = resource.LoadTexture("block/tree2.png")
}

func AddBlock(x, y float32, texture rl.Texture2D) {
	world[rl.NewVector2(x*TileSize, y*TileSize)] = Block{Rec: rl.NewRectangle(x*TileSize, y*TileSize, TileSize, TileSize), Texture: texture, Passable: false}
}

func RemoveBlock(x, y float32) {
	delete(world, rl.NewVector2(x*TileSize, y*TileSize))
}

func DrawWorld() {
	for _, block := range world {
		rl.DrawTextureEx(block.Texture, rl.NewVector2(block.Rec.X, block.Rec.Y), 0.0, 1.0, rl.White)
	}
}
