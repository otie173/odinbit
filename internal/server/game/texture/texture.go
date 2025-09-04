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

type TexturePack struct {
	textures map[string]Texture
}

func NewPack() *TexturePack {
	textures := make(map[string]Texture, 128)

	return &TexturePack{
		textures: textures,
	}
}

func (t *TexturePack) LoadTextures() {
	t.textures = make(map[string]Texture, 128)

	files, err := filepath.Glob("resources/textures/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".png")
		path := file

		texture := Texture{Id: counter, Path: path}
		t.textures[name] = texture
		counter++
	}
}

func (t *TexturePack) GetID(name string) int {
	val, ok := t.textures[name]
	if !ok {
		return -1
	}
	return val.Id
}

func (t *TexturePack) GetTextures() ([]byte, error) {
	binaryTextures, err := msgpack.Marshal(t.textures)
	if err != nil {
		return nil, err
	}
	return binaryTextures, nil
}
