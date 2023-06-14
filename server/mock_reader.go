package server

import (
	"io"
	"strings"
)

type MockReader struct {
	reader  io.ReadCloser
	Actions map[string]int
}

func (t *MockReader) SetKeepAlive(keepalive bool) error {
	t.Actions["setKeepAlive"]++
	return nil
}

func (t *MockReader) Close() error {
	t.Actions["close"]++
	return t.reader.Close()
}

func (t *MockReader) Reader() io.Reader {
	t.Actions["reader"]++
	return t.reader
}

func NewMockReader(lines ...string) *MockReader {
	someString := ""
	for _, line := range lines {
		someString += line + "\n"
	}
	myReader := strings.NewReader(someString)

	reader := io.NopCloser(myReader)
	return &MockReader{
		reader: reader,
		Actions: map[string]int{
			"setKeepAlive": 0,
			"close":        0,
			"reader":       0,
		},
	}
}
