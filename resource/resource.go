package resource

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	TextureID      map[int]rl.Texture2D
	TextureName    map[string]int
	textureCounter int = 1

	BarrierTexture rl.Texture2D
	PlayerTexture  rl.Texture2D

	GrassTextures    []rl.Texture2D
	TreeTextures     []rl.Texture2D
	StoneTextures    []rl.Texture2D
	MushroomTextures []rl.Texture2D
)

const (
	grassTextures    int = 8
	treeTextures     int = 4
	stoneTextures    int = 3
	mushroomTextures int = 2
)

func Load() {
	TextureID = make(map[int]rl.Texture2D, 256)
	TextureName = make(map[string]int, 256)

	BarrierTexture = loadTexture("barrier.png")
	PlayerTexture = loadTexture("player.png")

	GrassTextures = make([]rl.Texture2D, 0, grassTextures)
	TreeTextures = make([]rl.Texture2D, 0, treeTextures)
	StoneTextures = make([]rl.Texture2D, 0, stoneTextures)
	MushroomTextures = make([]rl.Texture2D, 0, mushroomTextures)

	loadTextureBase("grass", grassTextures, &GrassTextures)
	loadTextureBase("tree", treeTextures, &TreeTextures)
	loadTextureBase("stone", stoneTextures, &StoneTextures)
	loadTextureBase("mushroom", mushroomTextures, &MushroomTextures)
}

func loadTexture(textureName string) rl.Texture2D {
	path := fmt.Sprintf("resource/sprite/%s", textureName)
	loadedTexture := rl.LoadTexture(path)

	TextureID[textureCounter] = loadedTexture
	TextureName[textureName] = textureCounter
	textureCounter++

	return loadedTexture
}

func loadTextureBase(baseName string, count int, target *[]rl.Texture2D) {
	for i := 1; i <= count; i++ {
		textureName := fmt.Sprintf("%s%d.png", baseName, i)
		path := fmt.Sprintf("resource/sprite/%s", textureName)
		loadedTexture := rl.LoadTexture(path)
		*target = append(*target, loadedTexture)

		TextureID[textureCounter] = loadedTexture
		TextureName[textureName] = textureCounter
		textureCounter++
	}
}
