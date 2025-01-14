package setup

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var client *mongo.Client

func SetupDB(name string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	dbClient, err := mongo.Connect(context.TODO(), clientOptions)
	client = dbClient
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

func DisconnectDB() {
	client.Disconnect(context.TODO())
}

func GetCollection(name string) (*mongo.Collection, error) {
	if db == nil {
		return nil, fmt.Errorf("DataBase is not configured. Call SetupDB() first")
	}

	return db.Collection(name), nil
}
