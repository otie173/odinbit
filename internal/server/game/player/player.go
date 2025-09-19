package player

import (
	"net"
)

type Player struct {
	Conn     net.Conn
	Name     string
	CurrentX float32
	CurrentY float32
}

func NewPlayer(conn net.Conn, name string, x, y float32) *Player {
	player := &Player{
		Conn:     conn,
		Name:     name,
		CurrentX: x,
		CurrentY: y,
	}

	return player
}
