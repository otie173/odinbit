package client

import (
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/scene"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/server/core/ticker"
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
	netDispatcher := net.NewDispatcher(textureStorage)
	netLoader := net.NewLoader()
	netModule := net.New(netDispatcher, netLoader)
	sceneHandler := scene.New(screenWidth, screenHeight, common.Title, netModule)
	deviceHandler := device.New(sceneHandler)

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
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowUnfocused)
	rl.InitWindow(c.screenWidth, c.screenHeight, c.title)
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	rl.SetExitKey(0)
	camera.LoadCamera()
	world.Overworld.Textures = c.textureStorage
	scene.BkgTexture = rl.LoadTexture("resources/backgrounds/background1.png")
	texture.PlayerTexture = rl.LoadTexture("resources/textures/ghost.png")
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
	c.sceneHandler.Handle()
	c.deviceHandler.Handle()
	camera.UpdateCamera()

	if c.sceneHandler.GetScene() == common.Connect && c.netModule.IsReady() {
		c.sceneHandler.SetScene(common.Game)
	}
}

func (c *Client) Run() {
	defer rl.CloseWindow()

	go c.netModule.Run()
	for !rl.WindowShouldClose() {
		c.update()
	}
}
