package world

import (
	"github.com/otie173/odinbit/internal/server/game/texture"
)

type Block struct {
	TextureID int
	Passable  bool
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

func (w *World) GetWorld() ([]byte, error) {
	binaryWorld, err := w.storage.getWorld()
	if err != nil {
		return nil, err
	}
	return binaryWorld, nil
}
