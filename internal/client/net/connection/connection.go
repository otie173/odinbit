package connection

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Handler struct {
	connection net.Conn
	connected  bool
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) IsConnected() bool {
	return h.connected
}

func (h *Handler) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return err
	}

	h.connection = conn
	h.connected = true
	return nil
}

func (h *Handler) Disconnect() {
	h.connection.Close()
}

func (h *Handler) ConvertPacket(pktType packet.PacketType, pktData interface{}) ([]byte, error) {
	data, err := msgpack.Marshal(&pktData)
	if err != nil {
		return nil, err
	}
	packet := packet.Packet{Type: pktType, Payload: data}
	binaryPacket, err := msgpack.Marshal(&packet)
	if err != nil {
		return nil, err
	}
	return binaryPacket, nil
}

func (h *Handler) Write(data []byte) error {
	if _, err := h.connection.Write(data); err != nil {
		return err
	}
	return nil
}
