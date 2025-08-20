package texture

import rl "github.com/gen2brain/raylib-go/raylib"

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

func (s *Storage) LoadTexture() {

}
