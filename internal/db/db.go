package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client            *mongo.Client
	connectionTimeout = 10 * time.Second
)

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongodb:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func Collection() *mongo.Collection {
	if client == nil {
		log.Println("Database connection is not established")
	}
	return client.Database("url-shortener").Collection("urls")
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect database: %v", err)
	}
}
