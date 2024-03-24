// $GOPATH/src/github.com/harundurmus/go-to-do-app/internal/todo/service_test.go

package todo_test

import (
	"errors"
	"github.com/harundurmus/go-to-do-app/internal/todo"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestInsertOrUpdateTodo(t *testing.T) {
	// Initialize a todo logger
	mockLogger := zaptest.NewLogger(t)

	// Initialize a todo repository
	mockRepo := &todo.MockRepository{
		Data: make(map[string]todo.Todo),
	}

	// Initialize the service with the todo repository and logger
	service := todo.NewService(mockRepo, mockLogger)

	// Create a sample TODO item
	todoItem := todo.Todo{
		ID:   "1",
		Task: "Sample TODO",
	}

	// Test inserting a TODO item
	err := service.InsertOrUpdateTodo(todoItem)
	assert.NoError(t, err)
	assert.Equal(t, todoItem, mockRepo.Data["1"], "TODO item should be inserted")

	// Test updating a TODO item
	todoItem.Task = "Updated TODO"
	err = service.InsertOrUpdateTodo(todoItem)
	assert.NoError(t, err)
	assert.Equal(t, todoItem, mockRepo.Data["1"], "TODO item should be updated")

	// Test error handling
	mockRepo.Error = errors.New("todo error")
	err = service.InsertOrUpdateTodo(todoItem)
	assert.Error(t, err)
}

func TestDeleteTodo(t *testing.T) {
	// Initialize a todo logger
	mockLogger := zaptest.NewLogger(t)

	// Initialize a todo repository with some data
	mockData := map[string]todo.Todo{
		"1": {ID: "1", Title: "Sample TODO"},
		"2": {ID: "2", Title: "Another TODO"},
	}
	mockRepo := &todo.MockRepository{Data: mockData}

	// Initialize the service with the todo repository and logger
	service := todo.NewService(mockRepo, mockLogger)

	// Test deleting a TODO item
	err := service.DeleteTodo("1")
	assert.NoError(t, err)
	assert.NotContains(t, mockRepo.Data, "1", "TODO item should be deleted")

	// Test error handling
	mockRepo.Error = errors.New("todo error")
	err = service.DeleteTodo("2")
	assert.Error(t, err)
}
