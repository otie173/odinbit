package player

import (
	"log"
	"net"
	"os"
)

type Storage interface {
	AddPlayer(player *Player)
	RemovePlayer(player *Player)
}

type filestorage struct {
	players    []Player
	playersMap map[net.Conn]Player
}

func NewStorage(capacity int) *filestorage {
	players := make([]Player, capacity)
	playersMap := make(map[net.Conn]Player, capacity)

	dirPath := "players"
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			log.Printf("Error! Cant create directory for player storage: %v\n", err)
		}
	}

	return &filestorage{
		players:    players,
		playersMap: playersMap,
	}
}

func (s *filestorage) AddPlayer(player *Player) {

}

func (s *filestorage) RemovePlayer(player *Player) {

}
