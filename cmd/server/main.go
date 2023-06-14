package main

import (
	"fmt"
	"net/http"
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
		WebPort:     getWebPort(),
	}

	producer := producer.NewProducer(config)

	//bulding tcp server
	tcpServer, err := buildTcpServer(config, producer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tcpServer.Cancel()
	tcpServer.Start()

	//build http server
	httpServer, err := buildHttpServer(config, producer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer httpServer.Cancel()
	httpServer.Start()

}

func buildHttpServer(config settings.Settings, producer producer.Producer) (*server.HttpServer, error) {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", config.WebPort),
	}
	done := make(chan bool)
	return server.NewHttpServer(producer, httpServer, done)
}

func buildTcpServer(config settings.Settings, producer producer.Producer) (*server.TcpServer, error) {
	var err error
	listener, err := server.NewDeafultTcpListener(config)
	if err != nil {
		return nil, err
	}
	done := make(chan bool)
	return server.NewTCPServer(producer, listener, done)
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

func getWebPort() string {
	port := os.Getenv("WEB_PORT")
	if port == "" {
		port = "80"
	}
	return port
}
