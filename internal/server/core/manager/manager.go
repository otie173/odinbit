package manager

import (
	"log"

	"github.com/otie173/odinbit/internal/server/core/ticker"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

type Components struct {
	// Game things
	Textures  *texture.TexturePack
	Overworld *world.World
	Players   player.Storage

	WorldRenderer  *world.Renderer
	PlayerRenderer *player.Renderer

	// Network things
	Handler     *http.Handler
	Listener    *tcp.Listener
	Broadcaster *tcp.Broadcaster

	// System things
	Ticker *ticker.Ticker
}

type Manager struct {
	Components Components
}

func New(components Components) *Manager {
	return &Manager{
		Components: components,
	}
}

func (m *Manager) HandleNetwork(httpAddr, tcpAddr string) {
	go func() {
		if err := m.Components.Handler.Run(httpAddr); err != nil {
			log.Fatalf("Error! Cant run HTTP handler: %v\n", err)
		}
	}()
	log.Printf("HTTP handler running on: %s\n", httpAddr)

	go func() {
		if err := m.Components.Listener.Run(tcpAddr); err != nil {
			log.Fatalf("Error! Cant run TCP listener: %v\n", err)
		}
	}()
	log.Printf("TCP listener running on: %s\n", tcpAddr)
}

func (m *Manager) RenderGame() {
	m.Components.Ticker.Run(func() {
		m.Components.PlayerRenderer.Render()
		m.Components.WorldRenderer.Render()
	})
}
