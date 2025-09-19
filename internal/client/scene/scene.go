package scene

import (
	"image/color"
	"log"
	"math"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/player"
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
	uiScale      float32
}

func New(screenWidth, screenHeight int32, scene common.Scene, netHandler *net.Handler) *Handler {
	scaleX := float32(screenWidth) / float32(common.BaseRenderWidth)
	scaleY := float32(screenHeight) / float32(common.BaseRenderHeight)
	uiScale := float32(math.Min(float64(scaleX), float64(scaleY)))
	if uiScale <= 0 {
		uiScale = 1
	}
	return &Handler{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentScene: scene,
		netHandler:   netHandler,
		uiScale:      uiScale,
	}
}

func roundToInt32(v float32) int32 {
	return int32(math.Round(float64(v)))
}

func sizeToInt32(v float32) int32 {
	val := int32(math.Round(float64(v)))
	if val < 1 {
		return 1
	}
	return val
}

func (h *Handler) scale(v float32) float32 {
	return v * h.uiScale
}

func (h *Handler) scaledInt(v int32) int32 {
	scaled := int32(math.Round(float64(float32(v) * h.uiScale)))
	if scaled < 1 {
		return 1
	}
	return scaled
}

func (h *Handler) drawBackground() {
	if BkgTexture.ID == 0 {
		return
	}
	source := rl.NewRectangle(0, 0, float32(BkgTexture.Width), float32(BkgTexture.Height))
	dest := rl.NewRectangle(0, 0, float32(h.screenWidth), float32(h.screenHeight))
	rl.DrawTexturePro(BkgTexture, source, dest, rl.NewVector2(0, 0), 0, rl.White)
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
			h.drawBackground()
			panelWidth := h.scale(900)
			panelHeight := h.scale(550)
			panelX := float32(h.screenWidth)/2 - panelWidth/2
			panelY := float32(h.screenHeight)/2 - panelHeight/2
			buttonX := panelX + h.scale(40)
			buttonY := panelY + h.scale(75)
			buttonWidth := h.scale(820)
			buttonHeight := h.scale(100)
			buttonSpacing := h.scale(150)

			rl.DrawRectangle(roundToInt32(panelX), roundToInt32(panelY), sizeToInt32(panelWidth), sizeToInt32(panelHeight), transparentColor)
			raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, int64(h.scaledInt(32)))
			raygui.GroupBox(rl.NewRectangle(panelX, panelY, panelWidth, panelHeight), "Odinbit")
			if raygui.Button(rl.NewRectangle(buttonX, buttonY, buttonWidth, buttonHeight), "Singleplayer") {
				selectedMode = singleplayer
			}
			if raygui.Button(rl.NewRectangle(buttonX, buttonY+buttonSpacing, buttonWidth, buttonHeight), "Multiplayer") {
				selectedMode = multiplayer
				h.currentScene = common.Connect
			}
			if raygui.Button(rl.NewRectangle(buttonX, buttonY+buttonSpacing*2, buttonWidth, buttonHeight), "Exit") {
				rl.CloseWindow()
				os.Exit(0)
			}
		})
	case common.Connect:
		h.drawFunc(func() {
			h.drawBackground()
			panelWidth := h.scale(900)
			panelHeight := h.scale(550)
			panelX := float32(h.screenWidth)/2 - panelWidth/2
			panelY := float32(h.screenHeight)/2 - panelHeight/2
			fieldX := panelX + h.scale(40)
			fieldY := panelY + h.scale(75)
			fieldWidth := h.scale(820)
			fieldHeight := h.scale(80)
			fieldSpacing := h.scale(110)

			rl.DrawRectangle(roundToInt32(panelX), roundToInt32(panelY), sizeToInt32(panelWidth), sizeToInt32(panelHeight), transparentColor)
			raygui.GroupBox(rl.NewRectangle(panelX, panelY, panelWidth, panelHeight), "Connect")
			if raygui.TextBox(rl.NewRectangle(fieldX, fieldY, fieldWidth, fieldHeight), &nickname, 64, nicknameEdit) {
				nicknameEdit = !nicknameEdit
			}
			if raygui.TextBox(rl.NewRectangle(fieldX, fieldY+fieldSpacing, fieldWidth, fieldHeight), &httpAddress, 16, httpAddressEdit) {
				httpAddressEdit = !httpAddressEdit
			}
			if raygui.TextBox(rl.NewRectangle(fieldX, fieldY+fieldSpacing*2, fieldWidth, fieldHeight), &tcpAddress, 64, tcpAddressEdit) {
				tcpAddressEdit = !tcpAddressEdit
			}

			buttonWidth := h.scale(350)
			buttonHeight := h.scale(85)
			buttonX := float32(h.screenWidth)/2 - buttonWidth/2
			buttonY := panelY + h.scale(420)
			if raygui.Button(rl.NewRectangle(buttonX, buttonY, buttonWidth, buttonHeight), "Connect") {
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
						player.NetConnection = h.netHandler

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
		player.DrawPlayer()
		rl.EndMode2D()
		rl.EndDrawing()
	case common.ConnClosed:
		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		h.drawBackground()

		panelWidth := h.scale(900)
		panelHeight := h.scale(350)
		panelX := float32(h.screenWidth)/2 - panelWidth/2
		panelY := float32(h.screenHeight)/2 - panelHeight/2

		rl.DrawRectangle(roundToInt32(panelX), roundToInt32(panelY), sizeToInt32(panelWidth), sizeToInt32(panelHeight), transparentColor)
		raygui.GroupBox(rl.NewRectangle(panelX, panelY, panelWidth, panelHeight), "Notification")

		text := "Connection closed"
		fontSize := float32(h.scaledInt(32))
		spacing := h.scale(2)
		textSize := rl.MeasureTextEx(raygui.GetFont(), text, fontSize, spacing)
		textX := panelX + panelWidth/2 - textSize.X/2
		textY := panelY + h.scale(120)
		raygui.Label(rl.NewRectangle(textX, textY, textSize.X, textSize.Y), text)

		buttonWidth := h.scale(300)
		buttonHeight := h.scale(70)
		buttonX := float32(h.screenWidth)/2 - buttonWidth/2
		buttonY := panelY + panelHeight - h.scale(105)
		if raygui.Button(rl.NewRectangle(buttonX, buttonY, buttonWidth, buttonHeight), "Okay ;(") {
			h.SetScene(common.Title)
		}

		rl.EndDrawing()
	}
}
