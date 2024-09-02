package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = "mongodb+srv://priyansh3006:monti2022@cluster0.vevncst.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

func ConnectToDatabase() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//client ,err := mongo.Connect(ctx, options.Client().ApplyURI())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("error in intializing context")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("error in pinging database")
	}
	fmt.Println("Successfully connected to mongo db")
	return client
}

var Db *mongo.Client = ConnectToDatabase()

func GetCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database("CrudApi").Collection(collectionName)
	return collection
}
