package tcp

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/game/blocks"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/compress"
)

type Dispatcher struct {
	playerStorage  player.Storage
	textureHandler *texture.Handler
	blocksStorage  *blocks.Storage
	world          *world.World
}

func NewDispatcher(playerStorage player.Storage, textureHandler *texture.Handler, blocksStorage *blocks.Storage, world *world.World) *Dispatcher {
	return &Dispatcher{
		playerStorage:  playerStorage,
		textureHandler: textureHandler,
		blocksStorage:  blocksStorage,
		world:          world,
	}
}

func (d *Dispatcher) Dispatch(conn net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, pktData []byte) {
	switch pktCategory {
	case packet.CategoryWorld:
		switch pktOpcode {
		case packet.OpcodeWorldSetBlock:
			pktStructure := packet.WorldSetBlock{}

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal world set material data: %v\n", err)
			}

			d.world.AddBlock(uint8(pktStructure.BlockID), 0, int16(pktStructure.X), int16(pktStructure.Y))
			log.Printf("Игрок поставил материал %d на %d %d\n", pktStructure.BlockID, pktStructure.X, pktStructure.Y)
		default:
			log.Println("Неизвестный опкод")
		}
	case packet.CategoryPlayer:
		switch pktOpcode {
		/*
			case packet.OpcodePlayerHandshake:
				pkt-Structure := packet.PlayerHandshake{}

				if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
					log.Printf("Error! Cant unmarshal player handshake data: %v\n", err)
				}

				player := player.NewPlayer(conn, pktStructure.Username, float32(common.WorldSize/2), float32(common.WorldSize)/2)
				d.playerStorage.AddPlayer(player)

				if _, err := d.playerStorage.LoadPlayer(pktStructure.Username); err != nil {
					log.Printf("Error! Cant load player %s from database: %v\n", pktStructure.Username, err)
				}

				log.Printf("Hi, %s!\n", pktStructure.Username)
		*/
		case packet.OpcodeConnectRequest:
			// Нужно написать логику по которой
			// на такой пакет будет отправляться
			// в ответ все данные с сервера для игры
			reqStructure := packet.ConnectRequest{}

			if err := binary.Unmarshal(pktData, &reqStructure); err != nil {
				log.Printf("Error! cant unmarshal player connect request: %v\n", err)
			}

			player := player.NewPlayer(conn, reqStructure.Username, 512, 512)
			d.playerStorage.AddPlayer(player)

			log.Printf("Привет, %s!\n", reqStructure.Username)

			textureData, err := d.textureHandler.GetTextures()
			if err != nil {
				log.Printf("Error! Cant get texture in binary format: %v\n", err)
			}
			blocksData := d.blocksStorage.GetBlocks()

			resStructure := packet.ConnectResponse{
				TexturesData: textureData,
				BlocksData:   blocksData,
			}

			binaryStructure, err := binary.Marshal(&resStructure)
			if err != nil {
				log.Printf("Error! Cant marshal request structure to binary format: %v\n", err)
			}

			pktStructure := packet.Packet{
				Category: packet.CategoryConnection,
				Opcode:   packet.OpcodeConnectResponse,
				Payload:  binaryStructure,
			}

			binaryPkt, err := binary.Marshal(&pktStructure)
			if err != nil {
				log.Printf("Error! Cant marshal connect response packet to binary format: %v\n", err)
			}

			compressedPacket, err := compress.CompressPacket(binaryPkt)
			if err != nil {
				log.Printf("Error! Cant compress response packet: %v\n", err)
			}

			if _, err := conn.Write(compressedPacket); err != nil {
				log.Printf("Error! Cant write connection data to connection: %v\n", err)
			}
		case packet.OpcodePlayerMove:
			pktStructure := packet.PlayerMove{}

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal player move packet data: %v\n", err)
			}

			player := d.playerStorage.GetPlayer(conn)
			if player != nil {
				player.CurrentX = pktStructure.CurrentX
				player.CurrentY = pktStructure.CurrentY
				player.Flipped = pktStructure.Flipped
			} else {
				log.Printf("Error! Cant handle opcode move")
				conn.Close()
			}
		case packet.OpcodePlayerDisconnect:
			log.Printf("Info! Player was disconnected: %s", conn.RemoteAddr().String())
			d.playerStorage.RemovePlayer(conn)
		}
	case packet.CategoryInventory:
	}
}
