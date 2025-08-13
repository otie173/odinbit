package server

import (
	"log"

	"github.com/otie173/odinbit/internal/server/net/handler/http"
	"github.com/otie173/odinbit/internal/server/net/handler/tcp"
)

// TODO: Server { world, players, entities}
// нужно все разделить на компоненты
// чтобы передавать в аргументы функции
type Server struct {
	httpHandler *http.HTTPHandler
	tcpHandler  *tcp.TCPHandler
}

func New() *Server {
	return &Server{
		httpHandler: http.New(),
		tcpHandler:  tcp.New(),
	}
}

func (s *Server) Load() {
	s.httpHandler.Load()
	s.tcpHandler.Load()
}

func (s *Server) Run() {
	log.Println("Hello world!")

	go s.httpHandler.Run()
	s.tcpHandler.Run()
}
