package inventory

type Storage struct {
}

type storage struct{}

func NewStorage() *storage {
	return &storage{}
}
