package player

import (
	"log"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	cplayer "github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
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
	players := h.storage.GetPlayers()

	for _, currentPlayer := range players {
		renderPlayers := make([]cplayer.Player, 0, len(players)-1)
		for _, otherPlayer := range players {
			if currentPlayer != otherPlayer {
				clientPlayer := cplayer.Player{
					Name:     otherPlayer.Name,
					CurrentX: otherPlayer.CurrentX,
					CurrentY: otherPlayer.CurrentY,
					Flipped:  otherPlayer.Flipped,
				}

				renderPlayers = append(renderPlayers, clientPlayer)
			}
		}

		pktStructure := packet.PlayerUpdate{}

		binaryPlayers, err := binary.Marshal(renderPlayers)
		if err != nil {
			log.Printf("Error! Cant marshal render players to binary format: %v\n", err)
		}
		pktStructure.Players = binaryPlayers

		binaryStructure, err := binary.Marshal(&pktStructure)
		if err != nil {
			log.Printf("Error! Cant marshal player update structure to binary format: %v\n", err)
		}

		pkt := packet.Packet{
			Category: packet.CategoryPlayer,
			Opcode:   packet.OpcodePlayerUpdate,
			Payload:  binaryStructure,
		}

		binaryPkt, err := binary.Marshal(&pkt)
		if err != nil {
			log.Printf("Error! Cant marshal player update packet to binary format: %v\n", err)
		}

		compressedPkt, err := h.compressPacket(binaryPkt)
		if err != nil {
			log.Printf("Error! Cant compress player update binary packet: %v\n", err)
		}

		if _, err := currentPlayer.Conn.Write(compressedPkt); err != nil {
			log.Printf("Error! Cant write player update binary packet to player: %v\n", err)
		}
	}
}
