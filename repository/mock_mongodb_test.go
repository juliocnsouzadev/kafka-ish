package repository

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const (
	ErrorOutputType   OutputType = "error"
	SuccessOutputType OutputType = "success"
)

type OverrideCollection func(*mongo.Collection)

type OutputType string

type MongoMock[T any] struct {
	repo               Repository[T]
	mt                 *mtest.T
	overrideCollection OverrideCollection
	outputs            map[string]map[OutputType]any
}

func NewMongoMock[T any](t *testing.T, repo Repository[T], overrideCollection OverrideCollection) *MongoMock[T] {
	return &MongoMock[T]{
		repo:               repo,
		overrideCollection: overrideCollection,
		mt:                 mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock)),
		outputs:            make(map[string]map[OutputType]any),
	}
}

func (this *MongoMock[T]) GetOutput(label string, outpType OutputType) any {
	if data, ok := this.outputs[label]; ok {
		return data[outpType]
	}
	return nil
}

func (this *MongoMock[T]) InsertSuccess(label string, data T) {

	defer this.mt.Close()

	this.mt.Run(label, func(mt *mtest.T) {
		this.overrideCollection(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := this.repo.Insert(data)
		this.outputs[label] = make(map[OutputType]any)
		this.outputs[label][ErrorOutputType] = err
	})
}

func (this *MongoMock[T]) InsertFailDuplicatedKey(label string, data T) {

	defer this.mt.Close()

	this.mt.Run(label, func(mt *mtest.T) {
		this.overrideCollection(mt.Coll)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		err := this.repo.Insert(data)
		this.outputs[label] = make(map[OutputType]any)
		this.outputs[label][ErrorOutputType] = err
	})
}

func (this *MongoMock[T]) FindByTopicSuccess(label string, expectedData []bson.D, topic string) {
	defer this.mt.Close()

	this.mt.Run(label, func(mt *mtest.T) {
		this.overrideCollection(mt.Coll)

		responses := make([]bson.D, len(expectedData)+1)
		for i, data := range expectedData {

			if i == 0 {
				response := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, data)
				responses[i] = response
			} else {
				response := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, data)
				responses[i] = response
			}

		}
		//kill cursor
		responses[len(expectedData)] = mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(responses...)
		messages, err := this.repo.FindByTopic(topic)
		this.outputs[label] = make(map[OutputType]any)
		this.outputs[label][SuccessOutputType] = messages
		this.outputs[label][ErrorOutputType] = err
	})
}
