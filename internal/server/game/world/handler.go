package world

import (
	"log"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/game/player"
)

type Handler struct {
	World   *World
	players player.Storage
}

func NewHandler(world *World, players player.Storage) *Handler {
	return &Handler{
		World:   world,
		players: players,
	}
}

func (h *Handler) compressPacket(binaryPacket []byte) ([]byte, error) {
	compressedData, err := minlz.Encode(nil, binaryPacket, minlz.LevelSmallest)
	if err != nil {
		return nil, err
	}
	return compressedData, nil
}

func (h *Handler) Handle() {
	players := h.players.GetPlayers()
	for _, player := range players {
		binaryOverworldArea, err := h.World.GetWorldArea(player.CurrentX, player.CurrentY)
		if err != nil {
			log.Printf("Error! Cant get binary overworld area: %v\n", err)
		}

		pktStructure := packet.WorldUpdate{
			Blocks: binaryOverworldArea,
			StartX: int16(player.CurrentX - common.ViewRadius),
			StartY: int16(player.CurrentY - common.ViewRadius),
			EndX:   int16(player.CurrentX + common.ViewRadius),
			EndY:   int16(player.CurrentY + common.ViewRadius),
		}
		log.Println(pktStructure.StartX, pktStructure.StartY, pktStructure.EndX, pktStructure.EndY, player.CurrentX, player.CurrentY)

		binaryStructure, err := binary.Marshal(&pktStructure)
		if err != nil {
			log.Printf("Error! Cant marshal world update structure to binary format: %v\n", err)
		}

		pkt := packet.Packet{
			Category: packet.CategoryWorld,
			Opcode:   packet.OpcodeWorldUpdate,
			Payload:  binaryStructure,
		}

		binaryPkt, err := binary.Marshal(&pkt)
		if err != nil {
			log.Printf("Error! Cant marshal world update packet: %v\n", err)
		}

		compressedPkt, err := h.compressPacket(binaryPkt)
		if err != nil {
			log.Printf("Error! Cant compress world update packet: %v\n", err)
		}

		if _, err := player.Conn.Write(compressedPkt); err != nil {
			log.Printf("Error! Cant send binary packet of world area to player: %v\n", err)
		}
	}
}
