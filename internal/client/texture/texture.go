package texture

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Storage struct {
	texture map[string]int
	id      map[int]rl.Texture2D
}

func New() *Storage {
	texture := make(map[string]int, 128)
	id := make(map[int]rl.Texture2D)

	return &Storage{
		texture: texture,
		id:      id,
	}
}

func (s *Storage) LoadTexture(id int, path string) {
	texture := rl.LoadTexture(path)
	textureName := strings.TrimSuffix(filepath.Base(path), ".png")

	s.texture[textureName] = id
	s.id[id] = texture
}

func (s *Storage) GetById(blockId int) rl.Texture2D {
	return s.id[blockId]
}
