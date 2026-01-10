package device

import (
	//"math"

	//"log"
	"log"
	"math"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kelindar/binary"

	//"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/inventory"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/net/compress"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/scene"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

type Handler struct {
	sceneHandler     *scene.Handler
	netModule        *net.Module
	inventoryHandler *inventory.Handler
	textureStorage   *texture.Storage
}

func New(sceneHandler *scene.Handler, netModule *net.Module, inventoryHandler *inventory.Handler, textureStorage *texture.Storage) *Handler {
	return &Handler{
		sceneHandler:     sceneHandler,
		netModule:        netModule,
		inventoryHandler: inventoryHandler,
		textureStorage:   textureStorage,
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

	speed := 5 * rl.GetFrameTime()
	player.PlayerMu.Lock()
	newX := player.GamePlayer.CurrentX + moveX*speed
	newY := player.GamePlayer.CurrentY + moveY*speed

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
		atomic.StoreInt32(&player.PlayerMoved, 1)
		//player.UpdateServerPos()
	} else {
		atomic.StoreInt32(&player.PlayerMoved, 0)
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

	if rl.IsKeyPressed(rl.KeyE) && h.sceneHandler.GetScene() == common.Game {
		switch h.sceneHandler.InventoryOpen {
		case false:
			h.sceneHandler.InventoryOpen = true
		case true:
			h.sceneHandler.InventoryOpen = false
		}

		log.Printf("Состояние инвентаря теперь: %t\n", h.sceneHandler.InventoryOpen)
	}

	if rl.IsKeyPressed(rl.KeyEscape) && h.sceneHandler.GetScene() == common.Game {
		h.netModule.Disconnect()
		h.sceneHandler.SetScene(common.Title)
	}

	// Логика с материалами
	if h.sceneHandler.GetScene() == common.Game {
		if rl.IsKeyPressed(rl.KeyOne) {
			h.inventoryHandler.SetMaterial(common.Wood)
			log.Println("Теперь материал Wood")
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			h.inventoryHandler.SetMaterial(common.Stone)
			log.Println("Теперь материал Stone")
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			h.inventoryHandler.SetMaterial(common.Metal)
			log.Println("Теперь материал Metal")
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) && h.inventoryHandler.GetMaterial() != -1 {
			mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera.Camera)
			mouseX := int(math.Floor(float64(mousePos.X / common.TileSize)))
			mouseY := int(math.Floor(float64(mousePos.Y / common.TileSize)))
			if mouseX > -(common.WorldSize) && mouseX < common.WorldSize &&
				mouseY > -(common.WorldSize) && mouseY < common.WorldSize {
				log.Printf("Блок будет поставлен на %d %d\n", mouseX, mouseY)

				material := h.inventoryHandler.GetMaterial()
				materialID := -1

				switch material {
				case common.Wood:
					materialID = int(h.textureStorage.GetIdByName("wood_material"))
				case common.Stone:
					materialID = int(h.textureStorage.GetIdByName("stone_material"))
				case common.Metal:
					materialID = int(h.textureStorage.GetIdByName("metal_material"))
				}
				log.Println(materialID)

				pktStructure := packet.WorldSetBlock{
					BlockID: materialID,
					X:       mouseX,
					Y:       mouseY,
				}

				binaryStructure, err := binary.Marshal(&pktStructure)
				if err != nil {
					log.Printf("Error! Cant marshal world set block structure: %v\n", err)
				}

				pkt := packet.Packet{
					Category: packet.CategoryWorld,
					Opcode:   packet.OpcodeWorldSetBlock,
					Payload:  binaryStructure,
				}

				data, err := binary.Marshal(&pkt)
				if err != nil {
					log.Printf("Error! Cant marshal world set block packet: %v\n", err)
				}

				compressedPkt, err := compress.CompressPkt(data)
				if err != nil {
					log.Printf("Error! Cant compress binary world set block packet: %v\n", err)
				}

				if err := h.netModule.SendData(compressedPkt); err != nil {
					log.Printf("Error! Cant write world set block packet data to server: %v\n", err)
				}
			}
		}
	}
}
