package todo

import (
	"go.uber.org/zap"
)

type Repository interface {
	GetById(ID string) (*Todo, error)
	List() ([]Todo, error)
	Upsert(ID string, todo Todo) error
	Delete(ID string) error
}

type Service struct {
	repository Repository
	logger     *zap.Logger
}

func NewService(repository Repository, logger *zap.Logger) Service {
	return Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) InsertOrUpdateTodo(todo Todo) error {
	s.logger.Sugar().Debugf("creating todo with plate: %s", todo.ID)
	err := s.repository.Upsert("asdasd", todo)
	s.logger.Sugar().Debugf("todo created with plate: %s", todo.ID)
	if err != nil {
		s.logger.Sugar().Errorf("error creating todo with plate: %s", todo.ID)
		return err
	}
	s.logger.Sugar().Debugf("created todo with plate: %s", todo.ID)
	return nil
}
