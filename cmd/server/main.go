package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/juliocnsouzadev/kafka-ish/settings"
)

var prod producer.Producer

func main() {
	storageType := os.Args[1]

	if storageType == "" {
		fmt.Println("storage type not informed switchin to filestore")
		storageType = string(settings.FileStore)
	}

	prod = producer.NewProducer(settings.Settings{
		StorageType: settings.StorageType(strings.ToLower(storageType)),
	})

	defer func() {
		prod.Cancel()
	}()

}
