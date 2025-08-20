package packet

type Packet struct {
	Type    PacketType
	Payload []byte
}

type PacketType int

const (
	PingType PacketType = iota
	HandshakeType
	GetTexturesType
	GetWorldType
)

type Handshake struct {
	Username string
}
