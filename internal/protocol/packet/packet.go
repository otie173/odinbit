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
	CategoryTexture PacketCategory = iota
	CategoryWorld
	CategoryPlayer
	CategoryInventory
)

// Enum for texture opcodes
const (
	OpcodeTextureData PacketOpcode = 1000 + iota
)

// Enum for world opcodes
const (
	OpcodeWorldUpdate PacketOpcode = 2000 + iota
	OpcodeSetBlock
)

// Enum for player opcode
const (
	OpcodeHandshake PacketOpcode = 3000 + iota
	OpcodeMove
)
