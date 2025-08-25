package tcp

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/texture"
	"github.com/otie173/odinbit/internal/server/world"
	"github.com/vmihailenco/msgpack/v5"
)

type Dispatcher struct {
	textureHandler *texture.Handler
	worldHandler   *world.Handler
}

func NewDispatcher(textureHandler *texture.Handler, worldHandler *world.Handler) *Dispatcher {
	return &Dispatcher{
		textureHandler: textureHandler,
		worldHandler:   worldHandler,
	}
}

func (d *Dispatcher) Dispatch(conn net.Conn, pktType packet.PacketType, data []byte) {
	switch pktType {
	case packet.PingType:
		log.Println("Someone tried ping the server")
	case packet.HandshakeType:
		handshakePkt := packet.Handshake{}
		if err := msgpack.Unmarshal(data, &handshakePkt); err != nil {
			log.Printf("Error! Cant unmarshal packet: %v\n", err)
			return
		}
		log.Printf("Hello, %s\n", handshakePkt.Username)
	case packet.GetTexturesType:
		data, err := d.textureHandler.GetTextures()
		if err != nil {
			log.Printf("Error! Cant get textures: %v\n", err)
		}

		pkt := packet.Packet{
			Type:    packet.GetTexturesType,
			Payload: data,
		}

		binaryPkt, err := msgpack.Marshal(&pkt)
		if err != nil {
			log.Printf("Error! Cant marshal packet: %v\n", err)
		}

		if _, err := conn.Write(binaryPkt); err != nil {
			log.Printf("Error! Cant send binary packet to client: %v\n", err)
		}
	}
}
