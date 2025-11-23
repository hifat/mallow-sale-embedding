package inventoryHandler

import (
	"context"

	inventoryPb "github.com/hifat/mallow-sale-embedding/internal/inventory/pb"
	inventoryService "github.com/hifat/mallow-sale-embedding/internal/inventory/service"
)

type InventoryGrpc struct {
	inventoryPb.UnimplementedInventoryGrpcServiceServer
	inventorySvc inventoryService.IService
}

func NewGrpc(inventorySvc inventoryService.IService) *InventoryGrpc {
	return &InventoryGrpc{inventorySvc: inventorySvc}
}

func (g *InventoryGrpc) Search(ctx context.Context, req *inventoryPb.SearchReq) (*inventoryPb.InventoryResponse, error) {
	res, err := g.inventorySvc.Search(ctx, req.Search)
	if err != nil {
		return nil, err
	}

	return &inventoryPb.InventoryResponse{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}
