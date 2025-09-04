package player

import "net"

type Storage struct {
	players    []Player
	playersMap map[net.Conn]Player
}

func NewStorage(capacity int) *Storage {
	players := make([]Player, capacity)
	playersMap := make(map[net.Conn]Player, capacity)

	return &Storage{
		players:    players,
		playersMap: playersMap,
	}
}

func (s *Storage) addPlayer(player Player) {

}

func (s *Storage) removePlayer(player Player) {

}
