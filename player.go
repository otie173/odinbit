package main

import (
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
	cam = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), rl.NewVector2(float32(playerPosition.X), float32(playerPosition.Y)), 0.0, 5.0)
}

func unloadPlayer() {
	rl.UnloadTexture(player)
}

func drawPlayer() {
	if playerDirection {
		rl.DrawTextureRec(player, playerRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	} else if !playerDirection {
		rl.DrawTextureRec(player, playerFlippedRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	}
}
