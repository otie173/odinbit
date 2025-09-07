package packet

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
	GetTextures PacketOpcode = iota
	GetWorld
)
