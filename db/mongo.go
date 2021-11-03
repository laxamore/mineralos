package db

import (
	"context"
	"fmt"
	"os"
	"time"

	Log "github.com/laxamore/mineralos/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
	InsertOne(string, string, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(string, string, interface{}) (*mongo.DeleteResult, error)
	IndexesCreateOne(string, string, mongo.IndexModel) (string, error)
	IndexesDropOne(string, string, string) (bson.Raw, error)
	IndexesReplaceOne(string, string, mongo.IndexModel) (string, error)
	IndexesReplaceMany(string, string, []mongo.IndexModel) (string, error)
}

type MongoDB struct{}

func (a MongoDB) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		Log.Panicf("Error Connecting to MongoDB:\n%v", err)
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
		Log.Panicf("Error Connecting to MongoDB:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	return collection.InsertOne(ctx, filter)
}

func (a MongoDB) DeleteOne(db_name string, collection_name string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	var DeleteResult *mongo.DeleteResult

	if err != nil {
		return DeleteResult, fmt.Errorf("error connecting to mongodb:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	DeleteResult, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return DeleteResult, fmt.Errorf("DeleteOne Failed:\n%v", err)
	}
	return DeleteResult, err
}

func (a MongoDB) IndexesCreateOne(db_name string, collection_name string, indexModel mongo.IndexModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		return "", fmt.Errorf("error connecting to mongodb:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	return collection.Indexes().CreateOne(ctx, indexModel)
}

func (a MongoDB) IndexesDropOne(db_name string, collection_name string, indexName string) (bson.Raw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	return collection.Indexes().DropOne(ctx, indexName)
}

func (a MongoDB) IndexesReplaceOne(db_name string, collection_name string, indexModel mongo.IndexModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		return "", fmt.Errorf("error connecting to mongodb:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	res, err := collection.Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		var input map[string]interface{}
		filterBytes, _ := bson.Marshal(indexModel.Keys)
		bson.Unmarshal(filterBytes, &input)

		indexName := func(m map[string]interface{}) []string {
			keys := make([]string, len(m))
			i := 0
			for k := range m {
				keys[i] = k
				i++
			}
			return keys
		}(input)[0]

		_, err := collection.Indexes().DropOne(ctx, indexName+"_1")

		if err != nil {
			return "failed replacing index" + indexName + "_1", err
		}
		return collection.Indexes().CreateOne(ctx, indexModel)
	}
	return res, err
}

func (a MongoDB) IndexesReplaceMany(db_name string, collection_name string, indexModel []mongo.IndexModel) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := MongoClient(ctx)

	if err != nil {
		return []string{""}, fmt.Errorf("error connecting to mongodb:\n%v", err)
	}

	collection := client.Database(db_name).Collection(collection_name)
	res, err := collection.Indexes().CreateMany(ctx, indexModel)

	if err != nil {
		for i := range indexModel {
			var input map[string]interface{}
			filterBytes, _ := bson.Marshal(indexModel[i].Keys)
			bson.Unmarshal(filterBytes, &input)

			indexName := func(m map[string]interface{}) []string {
				keys := make([]string, len(m))
				i := 0
				for k := range m {
					keys[i] = k
					i++
				}
				return keys
			}(input)[0]

			_, err := collection.Indexes().DropOne(ctx, indexName+"_1")

			if err != nil {
				return []string{"failed replacing index " + indexName + "_1"}, err
			}
		}
		return collection.Indexes().CreateMany(ctx, indexModel)
	}
	return res, err
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
