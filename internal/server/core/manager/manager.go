package manager

import (
	"log"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/common"
	"github.com/otie173/odinbit/internal/server/core/ticker"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
	"github.com/vmihailenco/msgpack/v5"
)

type Components struct {
	// Game things
	Textures  *texture.TexturePack
	Overworld *world.World
	Players   player.Storage

	// Network things
	Handler     *http.Handler
	Listener    *tcp.Listener
	Broadcaster *tcp.Broadcaster

	// System things
	Ticker *ticker.Ticker
}

type Manager struct {
	Components Components
}

func New(components Components) *Manager {
	return &Manager{
		Components: components,
	}
}

func (m *Manager) HandleNetwork(httpAddr, tcpAddr string) {
	go func() {
		if err := m.Components.Handler.Run(httpAddr); err != nil {
			log.Fatalf("Error! Cant run HTTP handler: %v\n", err)
		}
	}()
	log.Printf("HTTP handler running on: %s\n", httpAddr)

	go func() {
		if err := m.Components.Listener.Run(tcpAddr); err != nil {
			log.Fatalf("Error! Cant run TCP listener: %v\n", err)
		}
	}()
	log.Printf("TCP listener running on: %s\n", tcpAddr)
}

func (m *Manager) HandleGame() {
	m.Components.Ticker.Run(func() {
		players := m.Components.Players.GetPlayers()
		for _, player := range players {
			binaryOverworldArea, err := m.Components.Overworld.GetWorldArea(player.X, player.Y)
			if err != nil {
				log.Printf("Error! Cant get binary overworld area: %v\n", err)
			}

			pktStructure := packet.WorldUpdate{
				Blocks: binaryOverworldArea,
				StartX: player.X - common.ViewRadius,
				StartY: player.Y - common.ViewRadius,
				EndX:   player.X + common.ViewRadius,
				EndY:   player.Y + common.ViewRadius,
			}
			log.Println(pktStructure.StartX, pktStructure.StartY, pktStructure.EndX, pktStructure.EndY, player.X, player.Y)

			binaryStructure, err := msgpack.Marshal(&pktStructure)
			if err != nil {
				log.Printf("Error! Cant marshal world update structure to binary format: %v\n", err)
			}

			pkt := packet.Packet{
				Category: packet.CategoryWorld,
				Opcode:   packet.OpcodeWorldUpdate,
				Payload:  binaryStructure,
			}

			binaryPkt, err := msgpack.Marshal(&pkt)
			if err != nil {
				log.Printf("Error! Cant marshal world update packet: %v\n", err)
			}

			if _, err := player.Conn.Write(binaryPkt); err != nil {
				log.Printf("Error! Cant send binary packet of world area to player: %v\n", err)
			}
		}
	})
}
