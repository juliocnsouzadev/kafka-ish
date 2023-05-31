package repository

import (
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
