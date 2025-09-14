package packet

type ServerTexture struct {
	Id   uint8
	Path string
}

type ServerBlock struct {
	TextureID uint8
	Passable  bool
}
