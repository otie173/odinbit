package packet

import (
	"github.com/otie173/odinbit/internal/server/game/world"
)

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

type Ping struct{}

type Handshake struct {
	Username string
}

type GetTextures struct {
	Textures map[string]ServerTexture
}

type GetWorld struct {
	Blocks []world.Block
}
