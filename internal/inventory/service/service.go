package service

import (
	"context"
	"log/slog"

	agentRepository "github.com/hifat/mallow-sale-embedding/internal/agent/repository"
	inventoryModule "github.com/hifat/mallow-sale-embedding/internal/inventory"
	inventoryRepository "github.com/hifat/mallow-sale-embedding/internal/inventory/repository"
)

type IService interface {
	Search(ctx context.Context, text string) (*inventoryModule.Response, error)
}

type service struct {
	agentRepo     agentRepository.IRepository
	inventoryRepo inventoryRepository.IRepository
}

func New(agentRepo agentRepository.IRepository, inventoryRepo inventoryRepository.IRepository) IService {
	return &service{
		agentRepo,
		inventoryRepo,
	}
}

func (s *service) Search(ctx context.Context, text string) (*inventoryModule.Response, error) {
	queryEmb, err := s.agentRepo.CreateEmbedding(ctx, []string{text})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	res, err := s.inventoryRepo.Search(ctx, queryEmb)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return res, nil
}
