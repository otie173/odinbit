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

func (i *Inventory) setMaterial(material common.Material) {
	i.currentMaterial = material
}
