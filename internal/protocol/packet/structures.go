package packet

// Structures for create opcode mappings
type TextureRequest struct{}

type TextureData struct {
	Textures map[string]ServerTexture
}

type PlayerHandshake struct {
	Username string
}

type WorldUpdate struct {
	_msgpack       struct{} `msgpack:",as_array"`
	Blocks         []byte
	StartX, StartY int16
	EndX, EndY     int16
}
