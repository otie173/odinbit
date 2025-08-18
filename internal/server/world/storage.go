package world

import "github.com/otie173/odinbit/internal/server/common"

type storage struct {
	blocks [common.WorldSize][common.WorldSize]block
}

func newStorage() *storage {
	return &storage{}
}

func (s *storage) addBlock(id int, passable bool, x, y int) {
	block := block{textureID: id, passable: passable}
	s.blocks[x][y] = block
}

func (s *storage) removeBlock(x, y int) {
	block := block{textureID: -1, passable: true}
	s.blocks[x][y] = block
}
