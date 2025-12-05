package inventory

import "github.com/otie173/odinbit/internal/client/common"

type Inventory struct {
	currentMaterial common.Material
	WoodCount       int
	StoneCount      int
	MetalCount      int
}

func NewInventory() *Inventory {
	return &Inventory{
		currentMaterial: -1,
		WoodCount:       0,
		StoneCount:      0,
		MetalCount:      0,
	}
}

func (i *Inventory) getMaterial() common.Material {
	return i.currentMaterial
}

func (i *Inventory) getMaterialCount(material common.Material) int {
	var count int = 0

	switch material {
	case common.Wood:
		count = i.WoodCount
	case common.Stone:
		count = i.StoneCount
	case common.Metal:
		count = i.MetalCount
	}

	return count
}

func (i *Inventory) setMaterial(material common.Material) {
	i.currentMaterial = material
}
