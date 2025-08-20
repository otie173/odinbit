package world

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/texture"
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
	textureStorage *texture.Storage
	blockStorage   *storage
}

func newGenerator(textureStorage *texture.Storage, blockStorage *storage) *generator {
	return &generator{
		textureStorage: textureStorage,
		blockStorage:   blockStorage,
	}
}

func (g *generator) generateBarrier() {
	for i := 0; i < common.WorldSize; i++ {
		g.blockStorage.blocks[i][0] = Block{textureID: g.textureStorage.GetID("barrier")}
	}
	for j := 0; j < common.WorldSize; j++ {
		g.blockStorage.blocks[0][j] = Block{textureID: g.textureStorage.GetID("barrier")}
	}

	for i := 0; i < common.WorldSize; i++ {
		g.blockStorage.blocks[i][common.WorldSize-1] = Block{textureID: g.textureStorage.GetID("barrier")}
	}
	for j := 0; j < common.WorldSize; j++ {
		g.blockStorage.blocks[common.WorldSize-1][j] = Block{textureID: g.textureStorage.GetID("barrier")}
	}
}

func (g *generator) generateResource(name string, multiplier float32, passable bool, textures ...int) {
	count := common.WorldSize * multiplier

	for i := 0; i <= int(count); i++ {
		var textureIndex int
		var textureName string

		textureIndex = rand.IntN(len(textures))
		textureName = fmt.Sprintf("%s%d.png", name, textures[textureIndex])

		x := rand.IntN(common.WorldSize)
		y := rand.IntN(common.WorldSize)
		g.blockStorage.blocks[x][y] = Block{
			textureID: g.textureStorage.GetID(strings.TrimSuffix(textureName, ".png")),
			passable:  passable,
		}
	}
}

func (g *generator) generateWorld() {
	g.generateResource("grass", grassMultiplier, true, 1, 2, 7, 8)
	g.generateResource("grass", tallGrassMultiplier, false, 3, 4, 5, 6)
	g.generateResource("tree", treeMultiplier, false, 4)
	g.generateResource("tree", deadTreeMultiplier, false, 2, 3)
	g.generateResource("stone", stoneMultiplier, false, 2, 3)
	g.generateResource("mushroom", mushroomMultiplier, false, 1, 2)
	g.generateBarrier()
}
