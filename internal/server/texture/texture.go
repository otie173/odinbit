package texture

import (
	"log"
	"path/filepath"
	"strings"
)

var (
	counter int = 1
	storage map[string]Texture
)

type Texture struct {
	id   int
	path string
}

func loadTexture(name, path string) {
	texture := Texture{id: counter, path: path}
	storage[name] = texture
	counter++
}

func LoadTextures() {
	storage = make(map[string]Texture, 128)

	files, err := filepath.Glob("resources/textures/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".png")
		path := filepath.Base(file)

		loadTexture(name, path)
	}
}

func GetID(name string) int {
	val, ok := storage[name]
	if !ok {
		return -1
	}
	return val.id
}
