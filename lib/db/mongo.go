package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB *mongo.Database
)

func SetupMongo(ctx context.Context, url, dbName string) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	MongoDB = client.Database(dbName)
}

func DB() *mongo.Database {
	return MongoDB
}
