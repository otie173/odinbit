package net

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
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
	case packet.CategoryWorld:
		switch pktOpcode {
		case packet.OpcodeWorldUpdate:
			log.Println("Received world area with size in bytes: ", len(pktData))

			var pktStructure packet.WorldUpdate
			var blocks []world.Block

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal binary world area: %v\n", err)
			}

			if err := binary.Unmarshal(pktStructure.Blocks, &blocks); err != nil {
				log.Printf("Error! Cant unmarshal packet structure data to overworld: %v\n", err)
			}
			world.Overworld.Blocks = blocks
			world.Overworld.StartX = pktStructure.StartX
			world.Overworld.StartY = pktStructure.StartY
			world.Overworld.EndX = pktStructure.EndX
			world.Overworld.EndY = pktStructure.EndY

			// log.Println(pktStructure.StartX, pktStructure.StartY, pktStructure.EndX, pktStructure.EndY)
		}
	}
}
