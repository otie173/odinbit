package scene

import (
	"image/color"
	"log"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	bkgColor         = color.RGBA{34, 35, 35, 255}
	transparentColor = color.RGBA{34, 34, 35, 200}
	BkgTexture       rl.Texture2D

	selectedMode int = 1

	nickname     string = "Nickname"
	nicknameEdit bool   = false

	httpAddress     string = "HTTP address"
	httpAddressEdit bool   = false

	tcpAddress     string = "TCP address"
	tcpAddressEdit bool   = false
)

const (
	singleplayer int = iota
	multiplayer
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
			rl.DrawTexture(BkgTexture, 0, 0, rl.White)
			rl.DrawRectangle(int32(x), int32(h.screenHeight/2-550/2), 900, 550, transparentColor)
			raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, 32)
			raygui.GroupBox(rl.NewRectangle(x, float32(h.screenHeight/2-550/2), 900, 550), "Odinbit")
			if raygui.Button(rl.NewRectangle(x+40, y, 820, 100), "Singleplayer") {
				selectedMode = singleplayer
			}
			if raygui.Button(rl.NewRectangle(x+40, y+150, 820, 100), "Multiplayer") {
				selectedMode = multiplayer
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
			if raygui.TextBox(rl.NewRectangle(x+40, y+110, 820, 80), &httpAddress, 16, httpAddressEdit) {
				httpAddressEdit = !httpAddressEdit
			}
			if raygui.TextBox(rl.NewRectangle(x+40, y+110*2, 820, 80), &tcpAddress, 64, tcpAddressEdit) {
				tcpAddressEdit = !tcpAddressEdit
			}

			if raygui.Button(rl.NewRectangle(float32(h.screenWidth/2-350/2), y+115*3, 350, 85), "Connect") {
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

					if err := h.netHandler.Connect(tcpAddress); err != nil {
						log.Printf("Error! Cant connect to server: %v\n", err)
						return
					} else {
						log.Printf("Success! Connected to %s\n", tcpAddress)
						h.currentScene = common.Game

						pktStructure := packet.PlayerHandshake{Username: nickname}
						binaryStructure, err := binary.Marshal(&pktStructure)
						if err != nil {
							log.Printf("Error! Cant marshal player handshake structure: %v\n", err)
						}

						pkt := packet.Packet{
							Category: packet.CategoryPlayer,
							Opcode:   packet.OpcodeHandshake,
							Payload:  binaryStructure,
						}

						data, err := binary.Marshal(&pkt)
						if err != nil {
							log.Printf("Error! Cant marshal player handshake packet: %v\n", err)
						}

						compressedPkt, err := net.CompressPkt(data)
						if err != nil {
							log.Printf("Error! Cant compress player handshake packet: %v\n", err)
						}

						h.netHandler.Write(compressedPkt)

						go h.netHandler.Handle()
					}
				}
			}
		})
	case common.Game:
		if selectedMode == multiplayer && !h.netHandler.IsConnected() {
			h.SetScene(common.ConnClosed)
		}

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.BeginMode2D(camera.Camera)
		index := 0
		for x := world.Overworld.StartX; x < world.Overworld.EndX; x++ {
			for y := world.Overworld.StartY; y < world.Overworld.EndY; y++ {
				if index < len(world.Overworld.Blocks) {
					block := world.Overworld.Blocks[index]
					texture := world.GetBlock(block.TextureID)
					rec := rl.NewRectangle(float32(x*12), float32(y*12), 12, 12)
					if block.TextureID != 0 {
						rl.DrawTextureRec(texture, rec, rl.NewVector2(rec.X, rec.Y), rl.White)
					}
					index++
				}
			}
		}
		rl.EndMode2D()
		rl.EndDrawing()
	case common.ConnClosed:
		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.DrawTexture(BkgTexture, 0, 0, rl.White)

		x := float32(h.screenWidth/2 - 900/2)
		groupBoxHeight := float32(350)
		groupBoxY := float32(h.screenHeight/2) - groupBoxHeight/2

		rl.DrawRectangle(int32(x), int32(groupBoxY), 900, int32(groupBoxHeight), transparentColor)
		raygui.GroupBox(rl.NewRectangle(x, groupBoxY, 900, groupBoxHeight), "Notification")

		text := "Connection closed"
		fontSize := int32(32)
		textSize := rl.MeasureTextEx(raygui.GetFont(), text, float32(fontSize), 2)
		textX := x + 450 - float32(textSize.X)/2
		textY := groupBoxY + 120
		raygui.Label(rl.NewRectangle(textX, textY, textSize.X, textSize.Y), text)

		buttonY := groupBoxY + groupBoxHeight - 105
		if raygui.Button(rl.NewRectangle(float32(h.screenWidth/2-300/2), buttonY, 300, 70), "Okay ;(") {
			h.SetScene(common.Title)
		}

		rl.EndDrawing()
	}
}
