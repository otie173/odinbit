package net

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/game/world"
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

func (d *Dispatcher) Dispatch(conn *net.Conn, pktType packet.PacketType, data []byte) {
	switch pktType {
	case packet.GetTexturesType:
		texturesPkt := packet.GetTextures{Textures: make(map[string]packet.ServerTexture, 128)}

		if err := msgpack.Unmarshal(data, &texturesPkt.Textures); err != nil {
			log.Printf("Error! Cant unmarshal textures from server: %v\n", err)
		}

		for _, texture := range texturesPkt.Textures {
			d.textureStorage.LoadTexture(texture.Id, texture.Path)
		}
	case packet.GetWorldType:
		var world [common.WorldSize][common.WorldSize]world.Block
		if err := msgpack.Unmarshal(data, &world); err != nil {
			log.Printf("Error! Cant unmarshal world from server: %v\n", err)
		}
		log.Println(world)
	}
}
