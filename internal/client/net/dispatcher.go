package net

import (
	"net"

	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

type Dispatcher struct {
	textureStorage *texture.Storage
}

func NewDispatcher(storage *texture.Storage) *Dispatcher {
	return &Dispatcher{
		textureStorage: storage,
	}
}

func (d *Dispatcher) Dispatch(conn *net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, data []byte) {
}
