package world

import (
	"math"

	"github.com/kelindar/binary"
	scommon "github.com/otie173/odinbit/internal/server/common"
	"github.com/vmihailenco/msgpack/v5"
)

type storage struct {
	blocks [scommon.WorldSize][scommon.WorldSize]Block
}

func newStorage() *storage {
	return &storage{}
}

func (s *storage) addBlock(id uint8, passable uint8, x, y int16) {
	block := Block{TextureID: id, Passable: passable}
	if s.blocks[x][y].TextureID == 0 || s.blocks[x][y].Passable == 1 {
		s.blocks[x][y] = block
	}
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

func computeAreaBounds(coord float32) (start, end int) {
	radius := float64(scommon.ViewRadius)
	start = int(math.Floor(float64(coord) - radius))
	end = int(math.Ceil(float64(coord) + radius))

	width := end - start
	if width <= 0 {
		width = 1
		end = start + width
	}

	if start < 0 {
		end += -start
		start = 0
	}

	if end > scommon.WorldSize {
		shift := end - scommon.WorldSize
		end = scommon.WorldSize
		start -= shift
		if start < 0 {
			start = 0
		}
	}

	if end > scommon.WorldSize {
		end = scommon.WorldSize
	}

	if end <= start {
		end = start + 1
		if end > scommon.WorldSize {
			end = scommon.WorldSize
		}
	}

	return start, end
}

func (s *storage) getWorldArea(x, y float32) ([]byte, scommon.AreaPositions, error) {
	startX, endX := computeAreaBounds(x)
	startY, endY := computeAreaBounds(y)

	areaWidth := endX - startX
	areaHeight := endY - startY
	blocks := make([]Block, 0, areaWidth*areaHeight)

	//log.Println(startX, endX, startY, endY)
	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			blocks = append(blocks, s.blocks[i][j])
		}
	}

	data, err := binary.Marshal(&blocks)
	if err != nil {
		return nil, scommon.AreaPositions{}, err
	}

	return data, scommon.AreaPositions{StartX: startX, StartY: startY, EndX: endX, EndY: endY}, nil
}
