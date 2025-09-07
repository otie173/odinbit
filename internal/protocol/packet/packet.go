package packet

import (
	"github.com/otie173/odinbit/internal/server/game/world"
)

type PacketCategory int
type PacketOpcode int

type Packet struct {
	Category PacketCategory
	Opcode   PacketOpcode
	Payload  []byte
}

// Enum for packet categories
const (
	Texture PacketCategory = iota
	World
	Player
	Inventory
)

// Enum for packet opcodes
const (
// Texture opcodes
// SOON

// World opcodes
// SOON

// Player opcodes
// SOON

// Inventory opcodes
// SOON
)

type Ping struct{}

type Handshake struct {
	Username string
}

type GetTextures struct {
	Textures map[string]ServerTexture
}

type UpdateWorld struct {
	Blocks []world.Block
}
