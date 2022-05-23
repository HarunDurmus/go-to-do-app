package todoapp

import (
	"context"
	"github.com/harundurmus/go-to-do-app/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	client *mongo.Client
	config config.MongoConfig
}

func (r repository) UpsertInitialData(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client *mongo.Client, config config.MongoConfig) Repository {
	return &repository{
		client: client,
		config: config,
	}
}
