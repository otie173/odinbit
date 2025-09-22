package net

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/player"
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
	case packet.CategoryPlayer:
		switch pktOpcode {
		case packet.OpcodePlayerUpdate:
			pktStructure := packet.PlayerUpdate{}

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal player update data to structure: %v\n", err)
			}

			netPlayers := make([]player.Player, 0, 16)
			if err := binary.Unmarshal(pktStructure.Players, &netPlayers); err != nil {
				log.Printf("Error! Cant unmarshal network players data: %v\n", err)
			}

			player.NetPlayersMu.Lock()
			player.NetworkPlayers = netPlayers
			player.NetPlayersMu.Unlock()

		}
	case packet.CategoryWorld:
		switch pktOpcode {
		case packet.OpcodeWorldUpdate:
			//log.Println("Received world area with size in bytes: ", len(pktData))

			var pktStructure packet.WorldUpdate
			var blocks []world.Block

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal binary world area: %v\n", err)
			}

			if err := binary.Unmarshal(pktStructure.Blocks, &blocks); err != nil {
				log.Printf("Error! Cant unmarshal packet structure data to overworld: %v\n", err)
			}
			world.OverworldMu.Lock()
			world.Overworld.Blocks = blocks
			world.Overworld.StartX = pktStructure.StartX
			world.Overworld.StartY = pktStructure.StartY
			world.Overworld.EndX = pktStructure.EndX
			world.Overworld.EndY = pktStructure.EndY
			world.OverworldMu.Unlock()

			// log.Println(pktStructure.StartX, pktStructure.StartY, pktStructure.EndX, pktStructure.EndY)
		}
	}
}
