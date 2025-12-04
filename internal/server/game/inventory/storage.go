package inventory

type InventoryModel struct {
	Id          int  `db:"id"`
	WoodCount   int  `db:"wood_count"`
	StoneCount  int  `db:"stone_count"`
	MetalCount  int  `db:"metal_count"`
	AxeOpen     bool `db:"axe_open"`
	PickaxeOpen bool `db:"pickaxe_open"`
	ShovelOpen  bool `db:"shovel_open"`
}
