package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/otie173/odinbit/internal/pkg/server"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/manager"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

func main() {
	textures := texture.NewPack()
	overworld := world.New(textures)

	textureHandler := texture.NewHandler(textures)
	worldHandler := world.NewHandler(overworld)

	router := http.NewRouter(chi.NewRouter())
	handler := http.NewHandler(router, textures, overworld)

	dispatcher := tcp.NewDispatcher(textureHandler, worldHandler)
	listener := tcp.NewListener(dispatcher)

	tps := 20
	tickDuration := time.Second / time.Duration(tps)
	ticker := time.NewTicker(tickDuration)

	components := manager.Components{
		Textures:  textures,
		Overworld: overworld,
		Handler:   handler,
		Listener:  listener,
		Ticker:    ticker,
	}
	manager := manager.New(components)

	server := server.New("0.0.0.0:8080", manager)
	server.Load()
	server.Run()
}
