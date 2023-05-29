package main

import (
	"os"

	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/juliocnsouzadev/kafka-ish/settings"
)

var prod producer.Producer

func main() {
	storageType := os.Args[1]

	prod = producer.NeProducer(settings.Settings{
		StorageType: settings.StorageType(storageType),
	})

	defer func() {
		prod.Cancel()
	}()

}
