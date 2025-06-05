package device

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/camera"
	"github.com/otie173/odinbit/internal/common"
	"github.com/otie173/odinbit/internal/player"
	"github.com/otie173/odinbit/internal/world"
)

func HandleKeyboard() {
	if rl.IsKeyPressed(rl.KeyW) {
		mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera.Cam)
		tileX := int(math.Floor(float64(mousePos.X / common.TileSize)))
		tileY := int(math.Floor(float64(mousePos.Y / common.TileSize)))

		if world.BlockExists(tileX, tileY) {
			world.CheckBehavior(tileX, tileY)
		}

		if world.IsValidPos(tileX, tileY) && !world.BlockExists(tileX, tileY) || world.IsValidPos(tileX, tileY) && world.IsPassable(tileX, tileY) {
			player.SetTilePos(tileX, tileY)
		}
	}
}
