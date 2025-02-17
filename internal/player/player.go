package player

import (
	"math"
	"odinbit/internal/entity"
	"odinbit/internal/world"
	"odinbit/resource"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	Roben                  Player       = Player{entity.Entity{HP: 3, Pos: rl.NewVector2(0.0, 0.0), Dir: entity.RIGHT}}
	playerRectangle        rl.Rectangle = rl.NewRectangle(Roben.Pos.X, Roben.Pos.Y, world.TileSize, world.TileSize)
	playerFlippedRectangle rl.Rectangle = rl.NewRectangle(Roben.Pos.X, Roben.Pos.Y, -world.TileSize, world.TileSize)
	Cam                    rl.Camera2D
	playerIdle             rl.Texture2D
)

type Player struct {
	entity.Entity
}

func LoadTexture() {
	playerIdle = resource.LoadTexture("entity/player_idle.png")
}

func RegisterCam() {
	Cam = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), Roben.Pos, 0.0, 6.0)
}

func roundToFixed(x float32, places int) float32 {
	shift := math.Pow(10, float64(places))
	return float32(math.Round(float64(x)*shift) / shift)
}

func UpdateCamera() {
	const lerpSpeed float32 = 0.05

	targetX := Roben.Pos.X + playerRectangle.Width/2
	targetY := Roben.Pos.Y + playerRectangle.Height/2

	newX := rl.Lerp(Cam.Target.X, targetX, lerpSpeed)
	newY := rl.Lerp(Cam.Target.Y, targetY, lerpSpeed)

	Cam.Target.X = roundToFixed(newX, 1)
	Cam.Target.Y = roundToFixed(newY, 1)
}

func UpdateTexture() {

}

func DrawPlayer() {
	switch Roben.Dir {
	case entity.RIGHT:
		rl.DrawTextureRec(playerIdle, playerRectangle, Roben.Pos, rl.White)
	case entity.LEFT:
		rl.DrawTextureRec(playerIdle, playerFlippedRectangle, Roben.Pos, rl.White)
	}
}
