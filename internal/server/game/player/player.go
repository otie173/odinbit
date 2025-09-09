package player

import "net"

type Player struct {
	Id   int
	Conn net.Conn
	Name int
	X    int
	Y    int
}

func New() *Player {
	return &Player{}
}
