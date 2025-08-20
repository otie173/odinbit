package server

import (
	"log"
	"net"
	"os"

	"github.com/otie173/odinbit/internal/server/net/dispatcher"
	"github.com/otie173/odinbit/internal/server/net/handler"
	"github.com/otie173/odinbit/internal/server/texture"
	"github.com/otie173/odinbit/internal/server/world"
)

type Server struct {
	addr        string
	textures    *texture.Storage
	world       *world.World
	mainHandler *handler.Handler
}

func New(addr string) *Server {
	textures := texture.New()
	wrld := world.New(textures)

	textureHandler := texture.NewHandler(textures)
	worldHandler := world.NewHandler(wrld)

	dispatcher := dispatcher.New(textureHandler, worldHandler)
	mainHandler := handler.New(dispatcher)

	return &Server{
		addr:        addr,
		textures:    textures,
		world:       wrld,
		mainHandler: mainHandler,
	}
}

func (s *Server) Load() {
	s.textures.LoadTextures()
	s.world.Generate()
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Printf("Server listening on: %s\n", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error with accept connection: %v\n", err)
			break
		}
		log.Printf("New connection accepted: %s\n", conn.RemoteAddr().String())
		go s.mainHandler.Handle(conn)
	}
}
