package inventory

type Inventory struct {
	Id         int
	WoodCount  int
	StoneCount int
	MetalCount int
}

func NewInventory() *Inventory {
	return &Inventory{}
}
