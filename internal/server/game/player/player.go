package player

import "net"

var (
	idCounter int = 0
)

type Player struct {
	Id   int
	Conn net.Conn
	Name string
	X    int
	Y    int
}

func NewPlayer(conn net.Conn, name string) *Player {
	player := &Player{
		Id:   idCounter,
		Name: name,
		X:    0,
		Y:    0,
	}
	idCounter++

	return player
}
