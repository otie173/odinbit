package player

type Player struct {
	Id   int
	Name int
}

func New() *Player {
	return &Player{}
}
