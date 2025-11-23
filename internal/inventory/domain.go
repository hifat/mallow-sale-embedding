package inventoryModule

type ReqInventory struct {
	Inventories []Response
	Embeddings  [][]float32
}

type Response struct {
	ID   string `fake:"{uuid}" json:"id"`
	Name string `fake:"{name}" json:"name"`
}
