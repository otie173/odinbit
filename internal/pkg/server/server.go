package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/otie173/odinbit/internal/server/core/manager"
)

type Server struct {
	addr    string
	manager *manager.Manager
}

func New(addr string, manager *manager.Manager) *Server {
	return &Server{
		addr:    addr,
		manager: manager,
	}
}

func (s *Server) Load() {
	s.manager.Components.Textures.LoadTextures()
	s.manager.Components.Overworld.Generate()
}

// FIXME: Add new fields with addr for http & tcp handlers
func (s *Server) Run() {
	httpAddr := "0.0.0.0:9999"
	tcpAddr := "0.0.0.0:8080"
	s.manager.HandleNetwork(httpAddr, tcpAddr)
	s.manager.HandleGame()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Bye! Shutting down server...")
}
