package world

type Handler struct {
	world *World
}

func NewHandler(world *World) *Handler {
	return &Handler{
		world: world,
	}
}
