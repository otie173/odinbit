package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/otie173/odinbit/internal/pkg/server"
	"github.com/otie173/odinbit/internal/server/core/manager"
	"github.com/otie173/odinbit/internal/server/core/ticker"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

func main() {
	textures := texture.NewPack()
	overworld := world.New(textures)
	playerStorage := player.NewStorage(16)

	textureHandler := texture.NewHandler(textures)
	worldHandler := world.NewHandler(overworld, playerStorage)

	router := http.NewRouter(chi.NewRouter())
	handler := http.NewHandler(router, textures, overworld)

	dispatcher := tcp.NewDispatcher(playerStorage, textureHandler, worldHandler)
	listener := tcp.NewListener(dispatcher)

	tps := 20
	ticker := ticker.New(tps)

	components := manager.Components{
		Textures:     textures,
		Overworld:    overworld,
		Players:      playerStorage,
		WorldHandler: worldHandler,
		Handler:      handler,
		Listener:     listener,
		Ticker:       ticker,
	}
	manager := manager.New(components)

	server := server.New("0.0.0.0:8080", manager)
	server.Load()
	server.Run()
}
