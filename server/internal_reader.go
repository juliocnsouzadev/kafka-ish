package server

import (
	"io"
	"net"
)

type InternalReader interface {
	SetKeepAlive(keepalive bool) error
	Close() error
	Reader() io.Reader
}

type DefaultTpcReader struct {
	reader *net.TCPConn
}

func (t *DefaultTpcReader) SetKeepAlive(keepalive bool) error {
	return t.reader.SetKeepAlive(keepalive)
}

func (t *DefaultTpcReader) Close() error {
	return t.reader.Close()
}

func (t *DefaultTpcReader) Reader() io.Reader {
	return t.reader
}
