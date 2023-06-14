package server

import (
	"encoding/json"
	"testing"

	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/stretchr/testify/require"
)

type HttpServerTestStage struct {
	t            *testing.T
	commands     []Command
	err          error
	mockProducer *producer.MockProducer
	httpServer   *MockHttpServer
}

func NewHttpServerTestStage(t *testing.T) (*HttpServerTestStage, *HttpServerTestStage, *HttpServerTestStage) {

	stage := &HttpServerTestStage{t: t}
	return stage, stage, stage
}

func (s *HttpServerTestStage) and() *HttpServerTestStage {
	return s
}

func (s *HttpServerTestStage) two_publish_commands() *HttpServerTestStage {
	s.commands = []Command{
		s._getCommand(CommandPublish),
		s._getCommand(CommandPublish),
	}
	return s
}

func (s *HttpServerTestStage) a_http_server() *HttpServerTestStage {
	s.mockProducer = producer.NewMockProducer()

	s.httpServer, s.err = NewMockHttpServer(s.mockProducer, s.commands)
	return s
}

func (s *HttpServerTestStage) http_server_starts() *HttpServerTestStage {
	require.NoError(s.t, s.err)
	s.httpServer.Start()

	return s
}

func (s *HttpServerTestStage) interactions_are_correct() *HttpServerTestStage {

	require.Equal(s.t, 2, s.mockProducer.Actions["publish"])
	require.Equal(s.t, 1, s.mockProducer.Actions["cancel"])
	return s
}

func (s *HttpServerTestStage) no_erros_ocurr() *HttpServerTestStage {
	require.NoError(s.t, s.err)
	return s
}

func (s *HttpServerTestStage) _getCommand(cmdType CommandType) Command {
	message, _ := json.Marshal(getMessage())
	command := Command{
		Type: cmdType,
		Body: string(message),
	}
	return command
}
