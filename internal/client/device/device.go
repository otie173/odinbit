package device

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/scene"
)

type Handler struct {
	sceneHandler *scene.Handler
}

func New(sceneHandler *scene.Handler) *Handler {
	return &Handler{
		sceneHandler: sceneHandler,
	}
}

func (h *Handler) Handle() {
	if h.sceneHandler.GetScene() == common.Connect && rl.IsKeyPressed(rl.KeyEscape) {
		h.sceneHandler.SetScene(common.Title)
	}

	if rl.IsKeyDown(rl.KeyW) && h.sceneHandler.GetScene() == common.Game {
		mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera.Camera)
		newPlayerX := float32(math.Floor(float64(mousePos.X) / 12))
		newPlayerY := float32(math.Floor(float64(mousePos.Y) / 12))
		oldPlayerX := math.Floor(float64(player.GamePlayer.CurrentX))
		oldPlayerY := math.Floor(float64(player.GamePlayer.CurrentY))
		log.Println(oldPlayerX, oldPlayerY, oldPlayerX, oldPlayerY)

		deltaX := math.Abs(float64(newPlayerX) - oldPlayerX)
		deltaY := math.Abs(float64(newPlayerY) - oldPlayerY)
		log.Println(deltaX, deltaY)

		if deltaX <= 4 && deltaY <= 4 {
			if newPlayerX > 0 && newPlayerX < 512 && newPlayerY > 0 && newPlayerY < 512 {
				player.ChangePos(newPlayerX, newPlayerY)
			}
		}
	}
}
