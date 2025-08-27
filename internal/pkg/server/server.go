package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
	"github.com/otie173/odinbit/internal/server/ticker"
)

type Server struct {
	addr        string
	textures    *texture.Storage
	world       *world.World
	httpHandler *http.Handler
	tcpHandler  *tcp.Listener
}

func New(addr string) *Server {
	textures := texture.New()
	wrld := world.New(textures)

	textureHandler := texture.NewHandler(textures)
	worldHandler := world.NewHandler(wrld)

	mux := chi.NewRouter()
	router := http.NewRouter(mux)
	httpHandler := http.NewHandler(router, textures, wrld)

	dispatcher := tcp.NewDispatcher(textureHandler, worldHandler)
	tcpHandler := tcp.NewListener(dispatcher)

	return &Server{
		addr:        addr,
		textures:    textures,
		world:       wrld,
		httpHandler: httpHandler,
		tcpHandler:  tcpHandler,
	}
}

func (s *Server) Load() {
	s.textures.LoadTextures()
	s.world.Generate()
}

// FIXME: Add new fields with addr for http & tcp handlers
func (s *Server) Run() {
	httpAddr := "0.0.0.0:9999"
	tcpAddr := "0.0.0.0:8080"

	go func() {
		if err := s.httpHandler.Run(httpAddr); err != nil {
			log.Fatalf("Error! Cant run http handler: %v\n", err)
		}
	}()
	log.Printf("HTTP handler listening on: %s\n", httpAddr)

	go func() {
		if err := s.tcpHandler.Run(tcpAddr); err != nil {
			log.Fatalf("Error! Cant run tcp handler: %v\n", err)
		}
	}()
	log.Printf("TCP handler listening on: %s\n", tcpAddr)

	ticker := ticker.New(20, func() {

	})
	ticker.Run()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down server...")
}
