package repository

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/juliocnsouzadev/kafka-ish/model"
)

func getMessage() model.Message {
	id := uuid.New().String()
	message := model.Message{
		Id:        id,
		Topic:     "tests_topic",
		Timestamp: time.Now(),
		Content:   "content",
	}
	return message
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
