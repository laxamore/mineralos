package db

import (
	"context"
	"fmt"
	"os"

	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBInterface interface {
	Find(*mongo.Client, string, string, interface{}) ([]map[string]interface{}, error)
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
	InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(*mongo.Client, string, string, interface{}) (*mongo.DeleteResult, error)
	IndexesCreateOne(*mongo.Client, string, string, mongo.IndexModel) (string, error)
	IndexesDropOne(*mongo.Client, string, string, string) (bson.Raw, error)
	IndexesReplaceOne(*mongo.Client, string, string, mongo.IndexModel) (string, error)
	IndexesReplaceMany(*mongo.Client, string, string, []mongo.IndexModel) (string, error)
	UpdateOne(*mongo.Client, string, string, interface{}, interface{}) (*mongo.UpdateResult, error)
}

type MongoDB struct{}

func (a MongoDB) Find(client *mongo.Client, db_name string, collection_name string, filter interface{}) ([]map[string]interface{}, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return []map[string]interface{}{{}}, err
	}

	var results []map[string]interface{}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			Log.Printf("%v", err)
		}

		results = append(results, elem)
	}

	return results, err
}

func (a MongoDB) FindOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) map[string]interface{} {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	var result map[string]interface{}
	collection.FindOne(ctx, filter).Decode(&result)
	return result
}

func (a MongoDB) InsertOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (*mongo.InsertOneResult, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	return collection.InsertOne(ctx, filter)
}

func (a MongoDB) DeleteOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (DeleteResult *mongo.DeleteResult, err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	ctx := context.TODO()

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

func (a MongoDB) IndexesCreateOne(client *mongo.Client, db_name string, collection_name string, indexModel mongo.IndexModel) (string, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	return collection.Indexes().CreateOne(ctx, indexModel)
}

func (a MongoDB) IndexesDropOne(client *mongo.Client, db_name string, collection_name string, indexName string) (bson.Raw, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	return collection.Indexes().DropOne(ctx, indexName)
}

func (a MongoDB) IndexesReplaceOne(client *mongo.Client, db_name string, collection_name string, indexModel mongo.IndexModel) (string, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

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

func (a MongoDB) IndexesReplaceMany(client *mongo.Client, db_name string, collection_name string, indexModel []mongo.IndexModel) ([]string, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

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

func (a MongoDB) UpdateOne(client *mongo.Client, db_name string, collection_name string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// client, err := MongoClient(ctx)

	// if err != nil {
	// 	Log.Printf("error Connecting to mongodb:\n%v", err)
	// }

	ctx := context.TODO()

	collection := client.Database(db_name).Collection(collection_name)
	return collection.UpdateOne(ctx, filter, update)
}

func MongoClient(ctx context.Context) (*mongo.Client, error) {
	DB_SERVER := os.Getenv("DB_SERVER")
	if os.Getenv("DOCKER") == "true" {
		DB_SERVER = "mongodb"
	}
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		DB_SERVER,
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
