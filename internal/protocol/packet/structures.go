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
	Blocks         []byte
	StartX, StartY int
	EndX, EndY     int
}
