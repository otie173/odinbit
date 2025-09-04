package tcp

import "github.com/otie173/odinbit/internal/server/game/player"

type Broadcaster struct{}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

func (b *Broadcaster) Broadcast(data []byte, players []player.Player) {

}
