package resource

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	BarrierTexture   rl.Texture2D
	PlayerTexture    rl.Texture2D
	WallTexture      rl.Texture2D
	WindowTexture    rl.Texture2D
	Door1Texture     rl.Texture2D
	Door2Texture     rl.Texture2D
	Floor1Texture    rl.Texture2D
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
	BarrierTexture = loadTexture("barrier.png")
	PlayerTexture = loadTexture("player.png")
	WallTexture = loadTexture("wall.png")
	WindowTexture = loadTexture("window.png")
	Door1Texture = loadTexture("door1.png")
	Door2Texture = loadTexture("door2.png")
	Floor1Texture = loadTexture("floor1.png")
	log.Println(Floor1Texture)

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

	return loadedTexture
}

func loadTextureBase(baseName string, count int, target *[]rl.Texture2D) {
	for i := 1; i <= count; i++ {
		textureName := fmt.Sprintf("%s%d.png", baseName, i)
		path := fmt.Sprintf("resource/sprite/%s", textureName)
		loadedTexture := rl.LoadTexture(path)
		*target = append(*target, loadedTexture)
	}
}
