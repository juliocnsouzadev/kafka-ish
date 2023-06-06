package producer

import (
	"log"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/repository"
)

type FileStoreProducer struct {
	repo repository.Repository[model.Message]
}

func NewFileStoreProducer() *FileStoreProducer {
	return &FileStoreProducer{}
}

func (p *FileStoreProducer) initRepo() {
	if p.repo != nil {
		return
	}
	p.repo = repository.NewFileMessageRepository()
}

func (p *FileStoreProducer) Publish(message model.Message) error {
	p.initRepo()

	return p.repo.Insert(message)
}

func (p *FileStoreProducer) Cancel() {
	defer func() {
		log.Println("end of filestore producer")
	}()
}
