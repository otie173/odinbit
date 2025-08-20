package texture

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type Texture struct {
	id   int
	path string
}

type Storage struct {
	counter int
	storage map[string]Texture
}

func New() *Storage {
	counter := 0
	storage := make(map[string]Texture, 128)

	return &Storage{
		counter: counter,
		storage: storage,
	}
}

func (s *Storage) LoadTextures() {
	s.storage = make(map[string]Texture, 128)

	files, err := filepath.Glob("resources/textures/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".png")
		path := filepath.Base(file)

		texture := Texture{id: s.counter, path: path}
		s.storage[name] = texture
		s.counter++
	}
}

func (s *Storage) GetID(name string) int {
	val, ok := s.storage[name]
	if !ok {
		return -1
	}
	return val.id
}

func (s *Storage) GetTextures() ([]byte, error) {
	binaryTextures, err := msgpack.Marshal(s.storage)
	if err != nil {
		return nil, err
	}
	return binaryTextures, nil
}
