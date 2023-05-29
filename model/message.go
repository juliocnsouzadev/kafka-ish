package model

import (
	"time"
)

type Message struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Topic     string    `bson:"topic,omitempty" json:"topic"`
	Timestamp time.Time `bson:"timestamp,omitempty" timestpajson:"timestamp"`
	Content   string    `bson:"content,omitempty" json:"content"`
}
