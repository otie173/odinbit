package handler

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/net/dispatcher"
	"github.com/vmihailenco/msgpack/v5"
)

type TCPHandler struct {
	dispatcher *dispatcher.Dispatcher
}

func NewHandler(dispatcher *dispatcher.Dispatcher) *TCPHandler {
	return &TCPHandler{
		dispatcher: dispatcher,
	}
}

func parsePacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := msgpack.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}
	return pkt, nil
}

func (h *TCPHandler) Handle(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection handling : %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Read buffer error from %s: %v\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}

		pkt, err := parsePacket(buffer[:n])
		if err != nil {
			log.Printf("Parse error from %s: %v\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}
		h.dispatcher.Dispatch(conn, pkt.Type, pkt.Payload)
	}
}
