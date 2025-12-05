package inventory

type InventoryModel struct {
	Id         int `db:"id"`
	WoodCount  int `db:"wood_count"`
	StoneCount int `db:"stone_count"`
	MetalCount int `db:"metal_count"`
}
