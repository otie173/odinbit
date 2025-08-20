package world

import (
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/texture"
)

type Block struct {
	textureID int
	passable  bool
}

type World struct {
	storage   *storage
	generator *generator
}

func New(textureStorage *texture.Storage) *World {
	blockStorage := newStorage()
	generator := newGenerator(textureStorage, blockStorage)

	return &World{
		storage:   blockStorage,
		generator: generator,
	}
}

func (w *World) AddBlock(id int, passable bool, x, y int) {
	w.storage.addBlock(id, passable, x, y)
}

func (w *World) RemoveBlock(x, y int) {
	w.storage.removeBlock(x, y)
}

func (w *World) Generate() {
	w.generator.generateWorld()
}

func (w *World) GetWorld() [common.WorldSize][common.WorldSize]Block {
	return w.storage.blocks
}
