package client

import (
	//"log"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/inventory"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/scene"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/server/core/ticker"
)

var (
	MainChan  = make(chan texture.Texture)
	ReadyChan = make(chan bool)
)

type Client struct {
	title                     string
	screenWidth, screenHeight int32
	deviceHandler             *device.Handler
	sceneHandler              *scene.Handler
	netModule                 *net.Module
	textureStorage            *texture.Storage
	ticker                    *ticker.Ticker
}

func New(title string, screenWidth, screenHeight int32) *Client {
	textureStorage := texture.New()
	netDispatcher := net.NewDispatcher(MainChan, ReadyChan, textureStorage)
	netLoader := net.NewLoader()
	netModule := net.New(netDispatcher, netLoader)
	inventoryHandler := inventory.NewHandler(inventory.NewInventory())
	sceneHandler := scene.New(screenWidth, screenHeight, common.Title, netModule, inventoryHandler, textureStorage)
	deviceHandler := device.New(sceneHandler, netModule, inventoryHandler, textureStorage)

	return &Client{
		title:          title,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		deviceHandler:  deviceHandler,
		sceneHandler:   sceneHandler,
		netModule:      netModule,
		textureStorage: textureStorage,
	}
}

func (c *Client) Load() {
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowUnfocused | rl.FlagFullscreenMode)
	rl.InitWindow(c.screenWidth, c.screenHeight, c.title)
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	rl.SetExitKey(0)
	camera.LoadCamera()
	world.Overworld.Textures = c.textureStorage
	scene.BkgTexture = rl.LoadTexture("resources/backgrounds/background1.png")
	texture.PlayerTexture = rl.LoadTexture("resources/textures/player.png")
	texture.WoodMaterial = rl.LoadTexture("resources/textures/wood_material.png")
	texture.StoneMaterial = rl.LoadTexture("resources/textures/stone_material.png")
	texture.MetalMaterial = rl.LoadTexture("resources/textures/metal_material.png")
	scene.SlotTexture = rl.LoadTexture("resources/ui/icons/slot.png")
	c.ticker = ticker.New(10)

	go func() {
		c.ticker.Run(func() {
			if atomic.LoadInt32(&player.PlayerMoved) == 1 {
				c.netModule.UpdateServerPos()
			}
		})
	}()
}

func (c *Client) update() {
	select {
	case ready, ok := <-ReadyChan:
		if !ok {
			ReadyChan = nil
			break
		}
		c.netModule.SetReady(ready)
	case texture, ok := <-MainChan:
		if !ok {
			MainChan = nil
			break
		}
		c.textureStorage.LoadTexture(texture.Id, texture.Path)
		//log.Println(texture)
	default:
		goto UPDATE
	}

UPDATE:
	c.sceneHandler.Handle()
	c.deviceHandler.Handle()
	camera.UpdateCamera()

	if c.sceneHandler.GetScene() == common.Connecting && c.netModule.IsReady() {
		c.sceneHandler.SetScene(common.Game)
	}
}

func (c *Client) Run() {
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		c.update()
	}
}
