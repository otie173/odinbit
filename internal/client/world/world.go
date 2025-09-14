package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/texture"
)

var (
	Overworld World
)

type World struct {
	Textures       *texture.Storage
	Blocks         []Block
	StartX, StartY int
	EndX, EndY     int
}

type Block struct {
	TextureID int
	Passable  bool
}

func GetBlock(id int) rl.Texture2D {
	return Overworld.Textures.GetById(id)
}
