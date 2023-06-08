package repository

import (
	"strings"
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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
	topic               string
	testLabel           string
	mock                *MongoMock[model.Message]
}

func NewMongoMessageTestStage(t *testing.T) (*MongoMessageTestStage, *MongoMessageTestStage, *MongoMessageTestStage) {
	repo := NewMongoMessageRepository(nil)
	overrideCollection := func(collection *mongo.Collection) {
		repo.collection = collection
	}
	stage := &MongoMessageTestStage{t: t, mock: NewMongoMock[model.Message](t, repo, overrideCollection)}
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
		for _, msg := range s.createdMessages {
			s.testLabel = "insert message"
			s.mock.InsertSuccess(s.testLabel, msg)
		}
		break
	case ExpectedMessageDuplicated:
		for _, msg := range s.createdMessages {
			s.testLabel = "fail duplicated key"
			s.mock.InsertFailDuplicatedKey(s.testLabel, msg)
		}
		break
	default:
		panic("invalid expected message type")
	}
	return s
}

func (s *MongoMessageTestStage) no_errors_happen() *MongoMessageTestStage {
	err := s.mock.GetOutput(s.testLabel, ErrorOutputType)
	require.Nil(s.t, err, s.testLabel)
	return s
}

func (s *MongoMessageTestStage) fail_duplicated_key() *MongoMessageTestStage {
	errRaw := s.mock.GetOutput(s.testLabel, ErrorOutputType)
	require.NotNil(s.t, errRaw, s.testLabel)

	err, ok := errRaw.(error)
	require.True(s.t, ok, s.testLabel)

	require.True(s.t, strings.Contains(err.Error(), duplicatedKeyErrorMessage), s.testLabel)
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
	s.testLabel = "find messages by topic"
	mockedData := make([]bson.D, len(s.createdMessages))
	for i, msg := range s.createdMessages {
		bsonData := bson.D{
			{Key: "_id", Value: msg.Id},
			{Key: "topic", Value: msg.Topic},
			{Key: "timestamp", Value: msg.Timestamp},
			{Key: "content", Value: msg.Content},
		}
		mockedData[i] = bsonData
	}
	s.mock.FindByTopicSuccess(s.testLabel, mockedData, topic)
	return s
}

func (s *MongoMessageTestStage) found_messages_matches_existing_messages() *MongoMessageTestStage {
	err := s.mock.GetOutput(s.testLabel, ErrorOutputType)
	require.Nil(s.t, err, s.testLabel)

	foundMessagesRaw := s.mock.GetOutput(s.testLabel, SuccessOutputType)
	require.NotNil(s.t, foundMessagesRaw, s.testLabel)

	foundMessages, ok := foundMessagesRaw.([]model.Message)
	require.True(s.t, ok, s.testLabel)

	require.Equal(s.t, len(s.createdMessages), len(foundMessages))
	for i, msg := range foundMessages {
		require.Equal(s.t, s.createdMessages[i].Id, msg.Id, s.testLabel)
		require.Equal(s.t, s.createdMessages[i].Topic, msg.Topic, s.testLabel)
		require.Equal(s.t, s.createdMessages[i].Content, msg.Content, s.testLabel)
	}
	return s
}
