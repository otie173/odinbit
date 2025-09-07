package net

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Dispatcher struct {
	textureStorage *texture.Storage
}

func NewDispatcher(storage *texture.Storage) *Dispatcher {
	return &Dispatcher{
		textureStorage: storage,
	}
}

func (d *Dispatcher) Dispatch(conn *net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, pktData []byte) {
	switch pktCategory {
	case packet.CategoryTexture:
		switch pktOpcode {
		case packet.OpcodeTextureData:
			var pktStructure packet.TextureData

			if err := msgpack.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal texture data: %v\n", err)
			}

			for _, texture := range pktStructure.Textures {
				d.textureStorage.LoadTexture(texture.Id, texture.Path)
			}
		}
	}
}
