package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hifat/mallow-sale-api/pkg/grpc/inventoryProto"
	inventoryModule "github.com/hifat/mallow-sale-embedding/internal/inventory"
	inventoryRepository "github.com/hifat/mallow-sale-embedding/internal/inventory/repository"
	"github.com/hifat/mallow-sale-embedding/pkg/config"
	"github.com/hifat/mallow-sale-embedding/pkg/vectordb"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/ollama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func apiKeyInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", apiKey)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

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

	grpcConn, err := grpc.NewClient(
		cfg.MLS.GRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(apiKeyInterceptor(cfg.MLS.GRPCKey)),
	)
	if err != nil {
		log.Fatalf("failed to connect grpc: %v", err)
	}
	defer grpcConn.Close()

	llm, err := ollama.New(
		ollama.WithModel("paraphrase-multilingual"),
		ollama.WithServerURL(cfg.Agent.OllamaHost),
	)
	if err != nil {
		log.Fatalf("failed to new model: %v", err)
	}

	ctx := context.Background()

	inventoryGrpc := inventoryProto.NewInventoryGrpcServiceClient(grpcConn)
	inventorRes, err := inventoryGrpc.Find(ctx, &inventoryProto.Query{})
	if err != nil {
		log.Fatalf("failed to find inventory: %v", err)
	}

	texts := make([]string, len(inventorRes.Items))
	for i, v := range inventorRes.Items {
		texts[i] = v.Name
	}

	fmt.Println("creating embedding...")
	embs, err := llm.CreateEmbedding(ctx, texts)
	if err != nil {
		log.Fatalf("failed to create embeddings: %v", err)
	}
	fmt.Println("created embedding")

	if err := qdClient.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: inventoryRepository.InventoryCol,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(len(embs[0])),
			Distance: qdrant.Distance_Cosine,
		}),
	}); err != nil {
		log.Fatalf("failed to create collection: %v", err)
	}

	newInventories := make([]inventoryModule.Response, 0, len(inventorRes.Items))
	for _, v := range inventorRes.Items {
		newInventories = append(newInventories, inventoryModule.Response{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	reqInventory := inventoryModule.ReqInventory{
		Inventories: newInventories,
		Embeddings:  embs,
	}

	newInventoryRepo := inventoryRepository.NewQdrant(qdClient)
	if err := newInventoryRepo.Upsert(ctx, &reqInventory); err != nil {
		log.Fatalf("failed to upsert inventory: %v", err)
	}
}
