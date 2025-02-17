package entity

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	LEFT Direction = iota
	RIGHT
)

type Direction int

type Entity struct {
	HP  int        `json:"hp"`
	Pos rl.Vector2 `json:"pos"`
	Dir Direction
}
