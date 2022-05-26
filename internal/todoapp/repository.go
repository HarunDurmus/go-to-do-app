package todoapp

import (
	"context"
	"github.com/harundurmus/go-to-do-app/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	client *mongo.Client
	config config.MongoConfig
}

func NewRepository(client *mongo.Client, config config.MongoConfig) Repository {
	return &repository{
		client: client,
		config: config,
	}
}
func (r repository) UpsertInitialData(ctx context.Context) error {
	var _ = r.client.Database(r.config.DBName).Collection(r.config.Collection)

	return nil
}

func (r repository) CreateTask(ctx context.Context) error {
	var _ = r.client.Database(r.config.DBName).Collection(r.config.Collection)

	return nil
}

func (r *repository) GetAll(ctx context.Context) (locations []*Location, err error) {
	coll := r.client.Database(r.config.DBName).Collection(r.config.Collection)
	cursor, _ := coll.Find(ctx, bson.M{})
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, err
	}
	return locations, err
}

func (r *repository) checkCollectionExist(ctx context.Context, coll *mongo.Collection) error {
	locations, err := r.GetAll(ctx)
	if len(locations) > 0 {
		return errors.New("database already initialized")
	}
	return err
}
