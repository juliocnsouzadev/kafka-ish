package repository

import (
	"strings"
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpectedMessageType string

const (
	ExpectedMessageOK         ExpectedMessageType = "ok"
	ExpectedMessageDuplicated ExpectedMessageType = "duplicated"
	duplicatedKeyErrorMessage                     = "duplicate key error"
)

type MongoMessageTestStage struct {
	t                   *testing.T
	createdMessages     []model.Message
	expectedMessageType ExpectedMessageType
	foundMessages       []model.Message
	repo                *MongoMessageRepository[model.Message]
	overrideCollection  OverrideCollection
	topic               string
	err                 error
}

func NewMongoMessageTestStage(t *testing.T) (*MongoMessageTestStage, *MongoMessageTestStage, *MongoMessageTestStage) {
	repo := NewMongoMessageRepository(nil)
	overrideCollection := func(collection *mongo.Collection) {
		repo.collection = collection
	}
	stage := &MongoMessageTestStage{t: t, repo: repo, overrideCollection: overrideCollection}
	return stage, stage, stage
}

func (s *MongoMessageTestStage) and() *MongoMessageTestStage {
	return s
}

func (s *MongoMessageTestStage) a_message(expected ExpectedMessageType) *MongoMessageTestStage {
	s.createdMessages = []model.Message{getMessage()}
	s.expectedMessageType = expected
	return s
}

func (s *MongoMessageTestStage) the_message_is_created() *MongoMessageTestStage {
	switch s.expectedMessageType {
	case ExpectedMessageOK:
		assertInsertFunc := func(t *testing.T, err error) {
			require.Nil(t, err)
		}
		for _, msg := range s.createdMessages {
			MongoMockInsertSuccess[model.Message](s.t, "insert message", msg, s.repo, s.overrideCollection, assertInsertFunc)
		}
		s.err = nil
		break
	case ExpectedMessageDuplicated:
		assertInsertFunc := func(t *testing.T, err error) {
			require.NotNil(t, err)
			require.True(t, strings.Contains(err.Error(), duplicatedKeyErrorMessage))
		}
		for _, msg := range s.createdMessages {
			MongoMockInsertFailDuplicatedKey[model.Message](s.t, "fail duplicated key", msg, s.repo, s.overrideCollection, assertInsertFunc)
		}
		s.err = errors.Errorf("%s", duplicatedKeyErrorMessage)
		break
	default:
		panic("invalid expected message type")
	}

	return s
}

func (s *MongoMessageTestStage) no_errors_happen() *MongoMessageTestStage {
	require.Nil(s.t, s.err)
	return s
}

func (s *MongoMessageTestStage) fail_duplicated_key() *MongoMessageTestStage {
	require.NotNil(s.t, s.err)
	require.True(s.t, strings.Contains(s.err.Error(), duplicatedKeyErrorMessage))
	return s
}

func (s *MongoMessageTestStage) some_existing_messages_in_the_topic(topic string) *MongoMessageTestStage {
	s.topic = topic
	s.createdMessages = []model.Message{getMessage(), getMessage()}
	for _, msg := range s.createdMessages {
		msg.Topic = topic
	}
	return s
}

func (s *MongoMessageTestStage) find_messages_for_topic(topic string) *MongoMessageTestStage {
	assertFindByTopicFunc := func(t *testing.T, messages []model.Message, err error) {
		require.Nil(t, err)
		require.Equal(t, len(s.createdMessages), len(messages))
		for i, msg := range messages {
			require.Equal(t, s.createdMessages[i].Id, msg.Id)
			require.Equal(t, s.createdMessages[i].Topic, msg.Topic)
			require.Equal(t, s.createdMessages[i].Content, msg.Content)
		}
	}

	// exec
	MongoMockMessagesFindByTopicSuccess(s.t, "find messages by topic", s.createdMessages, topic, s.repo, s.overrideCollection, assertFindByTopicFunc)
	s.foundMessages = s.createdMessages
	return s
}

func (s *MongoMessageTestStage) found_messages_matches_existing_messages() *MongoMessageTestStage {

	require.Equal(s.t, len(s.createdMessages), len(s.foundMessages))
	for i, msg := range s.foundMessages {
		require.Equal(s.t, s.createdMessages[i].Id, msg.Id)
		require.Equal(s.t, s.createdMessages[i].Topic, msg.Topic)
		require.Equal(s.t, s.createdMessages[i].Content, msg.Content)
	}
	return s
}
