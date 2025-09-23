package player

import (
	"log"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/server/common"
)

var (
	PlayerMoved int32
	PlayerMu    sync.Mutex
	GamePlayer  = Player{
		Name:     "otie173",
		CurrentX: float32(common.WorldSize / 2),
		CurrentY: float32(common.WorldSize / 2),
	}
	NetworkPlayers    []Player = make([]Player, 0, 16)
	NetworkPlayersRaw []Player = make([]Player, 0, 16)
	NetPlayersMu      sync.Mutex
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

func UpdateNetworkPlayers() {
	NetPlayersMu.Lock()
	defer NetPlayersMu.Unlock()

	delta := rl.GetFrameTime()

	if len(NetworkPlayers) == 0 && len(NetworkPlayersRaw) > 0 {
		NetworkPlayers = make([]Player, len(NetworkPlayersRaw))
		copy(NetworkPlayers, NetworkPlayersRaw)
		return
	}

	if len(NetworkPlayers) == 0 && len(NetworkPlayersRaw) == 0 {
		return
	}

	for i := range NetworkPlayers {
		for _, rawPlayer := range NetworkPlayersRaw {
			diffX := rawPlayer.CurrentX - NetworkPlayers[i].CurrentX
			diffY := rawPlayer.CurrentY - NetworkPlayers[i].CurrentY

			NetworkPlayers[i].CurrentX += diffX * delta * 8.0
			NetworkPlayers[i].CurrentY += diffY * delta * 8.0
			NetworkPlayers[i].Flipped = rawPlayer.Flipped
			break
		}
	}

}

func DrawNetworkPlayers() {
	if len(NetworkPlayers) > 0 {
		NetPlayersMu.Lock()
		for _, netPlayer := range NetworkPlayers {
			var playerRec rl.Rectangle
			if netPlayer.Flipped == 0 {
				playerRec = rl.NewRectangle(0, 0, -12, 12)
			} else {
				playerRec = rl.NewRectangle(0, 0, 12, 12)
			}

			playerVec := rl.NewVector2(netPlayer.CurrentX*12, netPlayer.CurrentY*12)
			rl.DrawTextureRec(texture.PlayerTexture, playerRec, playerVec, rl.White)
		}
		NetPlayersMu.Unlock()
	}
}

func DrawPlayer() {
	// log.Println("Рисую локального игрока:", GamePlayer.CurrentX, GamePlayer.CurrentY)
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
