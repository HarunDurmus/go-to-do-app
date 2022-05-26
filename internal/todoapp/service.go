package todoapp

import (
	"context"
	"go.uber.org/zap"
)

const earhRadiusKm = 6371

type (
	Repository interface {
		UpsertInitialData(ctx context.Context) error
		GetAll(ctx context.Context) (locations []*Location, err error)
	}

	Service struct {
		repository Repository
		logger     *zap.Logger
	}
)

func NewService(repository Repository, logger *zap.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) InitializeData(ctx context.Context) error {

	return s.repository.UpsertInitialData(ctx)
}

func (s *Service) CreateTaskData(ctx context.Context, model ToDoList) (*ToDoList, error) {
	return nil, nil
}
