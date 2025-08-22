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

func (d *Dispatcher) Dispatch(conn net.Conn, pktType packet.PacketType, data []byte) {
	switch pktType {
	case packet.GetTexturesType:
		log.Println("Receive textures form server!")
		texturesPkt := packet.GetTextures{}

		if err := msgpack.Unmarshal(data, &texturesPkt.Textures); err != nil {
			log.Printf("Error! Can't unmarshal textures from server: %v\n", err)
		}

		for _, texture := range texturesPkt.Textures {
			d.textureStorage.LoadTexture(texture.Id, texture.Path)
		}
	}
}
