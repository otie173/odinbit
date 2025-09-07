package packet

type OpcodeMapping interface {
	register()
}

type mapping struct {
	opcodesMapping map[PacketOpcode]any
}

func New() *mapping {
	return &mapping{}
}

func (m *mapping) Load() {
	m.register(packet.)
}

func (m *mapping) register(opcode PacketOpcode, structure any) {
	m.opcodesMapping[opcode] = structure
}

type RequestTexturesStruct struct {
}

type GetTexturesStruct struct {
	Textures []ServerTexture
}
