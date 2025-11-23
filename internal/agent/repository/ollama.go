package agentRepository

import (
	"context"

	"github.com/tmc/langchaingo/llms/ollama"
)

type ollamaRepository struct {
	llm *ollama.LLM
}

func NewOllama(llm *ollama.LLM) IRepository {
	return &ollamaRepository{
		llm,
	}
}

func (r *ollamaRepository) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
	return r.llm.CreateEmbedding(ctx, inputTexts)
}
