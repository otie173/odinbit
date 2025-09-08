package tcp

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
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

func (d *Dispatcher) Dispatch(conn net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, pktData []byte) {
	switch pktCategory {
	case packet.CategoryWorld:
	case packet.CategoryPlayer:
		switch pktOpcode {
		case packet.OpcodeHandshake:
			pktStructure := packet.PlayerHandshake{}

			if err := msgpack.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal player handshake data: %v\n", err)
			}
			log.Printf("Hi, %s!\n", pktStructure.Username)
		}
	case packet.CategoryInventory:
	}
}
