package texture

type Handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) GetTextures() ([]byte, error) {
	data, err := h.storage.GetTextures()
	if err != nil {
		return nil, err
	}
	return data, nil
}
