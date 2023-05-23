package repository

import (
	"context"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName   = storage.DbName
	collectionName = "messages"
)

type MongoMessageRepository[T model.Message] struct {
	client     *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
	collection *mongo.Collection
}

func NewMongoMessageRepository(client *mongo.Client) *MongoMessageRepository[model.Message] {
	r := &MongoMessageRepository[model.Message]{client: client}
	r.buildContext()
	r.buildCollection()
	return r
}

func (r *MongoMessageRepository[T]) Insert(message model.Message) error {
	defer r.cancel()

	_, err := r.collection.InsertOne(r.ctx, message)
	return err
}

func (r *MongoMessageRepository[T]) buildContext() {
	if r.ctx != nil || r.cancel != nil {
		return
	}
	r.ctx, r.cancel = context.WithTimeout(context.Background(), 5*time.Second)
}

func (r *MongoMessageRepository[T]) buildCollection() {
	if r.collection != nil || r.client == nil {
		return
	}
	r.collection = r.client.Database(databaseName).Collection(collectionName)
}
