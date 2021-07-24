package db

import (
	"context"

	"github.com/mises-id/sns/config/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB *mongo.Database
)

func SetupMongo(ctx context.Context) {
	client, err := mongo.Connect(ctx, options.Client().SetMaxPoolSize(30).ApplyURI(env.Envs.MongoURI))
	if err != nil {
		panic(err)
	}
	MongoDB = client.Database(env.Envs.DBName)
}

func DB() *mongo.Database {
	return MongoDB
}
