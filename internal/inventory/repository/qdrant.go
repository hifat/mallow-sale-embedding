package inventoryRepository

import (
	"context"
	"log"

	"github.com/google/uuid"
	inventoryModule "github.com/hifat/mallow-sale-embedding/internal/inventory"
	"github.com/qdrant/go-client/qdrant"
)

type qdrantRepository struct {
	db *qdrant.Client
}

func NewQdrant(db *qdrant.Client) IRepository {
	return &qdrantRepository{
		db,
	}
}

func (r *qdrantRepository) Search(ctx context.Context, queryEmb [][]float32) (*inventoryModule.Response, error) {
	limit := uint64(1)
	scThreshold := float32(0.8)
	searchResults, err := r.db.Query(ctx, &qdrant.QueryPoints{
		CollectionName: InventoryCol,
		Query:          qdrant.NewQuery(queryEmb[0]...),
		Limit:          &limit,
		ScoreThreshold: &scThreshold,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		log.Fatalf("failed to search: %v", err)
	}

	if len(searchResults) < 1 {
		return nil, nil
	}

	result := searchResults[0]

	if result.Payload != nil {
		return nil, nil
	}

	return &inventoryModule.Response{
		ID:   result.Payload["id"].GetStringValue(),
		Name: result.Payload["name"].GetStringValue(),
	}, nil
}

func (r *qdrantRepository) Upsert(ctx context.Context, req *inventoryModule.ReqInventory) error {
	points := make([]*qdrant.PointStruct, len(req.Embeddings))
	for i, emb := range req.Embeddings {
		point := &qdrant.PointStruct{
			Vectors: qdrant.NewVectors(emb...),
			Payload: qdrant.NewValueMap(map[string]any{
				"id":   req.Inventories[i].ID,
				"name": req.Inventories[i].Name,
			}),
		}

		if req.Inventories[i].ID != "" {
			point.Id = qdrant.NewIDUUID(req.Inventories[i].ID)
		} else {
			point.Id = qdrant.NewIDUUID(uuid.New().String())
		}

		points[i] = point
	}

	_, err := r.db.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: InventoryCol,
		Points:         points,
	})

	return err
}

func (r *qdrantRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: InventoryCol,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: []*qdrant.PointId{
						qdrant.NewIDUUID(id),
					},
				},
			},
		},
	})

	return err
}
