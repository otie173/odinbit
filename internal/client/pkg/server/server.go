package server

import (
	"log"
	"net"
	"os"

	"github.com/otie173/odinbit/internal/server/net/handler"
)

// TODO: Server { world, players, entities}
// нужно все разделить на компоненты
// чтобы передавать в аргументы функции
type Server struct {
	addr    string
	handler *handler.Handler
}

func New(addr string) *Server {
	return &Server{
		addr:    addr,
		handler: handler.New(),
	}
}

func (s *Server) Load() {
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
		log.Printf("New connection: %s\n", conn.RemoteAddr().String())

		go s.handler.Handle(conn)
	}
}
