package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/juliocnsouzadev/kafka-ish/server"
	"github.com/juliocnsouzadev/kafka-ish/settings"
)

var prod producer.Producer

func main() {

	config := settings.Settings{
		StorageType: getStorageType(),
		TcpPort:     getTcpPort(),
	}

	tcpServer, err := server.NewTCPServer(config)

	defer tcpServer.Cancel()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tcpServer.Start()
}

func getTcpPort() string {
	port := os.Getenv("TCP_PORT")
	if port == "" {
		port = "9001"
	}
	return port
}

func getStorageType() settings.StorageType {
	storageType := os.Args[1]

	if storageType == "" {
		fmt.Println("storage type not informed switchin to filestore")
		storageType = string(settings.FileStore)
	}
	return settings.StorageType(strings.ToLower(storageType))
}
