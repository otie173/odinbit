package player

import "net"

var (
	idCounter int = 0
)

type Player struct {
	Id   int
	Conn net.Conn
	Name string
	X    int16
	Y    int16
}

func NewPlayer(conn net.Conn, name string, x, y int16) *Player {
	player := &Player{
		Id:   idCounter,
		Conn: conn,
		Name: name,
		X:    x,
		Y:    y,
	}
	idCounter++

	return player
}
