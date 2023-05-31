package repository

import (
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type OverrideCollection func(*mongo.Collection)

type AssertInsert func(t *testing.T, err error)
type AssertFind[T any] func(t *testing.T, messages []T, err error)

func MongoMockInsertSuccess[T any](
	t *testing.T,
	label string,
	data T,
	repo Repository[T],
	overrideCollection OverrideCollection,
	assertFunc AssertInsert,
) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	defer mt.Close()

	mt.Run(label, func(mt *mtest.T) {
		overrideCollection(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.Insert(data)

		assertFunc(t, err)
	})
}

func MongoMockInsertFailDuplicatedKey[T any](
	t *testing.T,
	label string,
	data T,
	repo Repository[T],
	overrideCollection OverrideCollection,
	assertFunc AssertInsert,
) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	defer mt.Close()

	mt.Run(label, func(mt *mtest.T) {
		overrideCollection(mt.Coll)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		err := repo.Insert(data)

		assertFunc(t, err)
	})
}

func MongoMockMessagesFindByTopicSuccess(
	t *testing.T,
	label string,
	expectedData []model.Message,
	topic string,
	repo Repository[model.Message],
	overrideCollection OverrideCollection,
	assertFunc AssertFind[model.Message],
) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	defer mt.Close()

	mt.Run(label, func(mt *mtest.T) {
		overrideCollection(mt.Coll)

		responses := make([]bson.D, len(expectedData)+1)
		for i, data := range expectedData {

			bsonData := bson.D{
				{Key: "_id", Value: data.Id},
				{Key: "topic", Value: data.Topic},
				{Key: "timestamp", Value: data.Timestamp},
				{Key: "content", Value: data.Content},
			}
			if i == 0 {
				response := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bsonData)
				responses[i] = response
			} else {
				response := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bsonData)
				responses[i] = response
			}

		}
		//kill cursor
		responses[len(expectedData)] = mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(responses...)
		messages, err := repo.FindByTopic(topic)
		assertFunc(t, messages, err)
	})
}
