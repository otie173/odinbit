package player

import (
	"log"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kelindar/binary"
	"github.com/otie173/odinbit/internal/client/net"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/common"
)

var (
	PlayerMoved int32
	PlayerMu    sync.Mutex
	GamePlayer  = Player{
		Name:     "otie173",
		CurrentX: common.WorldSize / 2,
		CurrentY: common.WorldSize / 2,
	}
	NetworkPlayers []Player = make([]Player, 0, 16)
	NetConnection  *net.Handler
)

const (
	LocalPlayer PlayerType = iota
	NetworkPlayer
)

type PlayerType byte

type Player struct {
	Name     string
	CurrentX float32
	CurrentY float32
	Flipped  byte
}

func AddNetworkPlayer(player Player) {
	NetworkPlayers = append(NetworkPlayers, player)
	log.Println("Player was added: ", player)
}

func RemoveNetworkPlayer(removedPlayer Player) {
	log.Println("Player was removed: ", removedPlayer)
	players := make([]Player, 0, 16)

	for _, player := range NetworkPlayers {
		if player != removedPlayer {
			players = append(players, player)
		}
	}
	NetworkPlayers = players
}

func DrawNetworkPlayers() {
	playerRec := rl.NewRectangle(0, 0, 12, 12)

	if len(NetworkPlayers) > 0 {
		for _, player := range NetworkPlayers {
			playerVec := rl.NewVector2(player.CurrentX*12, player.CurrentY*12)
			rl.DrawTextureRec(texture.PlayerTexture, playerRec, playerVec, rl.White)
		}
	}
}

func UpdateServerPos() {
	PlayerMu.Lock()
	pktStructure := packet.PlayerMove{
		CurrentX: GamePlayer.CurrentX,
		CurrentY: GamePlayer.CurrentY,
	}
	PlayerMu.Unlock()

	binaryStructure, err := binary.Marshal(&pktStructure)
	if err != nil {
		log.Printf("Error! Cant marshal player move structure: %v\n", err)
	}

	pkt := packet.Packet{
		Category: packet.CategoryPlayer,
		Opcode:   packet.OpcodeMove,
		Payload:  binaryStructure,
	}

	data, err := binary.Marshal(&pkt)
	if err != nil {
		log.Printf("Error! Cant marshal player move packet: %v\n", err)
	}

	compressedPkt, err := net.CompressPkt(data)
	if err != nil {
		log.Printf("Error! Cant compress binary player move packet: %v\n", err)
	}

	if err := NetConnection.Write(compressedPkt); err != nil {
		log.Printf("Error! Cant write player move packet data to server: %v\n", err)
	}
}

func DrawPlayer() {
	var playerRec rl.Rectangle
	if GamePlayer.Flipped == 0 {
		playerRec = rl.NewRectangle(0, 0, -12, 12)
	} else {
		playerRec = rl.NewRectangle(0, 0, 12, 12)
	}

	PlayerMu.Lock()
	playerVec := rl.NewVector2(GamePlayer.CurrentX*12, GamePlayer.CurrentY*12)
	PlayerMu.Unlock()
	rl.DrawTextureRec(texture.PlayerTexture, playerRec, playerVec, rl.White)
	//rl.DrawTexture(texture.PlayerTexture, int32(GamePlayer.CurrentX*12), int32(GamePlayer.CurrentY*12), rl.White)
}
