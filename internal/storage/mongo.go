package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var MongoClient *mongo.Client

func ConnectMongoDB() *mongo.Client {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASS")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@bitbird.hxfijor.mongodb.net/?retryWrites=true&w=majority&appName=bitbird", mongoUser, mongoPassword)).SetServerAPIOptions(serverAPI)
	opts.Monitor = otelmongo.NewMonitor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB...")
	MongoClient = client
	return client
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}

func CloseMongoDB() {
	if err := MongoClient.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}
