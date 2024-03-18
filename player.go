package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player                 rl.Texture2D
	playerPosition         rl.Vector2   = rl.NewVector2(0.0, 0.0)
	playerRectangle        rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, TILE_SIZE, TILE_SIZE)
	playerFlippedRectangle rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, -TILE_SIZE, TILE_SIZE)
	playerDirection        bool         = false
	cam                    rl.Camera2D
	canMove                bool
	canBuild               bool
)

func loadPlayer() {
	player = rl.LoadTexture("assets/images/characters/player.png")
	cam = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), rl.NewVector2(float32(playerPosition.X), float32(playerPosition.Y)), 0.0, 6.0)
}

func unloadPlayer() {
	rl.UnloadTexture(player)
}

func roundToFixed(x float32, places int) float32 {
	shift := math.Pow(10, float64(places))
	return float32(math.Round(float64(x)*shift) / shift)
}

func updateCameraTarget(cam *rl.Camera2D, playerPosition rl.Vector2, playerRectangle rl.Rectangle) {
	targetX := playerPosition.X + playerRectangle.Width/2
	targetY := playerPosition.Y + playerRectangle.Height/2

	newX := rl.Vector2Lerp(cam.Target, rl.NewVector2(targetX, cam.Target.Y), 0.025).X
	newY := rl.Vector2Lerp(cam.Target, rl.NewVector2(cam.Target.X, targetY), 0.025).Y

	cam.Target.X = roundToFixed(newX, 1)
	cam.Target.Y = roundToFixed(newY, 1)
}

func drawPlayer() {
	if playerDirection {
		rl.DrawTextureRec(player, playerRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	} else if !playerDirection {
		rl.DrawTextureRec(player, playerFlippedRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	}
}
