package main

import (
	"fmt"
	"log"
	"net"

	inventoryDi "github.com/hifat/mallow-sale-embedding/internal/inventory/di"
	inventoryPb "github.com/hifat/mallow-sale-embedding/internal/inventory/pb"
	middlewareMiddleware "github.com/hifat/mallow-sale-embedding/internal/middleware/handler"
	"github.com/hifat/mallow-sale-embedding/pkg/config"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/ollama"
	"google.golang.org/grpc"
)

// var OllamaHost string
// var QdHost string
// var QdPort string
// var QdApiKey string
// var ApiKey string
// var GrpcPort string

// func init() {
// 	OllamaHost = os.Getenv("OLLAMA_HOST")
// 	QdHost = os.Getenv("QD_HOST")
// 	QdPort = os.Getenv("QD_PORT")
// 	QdApiKey = os.Getenv("QD_API_KEY")
// 	ApiKey = os.Getenv("API_KEY")
// 	GrpcPort = os.Getenv("GRPC_PORT")
// }

func main() {
	cfg, err := config.LoadConfig("./env/.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	agentCfg := cfg.Agent
	llm, err := ollama.New(
		ollama.WithModel("paraphrase-multilingual"),
		ollama.WithServerURL(agentCfg.OllamaHost),
	)
	if err != nil {
		log.Fatalf("failed to new model: %v", err)
	}

	qdCfg := cfg.QDB
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host:   qdCfg.Host,
		Port:   qdCfg.Port,
		APIKey: qdCfg.ApiKey,
		UseTLS: true,
		// Cloud:  true,
	})
	if err != nil {
		log.Fatalf("failed to connect to qdrant: %v", err)
	}
	defer qdClient.Close()

	m := middlewareMiddleware.New(&cfg.Auth)

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(m.AuthInterceptor),
	)

	ivtDi := inventoryDi.Init(llm, qdClient)

	inventoryPb.RegisterInventoryGrpcServiceServer(grpcSrv, ivtDi.InventoryGrpc)

	appCfg := cfg.App
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appCfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("listening on port :%d", appCfg.Port)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
