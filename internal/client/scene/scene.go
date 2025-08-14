package scene

import (
	"image/color"
	"log"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/net/connection"
)

var (
	bkgColor = color.RGBA{34, 35, 35, 255}

	nickname     string = "Nickname"
	nicknameEdit bool   = false

	password     string = "Password"
	passwordEdit bool   = false

	ip     string = "Address"
	ipEdit bool   = false
)

type Handler struct {
	screenWidth  int32
	screenHeight int32
	currentScene common.Scene
}

func New(screenWidth, screenHeight int32, scene common.Scene) *Handler {
	return &Handler{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentScene: scene,
	}
}

func (h *Handler) drawFunc(fn func()) {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	fn()
	rl.EndDrawing()
}

func (h *Handler) Update() {
	switch h.currentScene {
	case common.Title:
		h.drawFunc(func() {
			x := float32(h.screenWidth/2 - 900/2)
			y := float32(340)
			raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, 32)
			raygui.GroupBox(rl.NewRectangle(x, float32(h.screenHeight/2-550/2), 900, 550), "Odinbit")

			if raygui.Button(rl.NewRectangle(x+40, y, 820, 100), "Singleplayer") {
			}
			if raygui.Button(rl.NewRectangle(x+40, y+150, 820, 100), "Multiplayer") {
				h.currentScene = common.Connect
			}
			if raygui.Button(rl.NewRectangle(x+40, y+150*2, 820, 100), "Exit") {
				rl.CloseWindow()
				os.Exit(0)
			}
		})
	case common.Connect:
		h.drawFunc(func() {
			x := float32(h.screenWidth/2 - 900/2)
			y := float32(340)
			raygui.GroupBox(rl.NewRectangle(x, float32(h.screenHeight/2-550/2), 900, 550), "Connect")
			if raygui.TextBox(rl.NewRectangle(x+40, y, 820, 80), &nickname, 64, nicknameEdit) {
				nicknameEdit = !nicknameEdit
			}
			if raygui.TextBox(rl.NewRectangle(x+40, y+110, 820, 80), &password, 64, passwordEdit) {
				passwordEdit = !passwordEdit
			}
			if raygui.TextBox(rl.NewRectangle(x+40, y+110*2, 820, 80), &ip, 64, ipEdit) {
				ipEdit = !ipEdit
			}
			if raygui.Button(rl.NewRectangle(float32(h.screenWidth/2-350/2), y+115*3, 350, 85), "Connect") {
				if !connection.IsConnected() {
					if err := connection.Connect(ip); err != nil {
						log.Println(err)
						return
					}

					if connection.IsConnected() {
						connection.Write([]byte("Hello!"))
					}
				}
			}
		})
	}
}
