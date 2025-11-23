package repository

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ollamaRepository struct {
	llm *ollama.LLM
}

func NewQdrant(llm *ollama.LLM, db *qdrant.Client) IRepository {
	return &ollamaRepository{
		llm,
	}
}

func (r *ollamaRepository) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
	return r.llm.CreateEmbedding(ctx, inputTexts)
}
