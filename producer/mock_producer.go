package producer

import (
	"github.com/juliocnsouzadev/kafka-ish/model"
)

type MockProducer struct {
	Actions map[string]int
}

func NewMockProducer() *MockProducer {
	return &MockProducer{
		Actions: map[string]int{"publish": 0, "cancel": 0},
	}
}

func (p *MockProducer) Publish(message model.Message) error {
	p.Actions["publish"]++
	return nil
}

func (p *MockProducer) Cancel() {
	defer func() {
		p.Actions["cancel"]++
	}()
}
