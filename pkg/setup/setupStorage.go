package setup

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func SetupDB(name string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoClient")

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB")

	db = client.Database(name)

	return nil
}

func SetupCollection(name string) (*mongo.Collection, error) {
	if db == nil {
		return nil, fmt.Errorf("DataBase is not configured. Call SetupDB() first")
	}

	return db.Collection(name), nil
}
