package producer

import (
	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/settings"
)

type Producer interface {
	Publish(message model.Message) error
	Cancel()
}

func NewProducer(config settings.Settings) Producer {
	switch config.StorageType {
	case settings.MongoDB:
		return NewMongoProducer()
	case settings.FileStore:
		return NewFileStoreProducer()
	default:
		return nil
	}
}
