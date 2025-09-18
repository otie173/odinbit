package world

import (
	"log"

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

func (s *storage) addBlock(id uint8, passable uint8, x, y int16) {
	block := Block{TextureID: id, Passable: passable}
	s.blocks[x][y] = block
}

func (s *storage) removeBlock(x, y int16) {
	block := Block{TextureID: 0, Passable: 1}
	s.blocks[x][y] = block
}

func (s *storage) getWorld() ([]byte, error) {
	data, err := msgpack.Marshal(&s.blocks)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *storage) getWorldArea(x, y float32) ([]byte, error) {
	startX := int(x - common.ViewRadius)
	endX := int(x + common.ViewRadius)
	startY := int(y - common.ViewRadius)
	endY := int(y + common.ViewRadius)

	areaWidth := endX - startX + 1
	areaHeight := endY - startY + 1
	blocks := make([]Block, 0, areaWidth*areaHeight)

	log.Println(startX, endX, startY, endY)
	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			blocks = append(blocks, s.blocks[i][j])
		}
	}

	data, err := binary.Marshal(&blocks)
	if err != nil {
		return nil, err
	}
	return data, nil
}
