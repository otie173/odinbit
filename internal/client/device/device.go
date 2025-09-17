package device

import (
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
		playerX := float32(math.Floor(float64(mousePos.X) / 12))
		playerY := float32(math.Floor(float64(mousePos.Y) / 12))
		player.ChangePos(playerX, playerY)
	}
}
