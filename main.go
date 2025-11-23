package main

import (
	"context"
	"log"
	"os"

	inventoryDi "github.com/hifat/mallow-sale-embedding/internal/inventory/di"
	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/ollama"
)

var OllamaHost string
var AgentToken string
var QdHost string
var QdPort string
var QdApiKey string

func init() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	OllamaHost = os.Getenv("OLLAMA_HOST")
	AgentToken = os.Getenv("AGENT_TOKEN")
	QdHost = os.Getenv("QD_HOST")
	QdPort = os.Getenv("QD_PORT")
	QdApiKey = os.Getenv("QD_API_KEY")
}

func main() {
	llm, err := ollama.New(
		ollama.WithModel("paraphrase-multilingual"),
		ollama.WithServerURL(OllamaHost),
	)
	if err != nil {
		log.Fatalf("failed to new model: %v", err)
	}

	ctx := context.Background()

	// Create embeddings for the text
	texts := []string{
		"ฐานข้อมูลแบบฝัง",
		"ฉันรักคุณ",
		"สวัสดีตอนเช้า",
	}

	embs, err := llm.CreateEmbedding(ctx, texts)
	if err != nil {
		log.Fatalf("failed to create embeddings: %v", err)
	}

	// Connect to Qdrant
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host:   QdHost,
		Port:   6334,
		APIKey: QdApiKey,
		UseTLS: true,
		// Cloud:  true,

	})
	if err != nil {
		log.Fatalf("failed to connect to qdrant: %v", err)
	}
	defer qdClient.Close()

	colName := "embeddings"

	// // Ignore error if collection doesn't exist
	// _, err = collectionClient.Delete(ctx, &qdrant.DeleteCollection{
	// 	CollectionName: colName,
	// })

	// Create or recreate collection
	qdClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: colName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(len(embs[0])),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	// Insert embeddings into Qdrant
	points := make([]*qdrant.PointStruct, len(embs))
	for i, emb := range embs {
		points[i] = &qdrant.PointStruct{
			Id: qdrant.NewIDNum(uint64(i + 1)),
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{Data: emb},
				},
			},
			Payload: map[string]*qdrant.Value{
				"text": {Kind: &qdrant.Value_StringValue{StringValue: texts[i]}},
			},
		}
	}

	_, err = qdClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: colName,
		Points:         points,
	})
	if err != nil {
		log.Fatalf("failed to upsert points: %v", err)
	}

	log.Println("✓ Embeddings stored in Qdrant")

	queryText := "เราเลิกกันเถอะ"
	queryEmb, err := llm.CreateEmbedding(ctx, []string{queryText})
	if err != nil {
		log.Fatalf("failed to create query embedding: %v", err)
	}

	limit := uint64(1)
	scThreshold := float32(0.8)
	searchResults, err := qdClient.Query(ctx, &qdrant.QueryPoints{
		CollectionName: colName,
		Query:          qdrant.NewQuery(queryEmb[0]...),
		Limit:          &limit,
		ScoreThreshold: &scThreshold,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		log.Fatalf("failed to search: %v", err)
	}

	if len(searchResults) > 0 {
		result := searchResults[0]
		if result.Payload != nil {
			if textValue, ok := result.Payload["text"]; ok {
				if stringVal, ok := textValue.Kind.(*qdrant.Value_StringValue); ok {
					log.Printf("Query: '%s' -> Result: '%s'\n", queryText, stringVal.StringValue)
				}
			}
		}
	} else {
		log.Println("No results found")
	}

	inventoryDi.Init(llm, qdClient)
}
