package repository

import (
	"testing"

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

func TestFindByTopic_MongoDB_Success(t *testing.T) {
	//given
	messageOne := getMessage()
	messageTwo := getMessage()
	expectedMessages := []model.Message{messageOne, messageTwo}

	topic := "topic"

	repo, overrideCollectionFunc := getMongoRepoAndCollection()

	assertFindByTopicFunc := func(t *testing.T, messages []model.Message, err error) {
		require.Nil(t, err)
		require.Equal(t, len(expectedMessages), len(messages))
		for i, msg := range messages {
			require.Equal(t, expectedMessages[i].Id, msg.Id)
			require.Equal(t, expectedMessages[i].Topic, msg.Topic)
			require.Equal(t, expectedMessages[i].Content, msg.Content)
		}
	}

	// exec
	MongoMockMessagesFindByTopicSuccess(t, "find messages by topic", expectedMessages, topic, repo, overrideCollectionFunc, assertFindByTopicFunc)
}

func getMongoRepoAndCollection() (*MongoMessageRepository[model.Message], OverrideCollection) {
	repo := NewMongoMessageRepository(nil)

	OverrideCollection := func(collection *mongo.Collection) {
		repo.collection = collection
	}
	return repo, OverrideCollection
}
