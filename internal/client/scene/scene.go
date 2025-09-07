package scene

import (
	"image/color"
	"log"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
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
	world        *world.World
}

func New(screenWidth, screenHeight int32, scene common.Scene, netHandler *net.Handler, world *world.World) *Handler {
	return &Handler{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentScene: scene,
		netHandler:   netHandler,
		world:        world,
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
					data, err := h.netHandler.LoadTextures("http://0.0.0.0:9999")
					if err != nil {
						log.Printf("Error! Cant load textures from server: %v\n", err)
						return
					}

					pkt := packet.Packet{}
					if err := msgpack.Unmarshal(data, &pkt); err != nil {
						log.Printf("Error! Cant unmarshal body: %v\n", err)
						return
					}
					h.netHandler.Dispatch(nil, pkt.Category, pkt.Opcode, pkt.Payload)

					if err := h.netHandler.Connect(ip); err != nil {
						log.Printf("Error! Cant connect to server: %v\n", err)
						return
					} else {
						log.Printf("Success! Connected to %s\n", ip)
						h.currentScene = common.Game
						go h.netHandler.Handle()
					}
				}
			}
		})
	case common.Game:
		h.drawFunc(func() {

		})
	}
}
