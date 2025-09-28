package inventory

type InventoryItem struct {
	Opened bool
	Count  int16
}

type InventoryPage struct {
	Items [9]InventoryItem
}

type Inventory struct {
	Id             int
	InventoryPages [3]InventoryPage
}
