package client

import (
	"context"
	"log"

	"github.com/harundurmus/go-to-do-app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDb(config config.MongoConfig) *mongo.Client {
	credential := options.Credential{Username: config.Username, Password: config.Password}
	clientOptions := options.Client().ApplyURI(config.URL).SetAuth(credential)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	return client
}
