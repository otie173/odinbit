package world

import (
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/vmihailenco/msgpack/v5"
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

func (s *storage) getWorld() ([]byte, error) {
	data, err := msgpack.Marshal(&s.blocks)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *storage) getWorldArea(x, y int) ([]byte, error) {
	var blocks []Block

	for i := x - common.ViewRadius; i <= x+common.ViewRadius; i++ {
		for j := y - common.ViewRadius; j < y+common.ViewRadius; j++ {
			blocks = append(blocks, s.blocks[i][j])
		}
	}

	data, err := msgpack.Marshal(&blocks)
	if err != nil {
		return nil, err
	}
	return data, nil
}
