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
	CategoryMob
	CategoryInventory
)

// Enum for texture opcodes
const (
	OpcodeTextureData PacketOpcode = 1000 + iota
)

// Enum for world opcodes
const (
	OpcodeWorldUpdate PacketOpcode = 2000 + iota
	OpcodeWorldSetMaterial
	OpcodeWorldSetBlock
)

// Enum for player opcode
const (
	OpcodePlayerHandshake PacketOpcode = 3000 + iota
	OpcodePlayerMove
	OpcodePlayerUpdate
	OpcodePlayerDisconnect
)

// Enum for mob opcode
const (
// OpcodePlayerHandshake PacketOpcode = 4000 + iota
// OpcodePlayerMove
// OpcodePlayerUpdate
// OpcodePlayerDisconnect
)
