package material

import "github.com/otie173/odinbit/internal/server/game/blocks"

var (
	materials blocks.Materials
)

func LoadMaterials(loadedMaterials blocks.Materials) {
	materials = loadedMaterials
}

func GetMaterials() blocks.Materials {
	return materials
}
