package texture

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Storage struct {
	texture map[string]uint8
	id      map[uint8]rl.Texture2D
}

func New() *Storage {
	texture := make(map[string]uint8, 128)
	id := make(map[uint8]rl.Texture2D)

	return &Storage{
		texture: texture,
		id:      id,
	}
}

func (s *Storage) LoadTexture(id uint8, path string) {
	texture := rl.LoadTexture(path)
	textureName := strings.TrimSuffix(filepath.Base(path), ".png")

	s.texture[textureName] = id
	s.id[id] = texture
}

func (s *Storage) GetById(blockId uint8) rl.Texture2D {
	return s.id[blockId]
}
