package inventoryHandler

type Handler struct {
	InventoryGrpc *InventoryGrpc
}

func New(inventoryGrpc *InventoryGrpc) *Handler {
	return &Handler{inventoryGrpc}
}
