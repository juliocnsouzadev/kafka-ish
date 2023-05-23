package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInsert_MongoDB_Success(t *testing.T) {
	//given
	message := getMessage()

	repo, overrideCollectionFunc := getMongoRepoAndCollection()

	assertInsertFunc := func(t *testing.T, err error) {
		require.Nil(t, err)
	}

	// exec
	MongoMockInsertSuccess[model.Message](t, "insert message", message, repo, overrideCollectionFunc, assertInsertFunc)
}

func TestInsert_MongoDB_Fail_DulicatedKey(t *testing.T) {
	//given
	message := getMessage()

	repo, overrideCollectionFunc := getMongoRepoAndCollection()

	assertInsertFunc := func(t *testing.T, err error) {
		require.NotNil(t, err)
	}

	// exec
	MongoMockInsertFailDuplicatedKey[model.Message](t, "fail duplicated key", message, repo, overrideCollectionFunc, assertInsertFunc)
}

func getMessage() model.Message {
	id := uuid.New().String()
	message := model.Message{
		Id:        id,
		Topic:     "topic",
		Timestamp: time.Now(),
		Content:   "content",
	}
	return message
}

func getMongoRepoAndCollection() (*MongoMessageRepository[model.Message], OverrideCollection) {
	repo := NewMongoMessageRepository(nil)

	OverrideCollection := func(collection *mongo.Collection) {
		repo.collection = collection
	}
	return repo, OverrideCollection
}
