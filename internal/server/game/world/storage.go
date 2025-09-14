package world

import (
	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/vmihailenco/msgpack/v5"
)

type storage struct {
	blocks [common.WorldSize][common.WorldSize]Block
}

func newStorage() *storage {
	return &storage{}
}

func (s *storage) addBlock(id uint8, passable bool, x, y int16) {
	block := Block{TextureID: id, Passable: passable}
	s.blocks[x][y] = block
}

func (s *storage) removeBlock(x, y int16) {
	block := Block{TextureID: 0, Passable: true}
	s.blocks[x][y] = block
}

func (s *storage) getWorld() ([]byte, error) {
	data, err := msgpack.Marshal(&s.blocks)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *storage) getWorldArea(x, y int16) ([]byte, error) {
	var blocks []Block

	for i := x - common.ViewRadius; i <= x+common.ViewRadius; i++ {
		for j := y - common.ViewRadius; j < y+common.ViewRadius; j++ {
			blocks = append(blocks, s.blocks[i][j])
		}
	}

	data, err := binary.Marshal(&blocks)
	if err != nil {
		return nil, err
	}

	// data, err := msgpack.Marshal(&blocks)
	// if err != nil {
	// 	return nil, err
	// }
	return data, nil
}
