package inventoryRepository

import (
	"context"

	inventoryModule "github.com/hifat/mallow-sale-embedding/internal/inventory"
)

const InventoryCol string = "inventories"

type IRepository interface {
	Search(ctx context.Context, queryEmb [][]float32) (*inventoryModule.Response, error)
	Upsert(ctx context.Context, req *inventoryModule.ReqInventory) error
	BatchUpsert(ctx context.Context, reqs []*inventoryModule.ReqInventory, batchSize int) error
	DeleteByID(ctx context.Context, id string) error
}
