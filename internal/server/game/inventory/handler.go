package inventory

type Handler struct {
	inventory *Inventory
}

func NewHandler(inventory *Inventory) *Handler {
	return &Handler{
		inventory: inventory,
	}
}

func LoadInventory(playerID int) *Inventory {
	return &Inventory{}
}
