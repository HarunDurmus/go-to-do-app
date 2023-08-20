package couchbase_test

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/harundurmus/go-to-do-app/internal/todo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *CouchbaseTestSuite) TestInsertTodo() {
	s.T().Run("given a todo  then it should insert todo", func(t *testing.T) {
		expectedTodo := todo.Todo{
			Plate: "34 ts 20",
		}
		_, err := s.bucket.Collection(_todoCollection).Get(expectedTodo.Plate, nil)
		assert.NotNil(t, err)

		err = s.todoRepository.Upsert(expectedTodo.Plate, expectedTodo)
		assert.Nil(t, err)

		doc, err := s.bucket.Collection(_todoCollection).Get(expectedTodo.Plate, nil)
		assert.Nil(t, err)
		var actualTodo todo.Todo
		err = doc.Content(&actualTodo)
		assert.Nil(t, err)
		assert.Equal(t, &expectedTodo, &actualTodo)
	})
}

func (s *CouchbaseTestSuite) TestGetTodoByID() {
	s.T().Run("given a todo place  then it should return todo by given plate", func(t *testing.T) {
		expectedTodo := todo.Todo{
			Plate: "34 ts 50",
		}

		_, err := s.bucket.Collection(_todoCollection).Get(expectedTodo.Plate, nil)
		assert.NotNil(t, err)

		_, err = s.bucket.Collection(_todoCollection).Upsert(
			fmt.Sprint(expectedTodo.Plate),
			expectedTodo,
			&gocb.UpsertOptions{},
		)
		assert.Nil(t, err)

		actual, err := s.todoRepository.GetById(expectedTodo.Plate)
		assert.Nil(t, err)
		assert.Equal(t, &expectedTodo, actual)
	})
	s.T().Run("given a todo plate  then it should return err by given plate if it does not exist", func(t *testing.T) {
		expectedTodo := todo.Todo{
			Plate: "34 tks 23",
		}
		_, err := s.todoRepository.GetById(expectedTodo.Plate)
		assert.NotNil(t, err)
	})
}
