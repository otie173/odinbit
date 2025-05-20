package world

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/common"
	"github.com/otie173/odinbit/internal/player"
	"github.com/otie173/odinbit/resource"
)

var (
	world [512][512]block
)

const (
	grassMultiplier     = 50
	tallGrassMultiplier = 3
	treeMultiplier      = 20
	deadTreeMultiplier  = 1
	stoneMultiplier     = 5
	mushroomMultiplier  = 2
)

type block struct {
	textureID int
	passable  bool
}

func Load() {
	GenerateWorld()

	schema := [][]string{
		{"tree4.png", "tree4.png", "tree4.png"},
		{"tree4.png", "stone3.png", "tree4.png"},
		{"tree4.png", "tree4.png", "tree4.png"},
	}
	generateSchema(250, 250, schema)
}

func Draw() {
	playerTileX := int(player.Player.Pos.X) / common.TileSize
	playerTileY := int(player.Player.Pos.Y) / common.TileSize

	for i := playerTileX - common.BlocksView; i <= playerTileX+common.BlocksView; i++ {
		for j := playerTileY - common.BlocksView; j <= playerTileY+common.BlocksView; j++ {
			if i >= 0 && i < common.WorldSize && j >= 0 && j < common.WorldSize && world[i][j].textureID != 0 {
				rl.DrawTexture(resource.TextureID[world[i][j].textureID], int32(i*common.TileSize), int32(j*common.TileSize), rl.White)
			}
		}
	}
}

func BlockExists(x, y int) bool {
	return world[x][y].textureID != 0
}

func IsValidPos(x, y int) bool {
	return x >= 0 && x < common.WorldSize && y >= 0 && y < common.WorldSize
}

func IsPassable(x, y int) bool {
	return world[x][y].passable
}

func AddBlock(id int, passable bool, x, y int) {
	world[x][y].textureID = id
	world[x][y].passable = passable
}

func DeleteBlock(x, y int) {
	world[x][y].textureID = 0
}

func generateBarrier() {
	for i := 0; i < common.WorldSize; i++ {
		AddBlock(resource.TextureName["barrier.png"], false, i, 0)
	}
	for j := 0; j < common.WorldSize; j++ {
		AddBlock(resource.TextureName["barrier.png"], false, 0, j)
	}

	for i := 0; i < common.WorldSize; i++ {
		AddBlock(resource.TextureName["barrier.png"], false, i, common.WorldSize-1)
	}
	for j := 0; j < common.WorldSize; j++ {
		AddBlock(resource.TextureName["barrier.png"], false, common.WorldSize-1, j)
	}
}

func generateResource(resourceName string, multiplier int, passable bool, textures ...int) {
	count := common.WorldSize * multiplier

	for i := 0; i <= count; i++ {
		var textureIndex int
		var textureName string

		textureIndex = rand.IntN(len(textures))
		textureName = fmt.Sprintf("%s%d.png", resourceName, textures[textureIndex])

		x := rand.IntN(common.WorldSize)
		y := rand.IntN(common.WorldSize)
		AddBlock(resource.TextureName[textureName], passable, x, y)
	}
}

func generateSchema(x, y int, schema [][]string) {
	rows := len(schema)
	cols := len(schema[0])

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			AddBlock(resource.TextureName[schema[i][j]], false, x+i, y+j)
		}
	}
}

func GenerateWorld() {
	generateResource("grass", grassMultiplier, true, 1, 2, 7, 8)
	generateResource("grass", tallGrassMultiplier, false, 3, 4, 5, 6)
	generateResource("tree", treeMultiplier, false, 4)
	generateResource("tree", deadTreeMultiplier, false, 2, 3)
	generateResource("stone", stoneMultiplier, false, 2, 3)
	generateResource("mushroom", mushroomMultiplier, false, 1, 2)
	generateBarrier()
}
