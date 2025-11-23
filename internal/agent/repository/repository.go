package agentRepository

import "context"

type IRepository interface {
	CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
}
