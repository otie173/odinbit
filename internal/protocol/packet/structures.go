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
	StartX, StartY int16
	EndX, EndY     int16
}

type PlayerMove struct {
	CurrentX, TargetX float32
	CurrentY, TargetY float32
}
