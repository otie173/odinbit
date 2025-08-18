package world

type block struct {
	textureID int
	passable  bool
}

type World struct {
	storage   *storage
	generator *generator
}

func New() *World {
	storage := newStorage()
	generator := newGenerator(storage)

	return &World{storage: storage, generator: generator}
}

func (w *World) AddBlock(id int, passable bool, x, y int) {
	w.storage.addBlock(id, passable, x, y)
}

func (w *World) RemoveBlock(x, y int) {
	w.storage.removeBlock(x, y)
}

func (w *World) Generate() {
	w.generator.generateWorld()
}
