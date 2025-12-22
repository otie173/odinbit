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
	CurrentX float32
	CurrentY float32
	Flipped  byte
}

type PlayerUpdate struct {
	Players []byte
}

type WorldSetBlock struct {
	BlockID int
	X       int
	Y       int
}
