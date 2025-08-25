package world

import (
	"log"

	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/vmihailenco/msgpack/v5"
)

type Block struct {
	TextureID int
	Passable  bool
}

type World struct {
	blockStorage   *storage
	textureStorage *texture.Storage
}

func New(textureStorage *texture.Storage) *World {
	blockStorage := newStorage()

	return &World{
		blockStorage:   blockStorage,
		textureStorage: textureStorage,
	}
}

func (w *World) LoadStorage(data []byte) error {
	if err := msgpack.Unmarshal(data, &w.blockStorage); err != nil {
		return err
	}
	log.Println(w.blockStorage)
	return nil
}
