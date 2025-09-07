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

// Enum for texture packet opcodes
const (
	OpcodeTextureData PacketOpcode = 1000 + iota
)

// Enum for world packet opcodes
const (
	OpcodeSetBlock PacketOpcode = 2000 + iota
)
