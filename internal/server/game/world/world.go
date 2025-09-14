package world

import (
	"github.com/otie173/odinbit/internal/server/game/texture"
)

type Block struct {
	TextureID int
	Passable  bool
	X         int
	Y         int
}

type World struct {
	storage   *storage
	generator *generator
}

func New(textures *texture.TexturePack) *World {
	blockStorage := newStorage()
	generator := newGenerator(textures, blockStorage)

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

func (w *World) GetWorldArea(x, y int) ([]byte, error) {
	binaryWorldArea, err := w.storage.getWorldArea(x, y)
	if err != nil {
		return nil, err
	}
	return binaryWorldArea, nil
}
