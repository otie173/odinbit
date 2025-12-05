package inventory

import "github.com/otie173/odinbit/internal/client/common"

type Handler struct {
	inventory *Inventory
}

func NewHandler(inventory *Inventory) *Handler {
	return &Handler{
		inventory: inventory,
	}
}

func (h *Handler) GetMaterial() common.Material {
	return h.inventory.getMaterial()
}

func (h *Handler) GetMaterialCount(material common.Material) int {
	return h.inventory.getMaterialCount(material)
}

func (h *Handler) SetMaterial(material common.Material) {
	h.inventory.setMaterial(material)
}
