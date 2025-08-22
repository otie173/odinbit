package scene

import (
	"image/color"
	"log"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

var (
	bkgColor = color.RGBA{34, 35, 35, 255}

	nickname     string = "Nickname"
	nicknameEdit bool   = false

	ip     string = "Address"
	ipEdit bool   = false
)

type Handler struct {
	screenWidth  int32
	screenHeight int32
	currentScene common.Scene
	netHandler   *net.Handler
}

func New(screenWidth, screenHeight int32, scene common.Scene, netHandler *net.Handler) *Handler {
	return &Handler{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentScene: scene,
		netHandler:   netHandler,
	}
}

func (h *Handler) drawFunc(fn func()) {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	fn()
	rl.EndDrawing()
}

func (h *Handler) GetScene() common.Scene {
	return h.currentScene
}

func (h *Handler) SetScene(scene common.Scene) {
	h.currentScene = scene
}

func (h *Handler) Handle() {
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
			if raygui.TextBox(rl.NewRectangle(x+40, y+110, 820, 80), &ip, 64, ipEdit) {
				ipEdit = !ipEdit
			}
			if raygui.Button(rl.NewRectangle(float32(h.screenWidth/2-350/2), y+100*3, 350, 85), "Connect") {
				if !h.netHandler.IsConnected() {
					if err := h.netHandler.Connect(ip); err != nil {
						log.Println(err)
						return
					} else {
						log.Printf("Connected to %s\n", ip)
						go h.netHandler.Handle()

						pkt := packet.GetTextures{}
						data, err := h.netHandler.ConvertPacket(packet.GetTexturesType, pkt)
						if err != nil {
							log.Printf("Error convert packet to binary format: %v\n", err)
						}

						if err := h.netHandler.Write(data); err != nil {
							log.Printf("Error with write data to server: %v\n", err)
						}
					}
				}
			}
		})
	}
}
