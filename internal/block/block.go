package block

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Block struct {
	ID       int
	Texture  rl.Texture2D
	Passable bool
	Behavior func()
}
