package model

import (
	"time"
)

type Message struct {
	Id        string    `json:"id"`
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}
