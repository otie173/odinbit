package world

import (
	"odinbit/resource"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	world map[rl.Vector2]Block = make(map[rl.Vector2]Block, WorldSize*WorldSize)

	Barrier                                                rl.Texture2D
	Tree1, Tree2                                           rl.Texture2D
	Stone1, Stone2, Stone3, Stone4, Stone5, Stone6, Stone7 rl.Texture2D
)

const (
	TileSize  float32 = 10
	WorldSize int     = 16
)

type Block struct {
	Rec      rl.Rectangle
	Texture  rl.Texture2D
	Passable bool
}

func LoadTexture() {
	Barrier = resource.LoadTexture("block/barrier.png")

	Tree1 = resource.LoadTexture("block/tree1.png")
	Tree2 = resource.LoadTexture("block/tree2.png")

	Stone1 = resource.LoadTexture("block/stone1.png")
	Stone2 = resource.LoadTexture("block/stone2.png")
	Stone3 = resource.LoadTexture("block/stone3.png")
	Stone4 = resource.LoadTexture("block/stone4.png")
	Stone5 = resource.LoadTexture("block/stone5.png")
	Stone6 = resource.LoadTexture("block/stone6.png")
	Stone7 = resource.LoadTexture("block/stone7.png")
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
