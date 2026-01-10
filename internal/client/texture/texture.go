package texture

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	PlayerTexture rl.Texture2D
	WoodMaterial  rl.Texture2D
	StoneMaterial rl.Texture2D
	MetalMaterial rl.Texture2D
)

type Texture struct {
	Id   uint8
	Path string
}

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

func (s *Storage) GetIdByName(blockName string) uint8 {
	return s.texture[blockName]
}
