package world

import (
	"log"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/game/player"
)

type Renderer struct {
	World   *World
	players player.Storage
}

func NewRenderer(world *World, players player.Storage) *Renderer {
	return &Renderer{
		World:   world,
		players: players,
	}
}

func (r *Renderer) compressPacket(binaryPacket []byte) ([]byte, error) {
	compressedData, err := minlz.Encode(nil, binaryPacket, minlz.LevelSmallest)
	if err != nil {
		return nil, err
	}
	return compressedData, nil
}

func (r *Renderer) Render() {
	players := r.players.GetPlayers()
	if len(players) > 0 {
		for _, player := range players {
			binaryOverworldArea, area, err := r.World.GetWorldArea(player.CurrentX, player.CurrentY)
			if err != nil {
				log.Printf("Error! Cant get binary overworld area: %v\n", err)
			}

			pktStructure := packet.WorldUpdate{
				Blocks: binaryOverworldArea,
				StartX: int16(area.StartX),
				StartY: int16(area.StartY),
				EndX:   int16(area.EndX),
				EndY:   int16(area.EndY),
			}

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

			compressedPkt, err := r.compressPacket(binaryPkt)
			if err != nil {
				log.Printf("Error! Cant compress world update packet: %v\n", err)
			}

			if _, err := player.Conn.Write(compressedPkt); err != nil {
				log.Printf("Error! Cant send binary packet of world area to player: %v\n", err)
			}
		}
	}
}
