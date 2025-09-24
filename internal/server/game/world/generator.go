package world

import (
	"fmt"
	"math/rand/v2"

	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/game/texture"
)

const (
	grassMultiplier     float32 = 50
	tallGrassMultiplier float32 = 3
	treeMultiplier      float32 = 20
	deadTreeMultiplier  float32 = 1
	stoneMultiplier     float32 = 5
	mushroomMultiplier  float32 = 2
)

type generator struct {
	textures     *texture.TexturePack
	blockStorage *storage
}

func newGenerator(textures *texture.TexturePack, blockStorage *storage) *generator {
	return &generator{
		textures:     textures,
		blockStorage: blockStorage,
	}
}

func (g *generator) generateBarrier() {
	for i := range common.WorldSize {
		g.blockStorage.blocks[i][0] = Block{TextureID: g.textures.GetID("barrier")}
	}
	for j := range common.WorldSize {
		g.blockStorage.blocks[0][j] = Block{TextureID: g.textures.GetID("barrier")}
	}

	for i := range common.WorldSize {
		g.blockStorage.blocks[i][common.WorldSize-1] = Block{TextureID: g.textures.GetID("barrier")}
	}
	for j := range common.WorldSize {
		g.blockStorage.blocks[common.WorldSize-1][j] = Block{TextureID: g.textures.GetID("barrier")}
	}
}

func (g *generator) generateResource(name string, multiplier float32, passable uint8, textures ...int) {
	count := float32(common.WorldSize) * multiplier

	for i := 0; i <= int(count); i++ {
		var textureIndex int
		var textureName string

		textureIndex = rand.IntN(len(textures))
		textureName = fmt.Sprintf("%s%d", name, textures[textureIndex])

		x := rand.IntN(common.WorldSize)
		y := rand.IntN(common.WorldSize)
		g.blockStorage.addBlock(g.textures.GetID(textureName), passable, int16(x), int16(y))
	}
}

func (g *generator) generateSchema(schema [][]string, startX, startY int16) {
	for y, row := range schema {
		for x, block := range row {
			if block == "" {
				continue
			}

			blockID := g.textures.GetID(block)
			blockX := startX + int16(x)
			blockY := startY + int16(y)

			if blockX < 0 || int(blockX) > common.WorldSize ||
				blockY < 0 || int(blockY) > common.WorldSize {
				continue
			}

			g.blockStorage.blocks[blockX][blockY] = Block{TextureID: blockID, Passable: 0}
		}
	}
}

func (g *generator) generateWorld() {
	g.generateResource("grass", grassMultiplier, 1, 1, 2, 7, 8)
	g.generateResource("grass", tallGrassMultiplier, 0, 3, 4, 5, 6)
	g.generateResource("tree", treeMultiplier, 0, 4)
	g.generateResource("tree", deadTreeMultiplier, 0, 2, 3)
	g.generateResource("stone", stoneMultiplier, 0, 2, 3)
	g.generateResource("mushroom", mushroomMultiplier, 0, 1, 2)

	// houseSchema := [][]string{
	// 	{"wall", "wall", "wall", "wall", "wall"},
	// 	{"", "", "", "", "wall"},
	// 	{"wall", "", "", "", "wall"},
	// 	{"wall", "wall", "wall", "wall", "wall"},
	// }
	// g.generateSchema(houseSchema, 250, 250)
	block := Block{TextureID: g.textures.GetID("home"), Passable: 1}
	g.blockStorage.addBlock(block.TextureID, block.Passable, 258, 259)

	g.generateBarrier()
}
