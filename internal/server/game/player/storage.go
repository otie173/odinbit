package player

import (
	"log"
	"net"
	"os"
	"sync"
)

type Storage interface {
	GetPlayers() []*Player
	AddPlayer(player *Player)
	RemovePlayer(player *Player)
}

type filestorage struct {
	capacity   int
	mu         sync.Mutex
	players    []*Player
	playersMap map[net.Conn]*Player
}

func NewStorage(capacity int) *filestorage {
	players := make([]*Player, 0, capacity)
	playersMap := make(map[net.Conn]*Player, capacity)

	dirPath := "players"
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			log.Printf("Error! Cant create directory for player storage: %v\n", err)
		}
	}

	return &filestorage{
		capacity:   capacity,
		players:    players,
		playersMap: playersMap,
	}
}

func (s *filestorage) GetPlayers() []*Player {
	return s.players
}

func (s *filestorage) AddPlayer(player *Player) {
	s.mu.Lock()
	if len(s.players) < s.capacity {
		s.players = append(s.players, player)
	} else {
		log.Println("Error! Cant add new player to storage: no more space")
	}
	s.mu.Unlock()
}

func (s *filestorage) RemovePlayer(player *Player) {

}
