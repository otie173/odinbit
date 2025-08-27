package player

type Storage struct {
	storage []Player
}

func NewStorage(capacity int) *Storage {
	storage := make([]Player, capacity)

	return &Storage{
		storage: storage,
	}
}
