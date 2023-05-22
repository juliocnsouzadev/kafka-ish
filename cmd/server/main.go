package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	mongoClient, err := storage.ConnectToMongoDB()
	if err != nil {
		log.Println(fmt.Errorf("error connecting to mongoDB: %w", err))
		log.Panic(err)
	}
	client = mongoClient

	// create context with timeout to be able to disconnect from mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

}
