package world

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/block"
	"github.com/otie173/odinbit/internal/common"
	"github.com/otie173/odinbit/internal/player"
)

var (
	world [512][512]block.Block
)

const (
	grassMultiplier     = 50
	tallGrassMultiplier = 3
	treeMultiplier      = 20
	deadTreeMultiplier  = 1
	stoneMultiplier     = 5
	mushroomMultiplier  = 2
)

func Load() {
	GenerateWorld()

	schema := [][]block.Block{
		{block.Tree4, block.Tree4, block.Tree4},
		{block.Tree4, block.Stone3, block.Tree4},
		{block.Tree4, block.Tree4, block.Tree4},
	}
	generateSchema(250, 250, schema)
}

func Draw() {
	playerTileX := int(player.Player.Pos.X) / common.TileSize
	playerTileY := int(player.Player.Pos.Y) / common.TileSize

	for i := playerTileX - common.BlocksView; i <= playerTileX+common.BlocksView; i++ {
		for j := playerTileY - common.BlocksView; j <= playerTileY+common.BlocksView; j++ {
			if i >= 0 && i < common.WorldSize && j >= 0 && j < common.WorldSize && world[i][j].ID != 0 {
				rl.DrawTexture(world[i][j].Texture, int32(i*common.TileSize), int32(j*common.TileSize), rl.White)
			}
		}
	}
}

func BlockExists(x, y int) bool {
	return world[x][y].ID != 0
}

func IsValidPos(x, y int) bool {
	return x >= 0 && x < common.WorldSize && y >= 0 && y < common.WorldSize
}

func IsPassable(x, y int) bool {
	return world[x][y].Passable
}

func AddBlock(tile block.Block, x, y int) {
	world[x][y] = tile
}

func DeleteBlock(x, y int) {
	world[x][y].ID = 0
}

func generateBarrier() {
	for i := 0; i < common.WorldSize; i++ {
		AddBlock(block.Barrier, i, 0)
	}
	for j := 0; j < common.WorldSize; j++ {
		AddBlock(block.Barrier, 0, j)
	}

	for i := 0; i < common.WorldSize; i++ {
		AddBlock(block.Barrier, i, common.WorldSize-1)
	}
	for j := 0; j < common.WorldSize; j++ {
		AddBlock(block.Barrier, common.WorldSize-1, j)
	}
}

func generateResource(multiplier int, blocks []block.Block) {
	count := common.WorldSize * multiplier

	for i := 0; i <= count; i++ {
		var block int = rand.IntN(len(blocks))

		x := rand.IntN(common.WorldSize)
		y := rand.IntN(common.WorldSize)
		AddBlock(blocks[block], x, y)
	}
}

func generateSchema(x, y int, schema [][]block.Block) {
	rows := len(schema)
	cols := len(schema[0])

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			AddBlock(schema[i][j], x+i, y+j)
		}
	}
}

func CheckBehavior(x, y int) {
	world[x][y].Behavior()
}

func GenerateWorld() {
	generateResource(grassMultiplier, []block.Block{block.Grass1, block.Grass2, block.Grass7, block.Grass8})
	generateResource(tallGrassMultiplier, []block.Block{block.Grass3, block.Grass4, block.Grass5, block.Grass6})
	generateResource(treeMultiplier, []block.Block{block.Tree4})
	generateResource(deadTreeMultiplier, []block.Block{block.Tree2, block.Tree3})
	generateResource(stoneMultiplier, []block.Block{block.Stone2, block.Stone3})
	generateResource(mushroomMultiplier, []block.Block{block.Mushroom1, block.Mushroom2})
	generateBarrier()

	//generateResource("grass", grassMultiplier, true, 1, 2, 7, 8)
	//generateResource("grass", tallGrassMultiplier, false, 3, 4, 5, 6)
	//generateResource("tree", treeMultiplier, false, 4)
	//generateResource("tree", deadTreeMultiplier, false, 2, 3)
	//generateResource("stone", stoneMultiplier, false, 2, 3)
	//generateResource("mushroom", mushroomMultiplier, false, 1, 2)
	//generateBarrier()
}
