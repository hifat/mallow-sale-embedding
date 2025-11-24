package main

import (
	"context"
	"log"

	inventoryRepository "github.com/hifat/mallow-sale-embedding/internal/inventory/repository"
	"github.com/hifat/mallow-sale-embedding/pkg/config"
	vectordb "github.com/hifat/mallow-sale-embedding/pkg/vectorDB"
	"github.com/qdrant/go-client/qdrant"
)

func main() {
	cfg, err := config.LoadConfig("./env/.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	qdClient, qdCleanup, err := vectordb.ConnectQdrant(&cfg.QDB)
	if err != nil {
		log.Fatalf("failed to connect qdrant: %v", err)
	}
	defer qdCleanup()

	qdClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: inventoryRepository.InventoryCol,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			// Size:     uint64(len(embs[0])), Wait mls gRPC
			Distance: qdrant.Distance_Cosine,
		}),
	})
}
