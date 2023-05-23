package repository

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type OverrideCollection func(*mongo.Collection)

type AssertInsert func(t *testing.T, err error)

type MapModelForMockUpdate func(data any) (bson.D, error)

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
