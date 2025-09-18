package tcp

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/vmihailenco/msgpack/v5"
)

type Dispatcher struct {
	playerStorage  player.Storage
	textureHandler *texture.Handler
	worldHandler   *world.Handler
}

func NewDispatcher(playerStorage player.Storage, textureHandler *texture.Handler, worldHandler *world.Handler) *Dispatcher {
	return &Dispatcher{
		playerStorage:  playerStorage,
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

			player := player.NewPlayer(conn, pktStructure.Username, common.WorldSize/2, common.WorldSize/2)
			d.playerStorage.AddPlayer(player)
			log.Printf("Hi, %s!\n", pktStructure.Username)
		case packet.OpcodeMove:
			//log.Printf("Info! Player with conn %s is move\n", conn.RemoteAddr())

			pktStructure := packet.PlayerMove{}

			if err := binary.Unmarshal(pktData, &pktStructure); err != nil {
				log.Printf("Error! Cant unmarshal player move packet data: %v\n", err)
			}

			player := d.playerStorage.GetPlayer(conn)
			player.CurrentX = pktStructure.CurrentX
			player.TargetX = pktStructure.TargetX
			player.CurrentY = pktStructure.CurrentY
			player.TargetY = pktStructure.TargetY

		}
	case packet.CategoryInventory:
	}
}
