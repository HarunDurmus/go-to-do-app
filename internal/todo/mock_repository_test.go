// $GOPATH/src/github.com/harundurmus/go-to-do-app/internal/todo/repository_mock.go

package todo

// MockRepository is a mocks implementation of the Repository interface for testing purposes.
type MockRepository struct {
	Data  map[string]Todo // Data stores TODO items by their ID
	Error error           // Error simulates errors in repository operations
}

func (m *MockRepository) InsertOrUpdateTodo(todo Todo) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepository) DeleteTodo(ID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepository) GetById(ID string) (*Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepository) List() ([]Todo, error) {
	//TODO implement me
	panic("implement me")
}

// Upsert simulates the upsert operation in the mocks repository.
func (m *MockRepository) Upsert(ID string, todoItem Todo) error {
	// Simulate the upsert operation
	m.Data[ID] = todoItem
	return m.Error
}

// Delete simulates the delete operation in the mocks repository.
func (m *MockRepository) Delete(ID string) error {
	// Simulate the delete operation
	delete(m.Data, ID)
	return m.Error
}
