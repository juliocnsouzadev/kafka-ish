package repository

import (
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/stretchr/testify/require"
)

type FileMessageTestStage struct {
	t               *testing.T
	createdMessages []model.Message
	foundMessages   []model.Message
	repo            *FileMessageRepository[model.Message]
	topic           string
	err             error
}

func NewFileMessageTestStage(t *testing.T) (*FileMessageTestStage, *FileMessageTestStage, *FileMessageTestStage) {
	repo := NewFileMessageRepository()
	stage := &FileMessageTestStage{t: t, repo: repo}
	return stage, stage, stage
}

func (s *FileMessageTestStage) and() *FileMessageTestStage {
	return s
}

func (s *FileMessageTestStage) a_message() *FileMessageTestStage {
	s.createdMessages = []model.Message{getMessage()}
	return s
}

func (s *FileMessageTestStage) the_message_is_created() *FileMessageTestStage {
	for _, msg := range s.createdMessages {
		s.err = s.repo.Insert(msg)
		s.topic = msg.Topic
		if s.err != nil {
			break
		}
	}
	return s
}

func (s *FileMessageTestStage) no_errors_happen() *FileMessageTestStage {
	defer cleanup(s.repo.storagePath, s.topic)

	require.Nil(s.t, s.err)
	return s
}

func (s *FileMessageTestStage) some_existing_messages_in_the_topic(topic string) *FileMessageTestStage {
	s.topic = topic
	s.createdMessages = createMessages(s.repo)
	return s
}

func (s *FileMessageTestStage) find_messages_for_topic(topic string) *FileMessageTestStage {
	s.foundMessages, s.err = s.repo.FindByTopic(topic)
	return s
}

func (s *FileMessageTestStage) found_messages_matches_existing_messages() *FileMessageTestStage {
	defer cleanup(s.repo.storagePath, s.topic)

	require.Equal(s.t, len(s.createdMessages), len(s.foundMessages))
	for i, msg := range s.foundMessages {
		require.Equal(s.t, s.createdMessages[i].Id, msg.Id)
		require.Equal(s.t, s.createdMessages[i].Topic, msg.Topic)
		require.Equal(s.t, s.createdMessages[i].Content, msg.Content)
	}
	return s
}
