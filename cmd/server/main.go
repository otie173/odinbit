package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/otie173/odinbit/internal/pkg/server"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/manager"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

func main() {
	textures := texture.NewStorage()
	wrld := world.New(textures)

	textureHandler := texture.NewHandler(textures)
	worldHandler := world.NewHandler(wrld)

	mux := chi.NewRouter()
	router := http.NewRouter(mux)
	handler := http.NewHandler(router, textures, wrld)

	dispatcher := tcp.NewDispatcher(textureHandler, worldHandler)
	listener := tcp.NewListener(dispatcher)

	managerCfg := manager.Config{
		Textures: textures,
		World:    wrld,
		Handler:  handler,
		Listener: listener,
	}
	manager := manager.New(managerCfg)

	server := server.New("0.0.0.0:8080")
	server.Load()
	server.Run()
}
