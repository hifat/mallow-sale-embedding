package vectordb

import (
	"github.com/hifat/mallow-sale-embedding/pkg/config"
	"github.com/qdrant/go-client/qdrant"
)

func ConnectQdrant(cfg *config.QDB) (*qdrant.Client, func(), error) {
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.Host,
		Port:   cfg.Port,
		APIKey: cfg.ApiKey,
		UseTLS: true,
		// Cloud:  true,
	})
	if err != nil {
		return nil, nil, err
	}

	return qdClient, func() {
		qdClient.Close()
	}, nil
}
