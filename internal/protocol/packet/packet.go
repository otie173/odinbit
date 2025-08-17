package packet

type Packet int

const (
	Ping Packet = iota
	Auth
	GetTextures
	GetWorld
)
