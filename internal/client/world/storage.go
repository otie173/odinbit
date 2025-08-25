package world

import (
	"github.com/otie173/odinbit/internal/server/common"
)

type storage struct {
	blocks [common.WorldSize][common.WorldSize]Block
}

func newStorage() *storage {
	return &storage{}
}

func (s *storage) addBlock(id int, passable bool, x, y int) {
	block := Block{TextureID: id, Passable: passable}
	s.blocks[x][y] = block
}

func (s *storage) removeBlock(x, y int) {
	block := Block{TextureID: -1, Passable: true}
	s.blocks[x][y] = block
}
