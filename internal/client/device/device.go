package device

import (
	//"math"

	//"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	//"github.com/otie173/odinbit/internal/client/camera"
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

	scene := h.sceneHandler.GetScene()
	var moveX, moveY float32
	if rl.IsKeyDown(rl.KeyW) && scene == common.Game {
		moveY -= 1
	}
	if rl.IsKeyDown(rl.KeyA) && scene == common.Game {
		moveX -= 1
	}
	if rl.IsKeyDown(rl.KeyS) && scene == common.Game {
		moveY += 1
	}
	if rl.IsKeyDown(rl.KeyD) && scene == common.Game {
		moveX += 1
	}

	if moveX != 0 && moveY != 0 {
		lenght := float32(math.Sqrt(float64(moveX*moveX + moveY*moveY)))
		moveX /= lenght
		moveY /= lenght
	}

	speed := 5.0 * rl.GetFrameTime()
	player.PlayerMu.Lock()
	newX := player.GamePlayer.CurrentX + moveX*speed
	newY := player.GamePlayer.CurrentY + moveY*speed
	player.PlayerMu.Unlock()

	player.PlayerMu.Lock()
	if moveX != 0 {
		if moveX > 0 {
			player.GamePlayer.Flipped = 0
		} else if moveX < 0 {
			player.GamePlayer.Flipped = 1
		}
	}

	// player.GamePlayer.CurrentX = newX
	// player.GamePlayer.CurrentY = newY
	if moveX != 0 || moveY != 0 {
		player.GamePlayer.CurrentX = newX
		player.GamePlayer.CurrentY = newY
		player.PlayerMoved = true
		//player.UpdateServerPos()
	} else {
		player.PlayerMoved = false
	}
	player.PlayerMu.Unlock()

	//log.Println(player.GamePlayer.CurrentX, player.GamePlayer.CurrentY)
	// mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera.Camera)
	// newPlayerX := float32(math.Floor(float64(mousePos.X) / 12))
	// newPlayerY := float32(math.Floor(float64(mousePos.Y) / 12))
	// oldPlayerX := math.Floor(float64(player.GamePlayer.CurrentX))
	// oldPlayerY := math.Floor(float64(player.GamePlayer.CurrentY))
	//log.Println(oldPlayerX, oldPlayerY, oldPlayerX, oldPlayerY)

	// deltaX := math.Abs(float64(newPlayerX) - oldPlayerX)
	// deltaY := math.Abs(float64(newPlayerY) - oldPlayerY)
	//log.Println(deltaX, deltaY)

	// if deltaX <= 6 && deltaY <= 6 {
	// 	if newPlayerX > 0 && newPlayerX < 512 && newPlayerY > 0 && newPlayerY < 512 {
	// 		player.ChangePos(newPlayerX, newPlayerY)
	// 	}
	// }
	//}
}
