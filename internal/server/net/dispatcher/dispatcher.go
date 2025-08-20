package dispatcher

import (
	"log"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Dispatcher struct {
}

func New() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(pktType packet.PacketType, data []byte) {
	switch pktType {
	case packet.HandshakeType:
		handshakePkt := &packet.Handshake{}
		if err := msgpack.Unmarshal(data, &handshakePkt); err != nil {
			log.Printf("Error with unmarshal packet: %v\n", err)
			return
		}
		log.Printf("Hello, %s\n", handshakePkt.Username)
	}
}
