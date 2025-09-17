package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/camera"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/scene"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
)

type Client struct {
	title                     string
	screenWidth, screenHeight int32
	deviceHandler             *device.Handler
	sceneHandler              *scene.Handler
	netHandler                *net.Handler
	textureStorage            *texture.Storage
}

func New(title string, screenWidth, screenHeight int32) *Client {
	textureStorage := texture.New()
	netDispatcher := net.NewDispatcher(textureStorage)
	netLoader := net.NewLoader()
	netHandler := net.NewHandler(netDispatcher, netLoader)
	sceneHandler := scene.New(screenWidth, screenHeight, common.Title, netHandler)
	deviceHandler := device.New(sceneHandler)

	return &Client{
		title:          title,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		deviceHandler:  deviceHandler,
		sceneHandler:   sceneHandler,
		netHandler:     netHandler,
		textureStorage: textureStorage,
	}
}

func (c *Client) Load() {
	rl.InitWindow(c.screenWidth, c.screenHeight, c.title)
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowUnfocused | rl.FlagFullscreenMode)
	rl.ToggleFullscreen()
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	rl.SetExitKey(0)
	camera.LoadCamera()
	world.Overworld.Textures = c.textureStorage
	scene.BkgTexture = rl.LoadTexture("resources/backgrounds/background1.png")
	texture.PlayerTexture = rl.LoadTexture("resources/textures/player.png")
}

func (c *Client) update() {
	c.sceneHandler.Handle()
	c.deviceHandler.Handle()
	camera.UpdateCamera()
	player.UpdatePos()

	if c.sceneHandler.GetScene() == common.Game {
		player.UpdateServerPos(c.netHandler)
	}
}

func (c *Client) Run() {
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		c.update()
	}
}
