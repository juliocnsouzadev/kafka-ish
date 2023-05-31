package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/juliocnsouzadev/kafka-ish/model"
)

const (
	DefaultStoragePath = "../db-data/filestore"
)

type FileMessageRepository[T model.Message] struct {
	storagePath string
}

func NewFileMessageRepository() *FileMessageRepository[model.Message] {
	path := os.Getenv("FILE_STORAGE_PATH")
	if path == "" {
		path = DefaultStoragePath
	}

	r := &FileMessageRepository[model.Message]{storagePath: path}
	return r
}

func (r *FileMessageRepository[T]) Insert(message model.Message) error {
	file, err := r.openTopicFile(message.Topic)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	raw, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	raw = append(raw, byte('\n'))
	_, err = file.Write(raw)
	if err != nil {
		return fmt.Errorf("error writing message: %v", err)
	}
	return nil
}

func (r *FileMessageRepository[T]) FindByTopic(topic string) ([]model.Message, error) {
	file, err := r.openTopicFile(topic)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var entries []model.Message
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("EOF: %v\n", err)
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line: %v", err)
		}
		entry := model.Message{}
		err = json.Unmarshal(line, &entry)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling line: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *FileMessageRepository[T]) openTopicFile(topic string) (*os.File, error) {
	filePath := fmt.Sprintf("%s/%s.topic", r.storagePath, topic)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	return file, err
}
