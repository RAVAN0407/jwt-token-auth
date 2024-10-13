package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("No Mongodb URI in .env file")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("failed to establish connection with mongoDB: ", err)
	}
	log.Println("Connected to mongoDB")
	return client
}

var Client *mongo.Client = DBInstance()

func OpenConnection(c *mongo.Client, collection string) *mongo.Collection {
	// TODO Database value
	return c.Database("Cluster0").Collection(collection)
}
