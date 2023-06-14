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

type TcpServerTestStage struct {
	t            *testing.T
	lines        []string
	err          error
	mockProducer *producer.MockProducer
	mockListener *MockTcpListener
	tcpServer    *TcpServer
}

func NewTcpServerTestStage(t *testing.T) (*TcpServerTestStage, *TcpServerTestStage, *TcpServerTestStage) {

	stage := &TcpServerTestStage{t: t}
	return stage, stage, stage
}

func (s *TcpServerTestStage) and() *TcpServerTestStage {
	return s
}

func (s *TcpServerTestStage) two_publish_commands() *TcpServerTestStage {
	s.lines = []string{
		string(getCommand(CommandPublish)),
		string(getCommand(CommandPublish)),
	}
	return s
}

func (s *TcpServerTestStage) a_tcp_server() *TcpServerTestStage {
	s.mockProducer = producer.NewMockProducer()

	s.mockListener = NewMockTcpListener(s.lines...)

	done := make(chan bool)

	s.tcpServer, s.err = NewTCPServer(s.mockProducer, s.mockListener, done)
	return s
}

func (s *TcpServerTestStage) tcp_server_starts() *TcpServerTestStage {
	require.NoError(s.t, s.err)
	s.tcpServer.Start()

	return s
}

func (s *TcpServerTestStage) interactions_are_correct() *TcpServerTestStage {

	require.Equal(s.t, 1, s.mockListener.Actions["accept"])
	require.Equal(s.t, 1, s.mockListener.Actions["close"])

	require.Equal(s.t, 1, s.mockListener.reader.Actions["setKeepAlive"])
	require.Equal(s.t, 1, s.mockListener.reader.Actions["reader"])
	require.Equal(s.t, 1, s.mockListener.reader.Actions["close"])

	require.Equal(s.t, 2, s.mockProducer.Actions["publish"])
	require.Equal(s.t, 1, s.mockProducer.Actions["cancel"])
	return s
}

func (s *TcpServerTestStage) no_erros_ocurr() *TcpServerTestStage {
	require.NoError(s.t, s.err)
	return s
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
