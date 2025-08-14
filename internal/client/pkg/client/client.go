package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/client/scene"
)

type Client struct {
	title                     string
	screenWidth, screenHeight int32
	sceneHandler              *scene.Handler
}

func New(title string, screenWidth, screenHeight int32) *Client {
	return &Client{
		title:        title,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		sceneHandler: scene.New(screenWidth, screenHeight, common.Title),
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
	c.sceneHandler.Update()
}

func (c *Client) Run() {
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		c.update()
	}
}
