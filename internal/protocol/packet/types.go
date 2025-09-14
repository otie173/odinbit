package packet

type ServerTexture struct {
	Id   int
	Path string
}

type ServerBlock struct {
	TextureID int
	Passable  bool
}
