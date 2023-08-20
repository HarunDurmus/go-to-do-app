package todo

import "github.com/couchbase/gocb/v2"

const _todo = "todo"

type repository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

func NewRepository(cluster *gocb.Cluster, bucket *gocb.Bucket) Repository {
	return &repository{
		cluster: cluster,
		bucket:  bucket,
	}
}

func (r *repository) Upsert(ID string, todo Todo) error {
	_, err := r.bucket.Collection(_todo).Upsert(
		ID,
		todo,
		&gocb.UpsertOptions{},
	)
	return err
}

func (r *repository) Delete(ID string) error {
	_, err := r.bucket.Collection(_todo).Remove(
		ID,
		&gocb.RemoveOptions{},
	)
	return err
}

func (r *repository) List() ([]Todo, error) {
	query := "SELECT * FROM " + _todo
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

func (r *repository) GetById(ID string) (*Todo, error) {
	result, err := r.bucket.Collection(_todo).Get(
		ID,
		&gocb.GetOptions{},
	)
	if err != nil {
		return nil, err
	}

	var todo Todo
	err = result.Content(&todo)
	return &todo, err
}
