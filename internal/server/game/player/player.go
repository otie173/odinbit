package player

import "net"

var (
	idCounter int = 0
)

type Player struct {
	Id                int
	Conn              net.Conn
	Name              string
	CurrentX, TargetX float32
	CurrentY, TargetY float32
}

func NewPlayer(conn net.Conn, name string, x, y float32) *Player {
	player := &Player{
		Id:       idCounter,
		Conn:     conn,
		Name:     name,
		CurrentX: x,
		CurrentY: y,
	}
	idCounter++

	return player
}
