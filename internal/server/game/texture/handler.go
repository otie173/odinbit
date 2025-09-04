package texture

type Handler struct {
	textures *TexturePack
}

func NewHandler(textures *TexturePack) *Handler {
	return &Handler{
		textures: textures,
	}
}

func (h *Handler) GetTextures() ([]byte, error) {
	data, err := h.textures.GetTextures()
	if err != nil {
		return nil, err
	}
	return data, nil
}
