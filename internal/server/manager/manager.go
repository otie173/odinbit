package manager

import (
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

type Config struct {
	Textures *texture.Storage
	World    *world.World
	Handler  *http.Handler
	Listener *tcp.Listener
}

type Manager struct {
	Cfg Config
}

func New(cfg Config) *Manager {
	return &Manager{Cfg: cfg}
}
