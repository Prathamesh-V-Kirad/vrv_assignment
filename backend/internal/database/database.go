package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	dbName     = "golang_db" 
)


func Connect() *mongo.Client {

	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	database = client.Database(dbName)
	log.Println("Successfully connected to MongoDB")
	return client
}

func GetCollection(name string) *mongo.Collection {
	if database == nil {
		log.Fatal("Database connection not initialized")
	}
	return database.Collection(name)
}
