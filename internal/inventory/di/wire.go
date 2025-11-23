//go:build wireinject
// +build wireinject

package inventoryDi

import (
	"github.com/google/wire"
	agentRepository "github.com/hifat/mallow-sale-embedding/internal/agent/repository"
	inventoryHandler "github.com/hifat/mallow-sale-embedding/internal/inventory/handler"
	inventoryRepository "github.com/hifat/mallow-sale-embedding/internal/inventory/repository"
	inventoryService "github.com/hifat/mallow-sale-embedding/internal/inventory/service"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/ollama"
)

func Init(llm *ollama.LLM, db *qdrant.Client) *inventoryHandler.Handler {
	wire.Build(
		// Repository
		agentRepository.NewOllama,
		inventoryRepository.NewQdrant,

		// Service
		inventoryService.New,

		// Handler
		inventoryHandler.NewGrpc,
		inventoryHandler.New,
	)

	return &inventoryHandler.Handler{}
}
