package handler

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/net/dispatcher"
	"github.com/vmihailenco/msgpack/v5"
)

type Handler struct {
	dispatcher *dispatcher.Dispatcher
}

func New(dispatcher *dispatcher.Dispatcher) *Handler {
	return &Handler{
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

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection handling : %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Read buffer error from %s: %s\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}
		data = buffer[:n]

		pkt, err := parsePacket(data)
		if err != nil {
			log.Printf("Parse error from %s: %s\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}
		h.dispatcher.Dispatch(conn, pkt.Type, pkt.Payload)
	}
}
