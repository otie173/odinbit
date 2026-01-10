package net

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/client/world"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

type Dispatcher struct {
	mainChan       chan texture.Texture
	readyChan      chan bool
	textureStorage *texture.Storage
}

func NewDispatcher(mainChan chan texture.Texture, readyChan chan bool, storage *texture.Storage) *Dispatcher {
	return &Dispatcher{
		mainChan:       mainChan,
		readyChan:      readyChan,
		textureStorage: storage,
	}
}

func (d *Dispatcher) Dispatch(conn *net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, pktData []byte) {
	switch pktCategory {
	case packet.CategoryConnection:
		switch pktOpcode {
		case packet.OpcodeConnectResponse:
			log.Printf("Какая-то инфа с сервака пришла: %d bytes\n", len(pktData))

			pktStructure := packet.ConnectResponse{}
			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal connect response packet: %v\n", err)
			}

			var textures packet.TextureData
			if err := binary.Unmarshal(pktStructure.TexturesData, &textures); err != nil {
				log.Printf("Error! Cant unmarshal texture data: %v\n", err)
			}
			//log.Println(textures)

			for _, data := range textures.Textures {
				loadedTexture := texture.Texture{
					Id:   data.Id,
					Path: data.Path,
				}
				d.mainChan <- loadedTexture
			}
			d.readyChan <- true
			close(d.mainChan)
			close(d.readyChan)

			//log.Println(pktStructure.BlocksData)

		}
	case packet.CategoryTexture:
		switch pktOpcode {
		case packet.OpcodeTextureData:
			var pktStructure packet.TextureData

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
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
			player.NetworkPlayersRaw = netPlayers
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
