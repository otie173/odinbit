package world

import (
	"fmt"
	"log"
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

type block struct {
	textureID int
	passable  bool
}

type World struct {
	blocks [common.WorldSize][common.WorldSize]block
}

func New() *World {
	return &World{}
}

func (w *World) AddBlock(id int, passable bool, x, y int) {
	block := block{textureID: id, passable: passable}
	w.blocks[x][y] = block
}

func (w *World) RemoveBlock() {

}

func generateBarrier(w *World) {
	for i := 0; i < common.WorldSize; i++ {
		w.AddBlock(texture.GetID("barrier"), false, i, 0)
	}
	for j := 0; j < common.WorldSize; j++ {
		w.AddBlock(texture.GetID("barrier"), false, 0, j)
	}

	for i := 0; i < common.WorldSize; i++ {
		w.AddBlock(texture.GetID("barrier"), false, i, common.WorldSize-1)
	}
	for j := 0; j < common.WorldSize; j++ {
		w.AddBlock(texture.GetID("barrier"), false, common.WorldSize-1, j)
	}
}

func generateResource(w *World, name string, multiplier float32, passable bool, textures ...int) {
	count := common.WorldSize * multiplier

	for i := 0; i <= int(count); i++ {
		var textureIndex int
		var textureName string

		textureIndex = rand.IntN(len(textures))
		textureName = fmt.Sprintf("%s%d.png", name, textures[textureIndex])

		x := rand.IntN(common.WorldSize)
		y := rand.IntN(common.WorldSize)
		w.AddBlock(texture.GetID(strings.TrimSuffix(textureName, ".png")), passable, x, y)
	}
}

func (w *World) Generate() {
	generateResource(w, "grass", grassMultiplier, true, 1, 2, 7, 8)
	generateResource(w, "grass", tallGrassMultiplier, false, 3, 4, 5, 6)
	generateResource(w, "tree", treeMultiplier, false, 4)
	generateResource(w, "tree", deadTreeMultiplier, false, 2, 3)
	generateResource(w, "stone", stoneMultiplier, false, 2, 3)
	generateResource(w, "mushroom", mushroomMultiplier, false, 1, 2)
	generateBarrier(w)
	log.Println(w.blocks)
}
