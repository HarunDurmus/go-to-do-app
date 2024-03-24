package todo

import (
	"go.uber.org/zap"
)

// Service represents the service layer for managing TODOs.
type Service struct {
	repository Repository  // repository is used to interact with the database
	logger     *zap.Logger // logger is used for logging
}

// NewService creates a new instance of Service with the given repository and logger.
func NewService(repository Repository, logger *zap.Logger) Service {
	return Service{
		repository: repository,
		logger:     logger,
	}
}

// InsertOrUpdateTodo inserts or updates a TODO item in the database.
func (s *Service) InsertOrUpdateTodo(todo Todo) error {
	// Log the creation of the TODO item
	s.logger.Sugar().Debugf("creating todo with ID: %s", todo.ID)

	// Perform the upsert operation in the repository
	err := s.repository.Upsert(todo.ID, todo)

	// Check for errors
	if err != nil {
		// Log the error
		s.logger.Sugar().Errorf("error creating todo with ID: %s", todo.ID)
		return err
	}

	// Log successful creation of the TODO item
	s.logger.Sugar().Debugf("todo created with ID: %s", todo.ID)
	return nil
}

// DeleteTodo deletes a TODO item from the database by its ID.
func (s *Service) DeleteTodo(ID string) error {
	// Log the deletion of the TODO item
	s.logger.Sugar().Debugf("deleting todo with ID: %s", ID)

	// Perform the delete operation in the repository
	err := s.repository.Delete(ID)

	// Check for errors
	if err != nil {
		// Log the error
		s.logger.Sugar().Errorf("error deleting todo with ID: %s", ID)
		return err
	}

	// Log successful deletion of the TODO item
	s.logger.Sugar().Debugf("todo deleted with ID: %s", ID)
	return nil
}
