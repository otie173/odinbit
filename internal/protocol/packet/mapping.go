package packet

type Mapping struct {
	opcodesMapping map[PacketOpcode]any
}

func New() *Mapping {
	opcodesMapping := make(map[PacketOpcode]any)

	return &Mapping{
		opcodesMapping: opcodesMapping,
	}
}

func (m *Mapping) Load() {

}

func (m *Mapping) register(opcode PacketOpcode, structure any) {
	m.opcodesMapping[opcode] = structure
}
