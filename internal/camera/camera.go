package camera

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/common"
	"github.com/otie173/odinbit/internal/player"
)

var (
	Cam rl.Camera2D
)

func Load() {
	Cam = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), player.Player.Pos, 0.0, 8.0)
}

//func Update() {
//	Cam.Target = rl.NewVector2(player.Player.Pos.X+common.TileSize/2, player.Player.Pos.Y+common.TileSize/2)
//}

func roundToFixed(x float32, places int) float32 {
	shift := math.Pow(10, float64(places))
	return float32(math.Round(float64(x)*shift) / shift)
}

func Update() {
	// Локальные константы
	const lerpSpeed float32 = 0.05 // Это значение ближе к исходной версии

	// Вычисляем целевую позицию камеры
	targetX := player.Player.Pos.X + common.TileSize/2
	targetY := player.Player.Pos.Y + common.TileSize/2

	// Используем простую линейную интерполяцию, как в исходной версии
	newX := rl.Lerp(Cam.Target.X, targetX, lerpSpeed)
	newY := rl.Lerp(Cam.Target.Y, targetY, lerpSpeed)

	// Обновляем позицию камеры, округляя до одного знака после запятой
	Cam.Target.X = roundToFixed(newX, 1)
	Cam.Target.Y = roundToFixed(newY, 1)
}
