package inventoryHandler

import (
	"context"

	inventoryProto "github.com/hifat/mallow-sale-embedding/internal/inventory/proto"
	inventoryService "github.com/hifat/mallow-sale-embedding/internal/inventory/service"
)

type InventoryGrpc struct {
	inventoryProto.UnimplementedInventoryGrpcServiceServer
	inventorySvc inventoryService.IService
}

func NewGrpc(inventorySvc inventoryService.IService) *InventoryGrpc {
	return &InventoryGrpc{inventorySvc: inventorySvc}
}

func (g *InventoryGrpc) Search(ctx context.Context, req *inventoryProto.SearchReq) (*inventoryProto.InventoryResponse, error) {
	res, err := g.inventorySvc.Search(ctx, req.Search)
	if err != nil {
		return nil, err
	}

	return &inventoryProto.InventoryResponse{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}
