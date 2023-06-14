package server

import (
	"encoding/json"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/producer"
)

type MockHttpServer struct {
	Actions         map[string]int
	prod            producer.Producer
	commands        []Command
	commandHandlers map[CommandType]CommandHandler
}

func NewMockHttpServer(producer producer.Producer, commands []Command) (*MockHttpServer, error) {
	return &MockHttpServer{
		prod:     producer,
		commands: commands,
		Actions: map[string]int{
			"publish": 0,
		},
	}, nil
}

func (t *MockHttpServer) CommandHandlers() map[CommandType]CommandHandler {
	if t.commandHandlers == nil {
		t.commandHandlers = map[CommandType]CommandHandler{
			CommandPublish: t.publish,
			CommandConsume: t.consume,
			CommandClose:   t.close,
		}
	}
	return t.commandHandlers
}

func (t *MockHttpServer) publish(command Command) error {
	t.Actions["publish"]++
	message := model.Message{}

	if err := json.Unmarshal([]byte(command.Body), &message); err != nil {
		return err
	}

	return t.prod.Publish(message)
}

func (t *MockHttpServer) consume(command Command) error {
	panic("not implemented")
}

func (t *MockHttpServer) close(command Command) error {
	return nil
}

func (t *MockHttpServer) Cancel() {
	t.prod.Cancel()
}

func (t *MockHttpServer) Start() {
	defer t.Cancel()

	for _, command := range t.commands {
		handler := t.CommandHandlers()[command.Type]
		handler(command)
	}
}
