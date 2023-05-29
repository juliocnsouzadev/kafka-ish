package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/repository"
	"github.com/juliocnsouzadev/kafka-ish/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBProducer struct {
	client *mongo.Client
	repo   repository.Repository[model.Message]
}

func NewMongoProducer() *MongoDBProducer {
	return &MongoDBProducer{}
}

func (p *MongoDBProducer) initMongoClient() error {
	if p.client != nil {
		return nil
	}
	mongoClient, err := storage.ConnectToMongoDB()
	if err != nil {
		return fmt.Errorf("error connecting to mongoDB: %w", err)
	}
	p.client = mongoClient
	return nil
}

func (p *MongoDBProducer) initRepo() {
	if p.repo != nil {
		return
	}
	p.repo = repository.NewMongoMessageRepository(p.client)
}

func (p *MongoDBProducer) Publish(message model.Message) error {
	if err := p.initMongoClient(); err != nil {
		return fmt.Errorf("error initializing mongoDB client: %w", err)
	}
	p.initRepo()

	return p.repo.Insert(message)
}

func (p *MongoDBProducer) Cancel() {
	// create context with timeout to be able to disconnect from mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	defer func() {
		if err := p.client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()
}
