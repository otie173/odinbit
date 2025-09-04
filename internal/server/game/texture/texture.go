package texture

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	counter int = 0
)

type Texture struct {
	Id   int
	Path string
}

type Storage struct {
	textures map[string]Texture
}

func NewStorage() *Storage {
	textures := make(map[string]Texture, 128)

	return &Storage{
		textures: textures,
	}
}

func (s *Storage) LoadTextures() {
	s.textures = make(map[string]Texture, 128)

	files, err := filepath.Glob("resources/textures/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".png")
		path := file

		texture := Texture{Id: counter, Path: path}
		s.textures[name] = texture
		counter++
	}
}

func (s *Storage) GetID(name string) int {
	val, ok := s.textures[name]
	if !ok {
		return -1
	}
	return val.Id
}

func (s *Storage) GetTextures() ([]byte, error) {
	binaryTextures, err := msgpack.Marshal(s.textures)
	if err != nil {
		return nil, err
	}
	return binaryTextures, nil
}
