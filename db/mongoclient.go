package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoClient(ctx context.Context) (*mongo.Client, error) {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return client, err
	}

	return client, err
}
