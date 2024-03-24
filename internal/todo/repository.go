package todo

import "github.com/couchbase/gocb/v2"

// Repository interface defines methods for interacting with TODO data.
type Repository interface {
	GetById(ID string) (*Todo, error)  // GetById retrieves a TODO item by its ID.
	List() ([]Todo, error)             // List returns a list of all TODO items.
	Upsert(ID string, todo Todo) error // Upsert inserts or updates a TODO item.
	Delete(ID string) error            // Delete removes a TODO item by its ID.
}

const todoCollection = "todo" // todoCollection defines the name of the collection where TODO items are stored.

// repository represents the implementation of the Repository interface.
type repository struct {
	cluster *gocb.Cluster // cluster is the Couchbase cluster instance.
	bucket  *gocb.Bucket  // bucket is the Couchbase bucket instance.
}

// NewRepository creates a new instance of Repository with the given Couchbase cluster and bucket.
func NewRepository(cluster *gocb.Cluster, bucket *gocb.Bucket) Repository {
	return &repository{
		cluster: cluster,
		bucket:  bucket,
	}
}

// Upsert inserts or updates a TODO item in the database.
func (r *repository) Upsert(ID string, todo Todo) error {
	_, err := r.bucket.Collection(todoCollection).Upsert(
		ID,
		todo,
		&gocb.UpsertOptions{},
	)
	return err
}

// Delete removes a TODO item by its ID.
func (r *repository) Delete(ID string) error {
	_, err := r.bucket.Collection(todoCollection).Remove(
		ID,
		&gocb.RemoveOptions{},
	)
	return err
}

// List returns a list of all TODO items.
func (r *repository) List() ([]Todo, error) {
	query := "SELECT * FROM " + todoCollection
	rows, err := r.cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Row(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// GetById retrieves a TODO item by its ID.
func (r *repository) GetById(ID string) (*Todo, error) {
	result, err := r.bucket.Collection(todoCollection).Get(ID, nil)
	if err != nil {
		return nil, err
	}

	var todo Todo
	if err := result.Content(&todo); err != nil {
		return nil, err
	}

	return &todo, nil
}
