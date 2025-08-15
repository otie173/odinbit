package handler

import (
	"log"
	"net"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	log.Printf("Handling connection: %s\n", conn.RemoteAddr().String())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Problem with read data from %s with error %v\n", conn.RemoteAddr().String(), err)
		}

		data := string(buffer[:n])
		log.Printf("Received data from %s: %s\n", conn.RemoteAddr().String(), data)
	}
}
