package storage

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoUrl = "mongodb://mongo:27017"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: decode(os.Getenv("MONGO_USER")),
		Password: decode(os.Getenv("MONGO_PASSWORD")),
	})
	return mongo.Connect(context.TODO(), clientOptions)
}

func decode(value string) string {
	result, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Panic(err)
	}
	return string(result)
}
