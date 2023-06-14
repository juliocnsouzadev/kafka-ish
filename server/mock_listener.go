package server

import (
	"time"
)

type MockTcpListener struct {
	reader  *MockReader
	Actions map[string]int
}

func NewMockTcpListener(lines ...string) *MockTcpListener {
	return &MockTcpListener{
		reader: NewMockReader(lines...),
		Actions: map[string]int{
			"close":  0,
			"accept": 0,
		},
	}
}

func (t *MockTcpListener) Close() error {
	t.reader.Close()
	t.Actions["close"]++
	return nil
}

func (t *MockTcpListener) Accept(deadline time.Time) (InternalReader, error) {
	t.Actions["accept"]++
	return t.reader, nil
}
