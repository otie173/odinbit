package player

import "net"

type Player struct {
	Id   int
	Conn net.Conn
	Name int
}

func New() *Player {
	return &Player{}
}
