package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/device"
	"github.com/otie173/odinbit/internal/client/scene"
)

type Client struct {
	title                     string
	screenWidth, screenHeight int32
	deviceHandler             *device.Handler
	sceneHandler              *scene.Handler
}

func New(title string, screenWidth, screenHeight int32) *Client {
	return &Client{
		title:        title,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (c *Client) Load() {
	rl.InitWindow(c.screenWidth, c.screenHeight, c.title)
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowUnfocused | rl.FlagFullscreenMode)
	rl.ToggleFullscreen()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	c.sceneHandler = scene.New(c.screenWidth, c.screenHeight, common.Title)
	c.deviceHandler = device.New(c.sceneHandler)
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
