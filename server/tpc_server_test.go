package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {

	lines := []string{
		string(getCommand(CommandPublish)),
		string(getCommand(CommandPublish)),
	}

	mockProducer := producer.NewMockProducer()

	mockListener := NewMockTcpListener(lines...)
	defer mockListener.Close()

	done := make(chan bool)

	tcpServer, err := NewTCPServer(mockProducer, mockListener, done)
	require.NoError(t, err)

	tcpServer.Start()

	require.Equal(t, 1, mockListener.Actions["accept"])
	require.Equal(t, 1, mockListener.Actions["close"])

	require.Equal(t, 1, mockListener.reader.Actions["setKeepAlive"])
	require.Equal(t, 1, mockListener.reader.Actions["reader"])
	require.Equal(t, 1, mockListener.reader.Actions["close"])

	require.Equal(t, 2, mockProducer.Actions["publish"])
	require.Equal(t, 1, mockProducer.Actions["cancel"])

}

func getCommand(cmdType CommandType) []byte {
	message, _ := json.Marshal(getMessage())
	command := Command{
		Type: cmdType,
		Body: string(message),
	}
	commandBytes, _ := json.Marshal(command)
	return commandBytes
}

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
