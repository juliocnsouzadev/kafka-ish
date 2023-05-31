package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/stretchr/testify/require"
)

func TestInsert_File_Success(t *testing.T) {
	//given
	message := getMessage()
	repo := NewFileMessageRepository()

	defer cleanup(repo.storagePath, message.Topic)

	// exec
	err := repo.Insert(message)

	// assert
	require.Nil(t, err)
}

func TestFindByTopic_File_Success(t *testing.T) {
	//given
	topic := "tests_topic"

	repo := NewFileMessageRepository()
	defer cleanup(repo.storagePath, topic)

	expectedMessages := createMessages(repo)

	// exec
	messages, err := repo.FindByTopic(topic)

	// assert
	require.Nil(t, err)
	require.Equal(t, len(expectedMessages), len(messages))
	for i, msg := range messages {
		require.Equal(t, expectedMessages[i].Id, msg.Id)
		require.Equal(t, expectedMessages[i].Topic, msg.Topic)
		require.Equal(t, expectedMessages[i].Content, msg.Content)
	}

}

func createMessages(repo *FileMessageRepository[model.Message]) []model.Message {
	messageOne := getMessage()
	messageTwo := getMessage()
	expectedMessages := []model.Message{messageOne, messageTwo}
	for _, msg := range expectedMessages {
		err := repo.Insert(msg)
		if err != nil {
			panic(err)
		}
	}
	return expectedMessages
}

func cleanup(storagePath, topic string) {
	filePath := fmt.Sprintf("%s/%s.topic", storagePath, topic)
	e := os.Remove(filePath)
	if e != nil {
		fmt.Printf("error removing file %s\n", filePath)
	}
}
