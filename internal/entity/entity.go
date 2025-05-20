package entity

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	Right byte = iota
	Left
)

type Entity struct {
	Texture rl.Texture2D
	Pos     rl.Vector2
}
