package device

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
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
}
