package player

import "net"

type Storage interface {
	AddPlayer(player *Player)
	RemovePlayer(player *Player)
}

type StorageImpl struct {
	players    []Player
	playersMap map[net.Conn]Player
}

func NewStorage(capacity int) *StorageImpl {
	players := make([]Player, capacity)
	playersMap := make(map[net.Conn]Player, capacity)

	return &StorageImpl{
		players:    players,
		playersMap: playersMap,
	}
}

func (s *StorageImpl) AddPlayer(player *Player) {

}

func (s *StorageImpl) RemovePlayer(player *Player) {

}
