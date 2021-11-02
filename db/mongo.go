package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
	InsertOne(string, string, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(string, string, interface{}) (*mongo.DeleteResult, error)
}

type MongoDB struct{}

func (a MongoDB) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		log.Panicf("Error Connecting to MongoDB:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	var result map[string]interface{}
	collection.FindOne(ctx, filter).Decode(&result)
	return result
}

func (a MongoDB) InsertOne(db_name string, collection_name string, filter interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		log.Panicf("Error Connecting to MongoDB:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	return collection.InsertOne(ctx, filter)
}

func (a MongoDB) DeleteOne(db_name string, collection_name string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		log.Panicf("Error Connecting to MongoDB:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	return collection.DeleteOne(ctx, filter)
}

func MongoClient(ctx context.Context) (*mongo.Client, error) {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, err
}
