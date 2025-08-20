package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/net/connection"
	"github.com/otie173/odinbit/internal/client/scene"
	"github.com/otie173/odinbit/internal/client/texture"
)

type Client struct {
	title                     string
	screenWidth, screenHeight int32
	deviceHandler             *device.Handler
	sceneHandler              *scene.Handler
	connHandler               *connection.Handler
	textureStorage            *texture.Storage
}

func New(title string, screenWidth, screenHeight int32) *Client {
	connHandler := connection.New()
	sceneHandler := scene.New(screenWidth, screenHeight, common.Title, connHandler)
	deviceHandler := device.New(sceneHandler)
	textureStorage := texture.New()

	return &Client{
		title:          title,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		deviceHandler:  deviceHandler,
		sceneHandler:   sceneHandler,
		connHandler:    connHandler,
		textureStorage: textureStorage,
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
