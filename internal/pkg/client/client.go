package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/net"
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
	world                     *world.World
}

func New(title string, screenWidth, screenHeight int32) *Client {
	textureStorage := texture.New()
	netDispatcher := net.NewDispatcher(textureStorage)
	netLoader := net.NewLoader()
	netHandler := net.NewHandler(netDispatcher, netLoader)
	world := world.New(textureStorage)
	sceneHandler := scene.New(screenWidth, screenHeight, common.Title, netHandler, world)
	deviceHandler := device.New(sceneHandler)

	return &Client{
		title:          title,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		deviceHandler:  deviceHandler,
		sceneHandler:   sceneHandler,
		netHandler:     netHandler,
		textureStorage: textureStorage,
		world:          world,
	}
}

func (c *Client) Load() {
	rl.InitWindow(c.screenWidth, c.screenHeight, c.title)
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowUnfocused | rl.FlagFullscreenMode)
	rl.ToggleFullscreen()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)
}

func (c *Client) update() {
	c.sceneHandler.Handle()
	c.deviceHandler.Handle()
}

func (c *Client) Run() {
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		c.update()
	}
}
