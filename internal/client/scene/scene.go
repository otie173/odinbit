package scene

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/inventory"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	bkgColor         = color.RGBA{34, 35, 35, 255}
	transparentColor = color.RGBA{34, 34, 35, 200}
	BkgTexture       rl.Texture2D
	SlotTexture      rl.Texture2D
	CurrentPage      int = 1
	MaxPage          int = 2

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
	screenWidth      int32
	screenHeight     int32
	currentScene     common.Scene
	InventoryOpen    bool
	netModule        *net.Module
	uiScale          float32
	inventoryHandler *inventory.Handler
}

func New(screenWidth, screenHeight int32, scene common.Scene, netModule *net.Module, inventoryHandler *inventory.Handler) *Handler {
	scaleX := float32(screenWidth) / float32(common.BaseRenderWidth)
	scaleY := float32(screenHeight) / float32(common.BaseRenderHeight)
	uiScale := float32(math.Min(float64(scaleX), float64(scaleY)))
	if uiScale <= 0 {
		uiScale = 1
	}
	return &Handler{
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		currentScene:     scene,
		netModule:        netModule,
		uiScale:          uiScale,
		inventoryHandler: inventoryHandler,
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
				if !h.netModule.IsConnected() {
					data, err := h.netModule.LoadTextures("http://0.0.0.0:9999")
					if err != nil {
						log.Printf("Error! Cant load textures from server: %v\n", err)
						return
					}

					pkt := packet.Packet{}
					if err := msgpack.Unmarshal(data, &pkt); err != nil {
						log.Printf("Error! Cant unmarshal body: %v\n", err)
						return
					}
					h.netModule.Dispatch(nil, pkt.Category, pkt.Opcode, pkt.Payload)

					if err := h.netModule.Connect(tcpAddress); err != nil {
						log.Printf("Error! Cant connect to server: %v\n", err)
						return
					} else {
						log.Printf("Success! Connected to %s\n", tcpAddress)

						pktStructure := packet.PlayerHandshake{Username: nickname}
						binaryStructure, err := binary.Marshal(&pktStructure)
						if err != nil {
							log.Printf("Error! Cant marshal player handshake structure: %v\n", err)
						}

						pkt := packet.Packet{
							Category: packet.CategoryPlayer,
							Opcode:   packet.OpcodePlayerHandshake,
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

						h.netModule.SendData(compressedPkt)
					}
				}
			}
		})
	case common.Game:
		if selectedMode == multiplayer && !h.netModule.IsConnected() {
			h.SetScene(common.ConnClosed)
		}

		rl.BeginDrawing()
		rl.ClearBackground(bkgColor)
		rl.BeginMode2D(camera.Camera)
		world.OverworldMu.Lock()
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
		world.OverworldMu.Unlock()
		player.DrawPlayer()

		if h.netModule.IsConnected() && h.netModule.IsReady() {
			player.UpdateNetworkPlayers()
			player.DrawNetworkPlayers()
		}
		rl.EndMode2D()

		// рисовка инвентаря с ячейками
		if h.InventoryOpen {
			panelWidth := h.scale(550)
			panelHeight := h.scale(500)
			panelX := float32(h.screenWidth)/2 - panelWidth/2
			panelY := float32(h.screenHeight)/2 - panelHeight/2

			rl.DrawRectangle(roundToInt32(panelX), roundToInt32(panelY), sizeToInt32(panelWidth), sizeToInt32(panelHeight), transparentColor)
			raygui.GroupBox(rl.NewRectangle(panelX, panelY, panelWidth, panelHeight), "Blocks")

			slotScale := float32(7.0)
			slot1Pos := rl.NewVector2(float32(h.screenWidth)/2-210, float32(h.screenHeight)/2-170)
			slot2Pos := rl.NewVector2(slot1Pos.X+155, slot1Pos.Y)
			slot3Pos := rl.NewVector2(slot2Pos.X+155, slot1Pos.Y)
			slot4Pos := rl.NewVector2(slot1Pos.X, slot1Pos.Y+155)
			slot5Pos := rl.NewVector2(slot2Pos.X, slot4Pos.Y)
			slot6Pos := rl.NewVector2(slot3Pos.X, slot4Pos.Y)

			rl.DrawTextureEx(SlotTexture, slot1Pos, 0, slotScale, rl.White)
			rl.DrawTextureEx(SlotTexture, slot2Pos, 0, slotScale, rl.White)
			rl.DrawTextureEx(SlotTexture, slot3Pos, 0, slotScale, rl.White)
			rl.DrawTextureEx(SlotTexture, slot4Pos, 0, slotScale, rl.White)
			rl.DrawTextureEx(SlotTexture, slot5Pos, 0, slotScale, rl.White)
			rl.DrawTextureEx(SlotTexture, slot6Pos, 0, slotScale, rl.White)

			pageString := fmt.Sprintf("Page: %d/%d", CurrentPage, MaxPage)
			pageStringSize := rl.MeasureTextEx(raygui.GetFont(), pageString, 30, 2)
			pageStringPos := rl.NewVector2((float32(h.screenWidth)-float32(pageStringSize.X))/2, 700)
			pageStringRec := rl.NewRectangle(pageStringPos.X, pageStringPos.Y, pageStringSize.X, pageStringSize.Y)
			raygui.Label(pageStringRec, pageString)
		}

		recWidth := float32(0)
		playerX := player.GamePlayer.CurrentX
		playerY := player.GamePlayer.CurrentY

		if playerX >= 1000 && playerY >= 1000 {
			recWidth = 295
		} else {
			recWidth = 275
		}

		// ободок там где предметы
		rl.DrawRectangleV(rl.NewVector2(0, 0), rl.NewVector2(recWidth, 225), transparentColor)
		rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, recWidth, 225), 5, rl.White)

		textureScale := float32(5)
		texturePos := rl.NewVector2(15, 15)
		rl.DrawTextureEx(texture.WoodMaterial, texturePos, 0, textureScale, rl.White)
		materialCount := h.inventoryHandler.GetMaterialCount(common.Wood)
		countText := fmt.Sprintf("%d", materialCount)
		countPos := rl.NewVector2(90, 30)
		textFont := raygui.GetFont()
		rl.DrawTextEx(textFont, countText, countPos, 25, 2, rl.White)

		texturePos = rl.NewVector2(15, 75)
		rl.DrawTextureEx(texture.StoneMaterial, texturePos, 0, textureScale, rl.White)
		materialCount = h.inventoryHandler.GetMaterialCount(common.Stone)
		countText = fmt.Sprintf("%d", materialCount)
		countPos = rl.NewVector2(90, 90)
		rl.DrawTextEx(textFont, countText, countPos, 25, 2, rl.White)

		texturePos = rl.NewVector2(15, 120)
		rl.DrawTextureEx(texture.MetalMaterial, texturePos, 0, textureScale, rl.White)
		materialCount = h.inventoryHandler.GetMaterialCount(common.Metal)
		countText = fmt.Sprintf("%d", materialCount)
		countPos = rl.NewVector2(90, 140)
		rl.DrawTextEx(textFont, countText, countPos, 25, 2, rl.White)

		textContent := fmt.Sprintf("X: %.0f Y: %.0f", playerX, playerY)
		textPos := rl.NewVector2(25, 185)
		rl.DrawTextEx(textFont, textContent, textPos, 24, 2, rl.White)

		/*
			textFont := raygui.GetFont()
			// надпись материала
			textContent := fmt.Sprintf("Material: %s", common.MaterialMap[h.inventoryHandler.GetMaterial()])
			textPos := rl.NewVector2(10, 5)
			rl.DrawTextEx(textFont, textContent, textPos, 32, 2, rl.White)

			// надпись дерева
			textContent = fmt.Sprintf("Wood: ")
			textPos = rl.NewVector2(10, 45)
			rl.DrawTextEx(textFont, textContent, textPos, 32, 2, rl.White)

			// надпись камня
			textContent = fmt.Sprintf("Stone: ")
			textPos = rl.NewVector2(10, 85)
			rl.DrawTextEx(textFont, textContent, textPos, 32, 2, rl.White)

			// надпись металла
			textContent = fmt.Sprintf("Metal: ")
			textPos = rl.NewVector2(10, 125)
			rl.DrawTextEx(textFont, textContent, textPos, 32, 2, rl.White)
		*/

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
