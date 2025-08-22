package net

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Handler struct {
	connection net.Conn
	connected  bool
	dispatcher *Dispatcher
}

func NewHandler(dispatcher *Dispatcher) *Handler {
	return &Handler{
		dispatcher: dispatcher,
	}
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

func parsePacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := msgpack.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}
	return pkt, nil
}

func (h *Handler) Handle() {
	log.Println("Началась обработка соединения")
	buffer := make([]byte, 1024)

	for {
		n, err := h.connection.Read(buffer)
		if err != nil {
			log.Printf("Error with read buffer from server: %v\n", err)
		}

		pkt, err := parsePacket(buffer[:n])
		if err != nil {
			log.Printf("Error with parse packet from server: %v\n", err)
		}
		h.dispatcher.Dispatch(h.connection, pkt.Type, pkt.Payload)
	}
}

func (h *Handler) Write(data []byte) error {
	if _, err := h.connection.Write(data); err != nil {
		return err
	}
	return nil
}
