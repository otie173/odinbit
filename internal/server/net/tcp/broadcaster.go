package tcp

import (
	"log"

	"github.com/otie173/odinbit/internal/server/game/player"
)

type Broadcaster struct{}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

func (b *Broadcaster) Broadcast(data []byte, players []*player.Player) {
	for _, player := range players {
		if _, err := player.Conn.Write(data); err != nil {
			log.Printf("Error! Cant broadcast data to player: %v\n", err)
		}
	}
}
