package player

import (
	"log"
	"net"
	"sync"

	"github.com/jmoiron/sqlx"
)

type PlayerModel struct {
	Id   int     `db:"player_id"`
	Name string  `db:"player_name"`
	X    float32 `db:"player_x"`
	Y    float32 `db:"player_y"`
}

type Storage interface {
	GetPlayer(conn net.Conn) *Player
	GetPlayers() []*Player
	AddPlayer(player *Player)
	RemovePlayer(playerConn net.Conn)
}

type storage struct {
	db         *sqlx.DB
	mu         sync.Mutex
	players    []*Player
	capacity   int
	playersMap map[net.Conn]*Player
}

type filestorage struct {
	capacity   int
	mu         sync.Mutex
	players    []*Player
	playersMap map[net.Conn]*Player
}

func NewStorage(db *sqlx.DB, capacity int) *storage {
	players := make([]*Player, 0, capacity)
	playersMap := make(map[net.Conn]*Player, capacity)

	return &storage{
		db:         db,
		players:    players,
		capacity:   capacity,
		playersMap: playersMap,
	}
}

func (s *storage) GetPlayer(conn net.Conn) *Player {
	return s.playersMap[conn]
}

func (s *storage) GetPlayers() []*Player {
	return s.players
}

func (s *storage) LoadPlayer(name string) *Player {
	var player PlayerModel

	if err := s.db.Get(&player, "SELECT * FROM Player WHERE player_name = ?", name); err != nil {
		log.Printf("Error! Cant get player from database: %v\n", err)
	}
}

func (s *storage) AddPlayer(player *Player) {
	s.mu.Lock()
	if len(s.players) < s.capacity {
		s.players = append(s.players, player)
		s.playersMap[player.Conn] = player
	} else {
		log.Println("Error! Cant add new player to storage: no more space")
	}
	s.mu.Unlock()
}

func (s *storage) RemovePlayer(playerConn net.Conn) {
	removedPlayer := s.playersMap[playerConn]
	players := make([]*Player, 0, s.capacity)

	s.mu.Lock()
	for _, player := range s.players {
		if player != removedPlayer {
			players = append(players, player)
		}
	}
	s.players = players
	s.mu.Unlock()
}
